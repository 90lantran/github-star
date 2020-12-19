# github-star
This is simple API written in Go to get number of github stars for a list of organization/repository.

## API document
Details of endpoints, response and request are in this [swagger](https://90lantran.github.io/swagger-github-stars/).

## Unit-test
unit-test were writtent with [goconvey](https://github.com/smartystreets/goconvey) which gives details about testing scenarios of unit-test. You may need to [install](https://github.com/smartystreets/goconvey#installation) it first if you have not had it.

To run unit-test: 
```
$ make unit-test 
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
```
$ make local-server
```

## Deployment in minikube: 
The deployment config is in deployment/server.yaml. Here are the steps:
- Install hypervisor (pick one of them): 
    - hyperkit: brew install hyperkit) 
    - virtualbox: brew install --cask virtualbox
- Install minikube : 
```
$ brew install minikube
```
- Start minikuke with driver (it is virtualbox for me):
```
$ minikube start --vm-driver=virtualbox
```
- Create deployment and service:
```
$kubectl apply -f deployment/server.yaml
```
-  minikube service my-server-service

After the last command, it should open a browser like below which has the base-url for the API. We can hit base-url/health to do health check for API.

After done with it, we can tear down the cluster
```
$ minikube delete
```
## Future work:
- connect o github with go-github https://github.com/google/go-github, no authentication for now
- no authentication to use github API because I dont want to commit by credentails -> rate limit
- rate limit: how to avoid it need to log in

## Swagger docs:
I use go-swagger to generate API document.
Create swagger yaml: 
```
$ make swagger
```
Show swagger ui:
```
$ make serve-swagger
```

