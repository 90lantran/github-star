FROM golang:latest

RUN go get github.com/90lantran/github-star

EXPOSE 8080

ENTRYPOINT ["/go/bin/github-star"]