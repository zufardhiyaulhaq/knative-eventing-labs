package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/http2"
)

type Data struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
}

func main() {
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
		log.Fatalln(err)
	}

	for notificationInfo := range minioClient.ListenBucketNotification(context.Background(), settings.MinioBucket, "", "", []string{"s3:ObjectCreated:*"}) {
		for _, event := range notificationInfo.Records {
			if validateExtension(event.S3.Object.Key) {
				data := Data{
					Bucket: event.S3.Bucket.Name,
					Name:   event.S3.Object.Key,
					Size:   event.S3.Object.Size,
				}
				dataEncoded, _ := json.Marshal(data)
				fmt.Println(string(dataEncoded))

				client, err := cloudevents.NewClientHTTP()
				if err != nil {
					log.Fatalf("failed to create client, %v", err)
				}

				event := cloudevents.NewEvent()
				event.SetID(uuid.New().String())
				event.SetTime(time.Now())
				event.SetType(settings.CloudEventType)
				event.SetSource(settings.CloudEventSource)
				event.SetDataSchema(settings.CloudEventDataSchema)
				event.SetData(cloudevents.ApplicationJSON, dataEncoded)

				ctx := cloudevents.ContextWithTarget(context.Background(), settings.CloudEventSink)
				if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
					log.Printf("failed to send, %v", result)
				}
			} else {
				log.Printf("Extension not supported for file: %s", event.S3.Object.Key)
			}
		}
	}
}

func validateExtension(filename string) bool {
	var isValid bool = false

	filename = strings.ToLower(filepath.Ext(filename))
	extensionName := strings.Replace(filepath.Ext(filename), ".", "", -1)

	if extensionName == "png" || extensionName == "jpg" || extensionName == "jpeg" {
		isValid = true
	}

	return isValid
}
