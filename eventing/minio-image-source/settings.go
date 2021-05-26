package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	MinioServer string `envconfig:"MINIO_SERVER"`
	MinioBucket string `envconfig:"MINIO_BUCKET"`
	MinioKey    string `envconfig:"MINIO_KEY"`
	MinioSecret string `envconfig:"MINIO_SECRET"`

	ServiceName      string `envconfig:"SERVICE_NAME"`
	ServiceNamespace string `envconfig:"SERVICE_NAMESPACE"`

	CloudEventSink       string `envconfig:"K_SINK"`
	CloudEventType       string
	CloudEventSource     string
	CloudEventDataSchema string
}

func NewSettings() Settings {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		log.Fatalln(err)
	}

	settings.CloudEventType = "dev.zufardhiyaulhaq.eventing.image-processing.minio.bucket.image.create"
	settings.CloudEventDataSchema = "https://minio.zufardhiyaulhaq.com/publics/schema/minio-data-processing-schema.json"
	settings.CloudEventSource = buildEventSource(settings.ServiceNamespace, settings.ServiceName)
	return settings
}

func buildEventSource(namespace, name string) string {
	var eventType string

	eventType = fmt.Sprintf("https://dev.zufardhiyaulhaq.eventing.%s.%s", namespace, name)

	return eventType
}
