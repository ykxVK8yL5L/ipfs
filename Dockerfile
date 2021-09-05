FROM golang:1.16

WORKDIR /go/src/app
COPY ./app /go/src/app

RUN go mod vendor

VOLUME ["/go/src/app/config"]
ENTRYPOINT ["go", "run", "main.go"]


