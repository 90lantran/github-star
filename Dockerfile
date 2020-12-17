FROM golang:latest as builder

WORKDIR /go/src

COPY . .

RUN mkdir -p bin

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o bin/github-star cmd/*.go

#####
FROM blang/golang-alpine:latest

WORKDIR /bin

COPY --from=builder /go/src/bin/github-star .

EXPOSE 8080

ENTRYPOINT ["/bin/github-star"]