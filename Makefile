.PHONY: local-server
local-server:
	@echo "running server locally ..."
	docker-compose up server

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=on go mod vendor  && GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml

unit-test:
	go test -p=1 -cover $(PACKAGES) > unit-test.out;\
	code=$$?;\
	go-junit-report < unit-test.out > unit-report.xml; \
	cat unit-test.out; \
	grep -e 'FAIL' unit-test.out; \
	exit $${code}

fmt:
	@echo "formatting code using go fmt"
	go fmt ./...

vet: 
	@echo "vetting code using go vet"
	go vet ./...

# .PHONY: golint
# golint:
# 	@echo "running golangci-lint"
# 	go lint