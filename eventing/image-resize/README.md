# Image Resize
Image Resize is serverless application that receive CloudEvents from Minio Image Source, get the image from Minio, Resize, and upload to the Minio bucket.

## Building
```bash
cd echo
docker build -t {username}/image-resize . 
docker push {username}/image-resize
```

## Deploying
Minio Image Source application require this variable to be exist:
- **MINIO_SERVER**, Minio server
- **MINIO_SOURCE_BUCKET**, Minio source bucket where getting the image
- **MINIO_DESTINATION_BUCKET**, Minio destination bucket where uploading the resize image
- **MINIO_KEY**, Minio key
- **MINIO_SECRET**, Minio secret

please check `.env` for more information. we can create Kubernetes secret and mount the secret as a environment variable in the application.
```bash
kubectl apply -f secret.yaml
```

Deploying the application
```bash
kubectl apply -f service.yaml

kn service list
NAME            URL                                                    LATEST                AGE   CONDITIONS   READY   REASON
image-resize    http://image-resize.serverless.svc.cluster.local       image-resize-00001    8m43s   3 OK / 3     True 
```

By default, image-resize is not subscribe to any event, we need to tell Knative to forward the event to the image-resize application via Trigger object. For more information, please check trigger.yaml.
```bash
kubectl apply -f trigger.yaml

kn trigger list
NAME                                       BROKER    SINK                 AGE   CONDITIONS   READY   REASON
image-resize-trigger                       default   ksvc:image-resize    17m   5 OK / 5     True
```
