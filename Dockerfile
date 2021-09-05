FROM golang:1.16

WORKDIR /go/src/app
COPY ./app /go/src/app

RUN go mod vendor
RUN go run main.go

VOLUME ["/go/src/app/config"]
EXPOSE 8080