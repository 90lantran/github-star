PACKAGES:=$(shell go list ./... | grep 'api')

build:
	go build -o server cmd/*.go

local-server: create-image
	@echo "running docker image of server ,listen at localhost:8080 ..."
	docker run -p 8080:8080 github-stars

create-image:
	docker build --tag github-stars .

push-new-image: create-image
	docker login && docker tag github-stars:latest 90lantran/my-server:latest  && docker push 90lantran/my-server:latest

unit-test:
	@echo "run unit-test"
	go clean -testcache && go test -cover $(PACKAGES) -coverprofile report/unit-test.out

code-coverage: unit-test
	@echo display code coverage
	go tool cover -html=report/unit-test.out

fmt:
	@echo "formatting code using go fmt"
	go fmt ./...

vet: 
	@echo "vetting code using go vet"
	go vet ./...

check-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	@echo "Create swagger yaml"
	GO111MODULE=on go mod vendor  && swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml

k8s-deploy:
	minikube start --vm-driver=virtualbox;\
	kubectl apply -f deployment/server.yaml;\
	sleep 10;\
	minikube service my-server-service;

k8s-delete:
	minikube delete

	


