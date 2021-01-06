# github-star
This API can get number of stars for a list of format `organization/repository` from github. If an element of the list is `not registered` as an orgization in github or a repository is `not belong` to a organization or `invallid format`, it will be not be processed.
<br/>
In order to connect to github, I use [go-github](https://github.com/google/go-github). This provides a client to list all repositories belong to an organization.

## API document
Details of endpoints, responses and requests are in this [swagger](https://90lantran.github.io/swagger-github-stars/). It is `highly recommended` to take a look at this swagger to have general idea how this API works.

## Unit test
unit-test were written with [goconvey](https://github.com/smartystreets/goconvey) which gives details about testing scenarios of unit-test. You may need to [install](https://github.com/smartystreets/goconvey#installation) it first if you have not had it.

To run unit-test: 
```
$ make unit-test
run unit-test
go clean -testcache && go test -cover github.com/90lantran/github-star/pkg/api -coverprofile report/unit-test.out
ok      github.com/90lantran/github-star/pkg/api        1.237s  coverage: 96.4% of statements 
```

To show code-coverage:
```
$ make code-coverage
```

## Build docker image
To dockerzied the API:
```
$ make create-image
```

## Run API locally
To run server as docker container
```
$ make local-server
```

## Deployment in minikube: 
The deployment config is in `deployment/server.yaml`.
<br/>
1.If you already had minikube and virtualbox installed, you can deploy to minikube right away.
```
$ make k8s-deploy
minikube delete
🔥  Deleting "minikube" in virtualbox ...
💀  Removed all traces of the "minikube" cluster.
(base) github-star$ make k8s-deploy
minikube start --vm-driver=virtualbox;\
	kubectl apply -f deployment/server.yaml;\
	sleep 10;\
	minikube service my-server-service;
😄  minikube v1.15.1 on Darwin 10.13.6
✨  Using the virtualbox driver based on user configuration
👍  Starting control plane node minikube in cluster minikube
🔥  Creating virtualbox VM (CPUs=2, Memory=2200MB, Disk=20000MB) ...
🐳  Preparing Kubernetes v1.19.4 on Docker 19.03.13 ...
🔎  Verifying Kubernetes components...
🌟  Enabled addons: default-storageclass, storage-provisioner
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
deployment.apps/my-server created
service/my-server-service created
|-----------|-------------------|-------------|-----------------------------|
| NAMESPACE |       NAME        | TARGET PORT |             URL             |
|-----------|-------------------|-------------|-----------------------------|
| default   | my-server-service |        8080 | http://192.168.99.106:30000 |
|-----------|-------------------|-------------|-----------------------------|
🎉  Opening service default/my-server-service in default browser...

```
2.If you need to install minikube and hypervisor, you can follow these steps. It will take like 15-20mins.
- Install hypervisor (pick one of them): 
    - hyperkit: brew install hyperkit) 
    - virtualbox: brew install --cask virtualbox
- Install minikube : 
```
$ brew install minikube
```
- Start minikube with driver (virtualbox or hyperkit):
```
$ minikube start --vm-driver=virtualbox
```
- Create deployment and service:
```
$ kubectl apply -f deployment/server.yaml
```
- Start minikube
```
$ service my-server-service
```

After the last command, it should open your browser and hit the `base-url`=`http://public.ip.minikube:30000`of the API. We can hit other endpoints  base-url/health and base-url/get-stars. For some reasons, you may see **This site can’t be reached** (maybe slow network). It just means minikube deployment or service is not ready to serve, please be patient and give it some time and refesh the page until you see 
```json
{
    message: "Server is deployed in minikube"
}
```

When you are done with using minikube, we can tear down the cluster
```
$ make k8s-delete
```

## Future work:
- Connect to github with [go-github](https://github.com/google/go-github) with no authentication for now. It will affect the rate limit when interacting with github APIs.
- Add [zap-log](https://github.com/uber-go/zap) from Uber for structured logs for monitoring
- Add Jenkinsfile for CICD

## Swagger docs:
I use [go-swagger](https://github.com/go-swagger/go-swagger) to generate API document.
Create swagger yaml: 
```
$ make swagger
```
Show swagger UI in browser:
```
$ make serve-swagger
```
