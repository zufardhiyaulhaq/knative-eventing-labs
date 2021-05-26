# Serving
In this lab, we will try to deploying echo application. To deploy to knative, you can apply YAML files to the knative via Kubernetes API or using `kn` command line.

## Building
You can write application with language that you like. In this example, we have echo application in the echo directory. To build the image, you can run
```bash
cd echo
docker build -t {username}/echo . 
docker push {username}/echo
```

Knative by default create environment variables that you can use, for example:
- **PORT**, the default port that should be configured for your application. the default in 8080
- **K_SERVICE**, knative service name, from the example, should be `echo`.
- **K_REVISION**, knative revision name, from the example, should be `echo-00001`.

## Deploying
### Using kn
Knative provide CLI to easily deploy your application.
```bash
kn service create echo --image zufardhiyaulhaq/echo --namespace serverless
Creating service 'echo' in namespace 'serverless':

  0.581s The Route is still working to reflect the latest desired specification.
  1.131s ...
  1.461s Configuration "echo" is waiting for a Revision to become ready.
 30.562s ...
 31.237s Ingress has not yet been reconciled.
 35.079s Waiting for load balancer to be ready
 38.670s Certificate route-5fa58ba0-0d6d-4b40-82f0-551e1297a219 is not ready.
 52.207s Ingress has not yet been reconciled.
 53.687s Waiting for load balancer to be ready
 57.697s Ready to serve.

Service 'echo' created to latest revision 'echo-00001' is available at URL:
https://echo.serverless.knative.zufardhiyaulhaq.tech
```

If you try to access the service
```bash
curl https://echo.serverless.knative.zufardhiyaulhaq.tech                                                    
[knative-labs] Hi from echo (rev: echo-00001)
```

### Using Kubernetes API
If you like to deploy your application in Kubernetes way, you can also deploy via YAML file
```bash
kubectl apply -f service.yaml
kubectl get ksvc
NAME             URL                                                     LATESTCREATED          LATESTREADY            READY   REASON
echo             https://echo.serverless.knative.zufardhiyaulhaq.tech    echo-00001             echo-00001             True  
```
