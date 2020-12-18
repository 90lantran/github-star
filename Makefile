PACKAGES:=$(shell go list ./... | grep 'api')

.PHONY: local-server
local-server: create-image
	@echo "running docker image of server ,listen at localhost:8080 ..."
	docker run -p 8080:8080 github-stars

create-image:
	docker build --tag github-stars .

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
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	@echo "Create swagger yaml"
	GO111MODULE=on go mod vendor  && swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml

