FROM golang:1.16 AS builder
WORKDIR /go/src/app
COPY ./app /go/src/app
RUN go mod vendor
RUN go build

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /go/src/app ./
CMD ["./app"]  
