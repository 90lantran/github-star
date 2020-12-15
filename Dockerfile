FROM golang:latest

RUN go get github.com/90lantran/github-star

ENTRYPOINT ["/go/bin/github-star"]

EXPOSE 8080