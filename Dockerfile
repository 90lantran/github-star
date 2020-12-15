FROM golang:latest

RUN go get github.com/90lantran/github-star

ENTRYPOINT ["/go/bin/github-star/cmd"]

EXPOSE 8080