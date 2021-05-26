# Eventing

In this lab, we will try to create simple Image processing application that using CloudEvents. There is 4 application in this use cases:
- **minio-image-source**, custom event publisher/source that watch event in the Minio bucket and send CloudEvents notification if image is created in the bucket.
- **image-resize**,  application that receive CloudEvents from Minio Image Source, get the image from Minio, Resize image, and upload to the Minio bucket.
- **image-grayscale**, application that receive CloudEvents from Image Resize, get the image from Minio, apply grayscale to the image, and upload to the Minio bucket.
- **event-display**, application written by Knative team to accept CloudNative events and display in the log.

## Topology
![image-processing](/static/image/image-processing-knative.png?raw=true)

## Integration Testing
Test upload the the bucket
```bash
mc cp static/image/image-processing-knative.png serverless/image-processing/
...age/image-processing-knative.png:  39.04 KiB / 39.04 KiB ┃▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓┃ 664.77 KiB/s 0s

mc ls serverless/image-processing/
[2021-05-26 05:36:08 CEST]  39KiB image-processing-knative.png
```

Image processing will upload two image that resize and grayscaled.
```bash
mc ls serverless/image-processing-output/
[2021-05-26 05:36:29 CEST]  33KiB image-processing-knative-resize-grayscale.png
[2021-05-26 05:36:20 CEST]  33KiB image-processing-knative-resize.png
```

We can check the CloudEvents log in event-display
```bash
kubectl logs deployment/event-display-00001-deployment -c user-container

☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.zufardhiyaulhaq.eventing.image-processing.minio.bucket.image.create
  source: https://dev.zufardhiyaulhaq.eventing.serverless.minio-image-source
  id: ff08584e-65db-427e-bd10-498bd9bf8a27
  time: 2021-05-26T03:41:22.946004518Z
  dataschema: https://minio.zufardhiyaulhaq.com/publics/schema/minio-data-processing-schema.json
  datacontenttype: application/json
Extensions,
  knativearrivaltime: 2021-05-26T03:41:23.001386854Z
Data,
  {
    "bucket": "image-processing",
    "name": "image-processing-knative.png",
    "size": 39980
  }
☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.zufardhiyaulhaq.eventing.image-processing.resize
  source: https://dev.zufardhiyaulhaq.eventing.serverless.image-resize
  id: a5a436b7-4486-42b0-bbde-1e408b858df1
  time: 2021-05-26T03:41:33.005512503Z
  dataschema: https://minio.zufardhiyaulhaq.com/publics/schema/minio-data-processing-schema.json
  datacontenttype: application/json
Extensions,
  knativearrivaltime: 2021-05-26T03:41:33.083113248Z
Data,
  {
    "bucket": "image-processing-output",
    "name": "image-processing-knative-resize.png",
    "size": -1
  }
☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.zufardhiyaulhaq.eventing.image-processing.grayscale
  source: https://dev.zufardhiyaulhaq.eventing.serverless.image-grayscale
  id: a591ad4e-05b9-4b2d-8fff-ec4e51a8b39f
  time: 2021-05-26T03:41:34.026350979Z
  dataschema: https://minio.zufardhiyaulhaq.com/publics/schema/minio-data-processing-schema.json
  datacontenttype: application/json
Extensions,
  knativearrivaltime: 2021-05-26T03:41:34.100628621Z
Data,
  {
    "bucket": "image-processing-output",
    "name": "image-processing-knative-resize-grayscale.png",
    "size": -1
  }
```
