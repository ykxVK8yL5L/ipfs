FROM golang:1.16

WORKDIR /go/src/app
COPY ./app /go/src/app

RUN go mod vendor
RUN go build main.go
CMD["./main"]

VOLUME ["/go/src/app/config"]
