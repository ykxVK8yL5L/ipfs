FROM golang:1.16 AS builder
WORKDIR /go/src/app
COPY ./app /go/src/app
RUN go mod vendor
RUN CGO_ENABLED=1 GOOS=linux go build -a -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /go/src/app/app ./
COPY ./app /root
CMD ["./app"]  
