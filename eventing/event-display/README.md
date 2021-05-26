# Event Display

Event Display is application written by Knative team to accept CloudNative events and display in the log.

## Deploying
```bash
kubectl apply -f service.yaml

kn service list
NAME            URL                                                    LATEST                AGE   CONDITIONS   READY   REASON
event-display   http://event-display.serverless.svc.cluster.local      event-display-00001   46s   3 OK / 3     True  
```

to send the event from the broker to the Event Display service, apply the trigger object
```bash
kubectl apply -f trigger.yaml

kn trigger list
NAME                                       BROKER    SINK                 AGE   CONDITIONS   READY   REASON
grayscale-event-display-trigger            default   ksvc:event-display   9s    5 OK / 5     True    
minio-image-source-event-display-trigger   default   ksvc:event-display   9s    5 OK / 5     True    
resize-event-display-trigger               default   ksvc:event-display   9s    5 OK / 5     True 
```

If other application trigger and event, you can check from the logs
```bash
kubectl logs deployment/event-display-00001-deployment -c user-container
☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.zufardhiyaulhaq.eventing.image-processing.minio.bucket.image.create
  source: https://dev.zufardhiyaulhaq.eventing.serverless.minio-image-source
  id: 78df2a38-0561-4f93-97ce-052dc6215a63
  time: 2021-05-26T01:40:39.728447562Z
  dataschema: https://minio.zufardhiyaulhaq.com/publics/schema/minio-data-processing-schema.json
  datacontenttype: application/json
Extensions,
  knativearrivaltime: 2021-05-26T01:40:39.774943775Z
Data,
  {
    "bucket": "image-processing",
    "name": "photo-18-may-2021.jpg",
    "size": 272849
  }
```
