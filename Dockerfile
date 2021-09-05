# syntax=docker/dockerfile:1
FROM golang:1.16 AS builder
WORKDIR /go/src/app
COPY ./app /go/src/app
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY ./app /root
COPY --from=builder /go/src/app/app ./
CMD ["./app"]  
