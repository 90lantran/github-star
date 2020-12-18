PACKAGES:=$(shell go list ./... | grep 'api')

.PHONY: local-server
local-server:
	@echo "running server locally,listen at :8080 ..."
	go run cmd/*.go

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
	GO111MODULE=on go mod vendor  && GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml