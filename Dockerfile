FROM golang:alpine

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o /bin/github-star cmd/*.go

EXPOSE 8080

ENTRYPOINT ["/bin/github-star"]