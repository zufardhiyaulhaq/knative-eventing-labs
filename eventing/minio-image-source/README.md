# Minio Image Source
Minio Image Source is an custom event publisher/source that watch event in the Minio bucket and send CloudEvents notification if image is created in the bucket.

## Building
```bash
cd echo
docker build -t {username}/minio-image-source . 
docker push {username}/minio-image-source
```

## Deploying
Minio Image Source application require this variable to be exist:
- **MINIO_SERVER**, Minio server
- **MINIO_BUCKET**, Minio bucket
- **MINIO_KEY**, Minio key
- **MINIO_SECRET**, Minio secret

please check `.env` for more information. we can create Kubernetes secret and mount the secret as a environment variable in the application.
```bash
kubectl apply -f secret.yaml
```

Minio Image Source is an event publisher in the event-driven concept. Knative expose lot of event publisher, but if we want to create our own publisher, we can use [ContainerSource](https://knative.dev/docs/eventing/sources/containersource/) in the Knative concept. This application applied the concept. To apply the publisher
```bash
kubectl apply -f source.yaml
```

we can check with `kn` command line
```bash
kn source list
NAME                             TYPE              RESOURCE                               SINK             READY
minio-image-source               ContainerSource   containersources.sources.knative.dev   broker:default   True
minio-image-source-sinkbinding   SinkBinding       sinkbindings.sources.knative.dev       broker:default   True
```
