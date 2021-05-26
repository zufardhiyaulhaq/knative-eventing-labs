# Eventing

In this lab, we will try to create simple Image processing application that using CloudEvents. There is 4 application in this use cases:
- **minio-image-source**, custom event publisher/source that watch event in the Minio bucket and send CloudEvents notification if image is created in the bucket.
- **image-resize**,  application that receive CloudEvents from Minio Image Source, get the image from Minio, Resize image, and upload to the Minio bucket.
- **image-grayscale**, application that receive CloudEvents from Image Resize, get the image from Minio, apply grayscale to the image, and upload to the Minio bucket.
- **event-display**, application written by Knative team to accept CloudNative events and display in the log.

## Topology
![image-processing](/static/image/image-processing-knative.png?raw=true)
