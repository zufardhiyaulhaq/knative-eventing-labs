package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"log"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/net/http2"
)

type Data struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
}

func receive(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	settings := NewSettings()

	minioClient, err := minio.New(settings.MinioServer, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.MinioKey, settings.MinioSecret, ""),
		Secure: true,
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				NextProtos:         []string{"h2"},
			},
		},
	})

	if err != nil {
		log.Printf("Error opening minio client: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(500, "failed to open minio client: %s", err)
	}

	schemaLoader := gojsonschema.NewReferenceLoader(event.Context.GetDataSchema())
	documentLoader := gojsonschema.NewStringLoader(string(event.DataEncoded))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if !result.Valid() {
		log.Printf("failed validating event data from schema: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(400, "failed validating event data from schema: %s", err)
	}

	eventData := &Data{}
	if err := event.DataAs(eventData); err != nil {
		log.Printf("Error while extracting cloudevent Data: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(500, "failed to convert event data: %s", err)
	}

	imageObject, err := minioClient.GetObject(context.Background(), eventData.Bucket, eventData.Name, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Cannot get %s/%s in Minio: %s\n", eventData.Bucket, eventData.Name, err.Error())
		return nil, cloudevents.NewHTTPResult(500, "failed to get image: %s", err)
	}

	image, err := imaging.Decode(imageObject)
	if err != nil {
		log.Printf("Cannot decode image: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(500, "failed to decode image: %s", err)
	}

	imageGrayscale := imaging.Grayscale(image)

	buffer := new(bytes.Buffer)
	err = imaging.Encode(buffer, imageGrayscale, 1)
	if err != nil {
		log.Printf("Cannot encode image: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(500, "failed to encode image: %s", err)
	}

	imageObject = nil
	image = nil
	imageGrayscale = nil

	extensionName := filepath.Ext(eventData.Name)
	name := strings.TrimSuffix(eventData.Name, filepath.Ext(eventData.Name)) + "-grayscale"
	fullName := name + extensionName

	_, err = minioClient.PutObject(context.Background(), settings.MinioDestinationBucket, fullName, buffer, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Printf("Cannot upload image: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(500, "failed to put image: %s", err)
	}

	buffer.Reset()
	debug.FreeOSMemory()

	newEvent := cloudevents.NewEvent()
	newEvent.SetID(uuid.New().String())
	newEvent.SetTime(time.Now())
	newEvent.SetType(settings.CloudEventType)
	newEvent.SetSource(settings.CloudEventSource)
	newEvent.SetDataSchema(settings.CloudEventDataSchema)

	data := Data{
		Bucket: settings.MinioDestinationBucket,
		Name:   fullName,
		Size:   -1,
	}

	if err := newEvent.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}

	return &newEvent, nil
}

func main() {
	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatal(client.StartReceiver(context.Background(), receive))
}
