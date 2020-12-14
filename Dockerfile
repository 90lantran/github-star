FROM golang:latest

RUN go get github.com/90lantran/github-star

CMD /go/bin/hellorest

EXPOSE 8080