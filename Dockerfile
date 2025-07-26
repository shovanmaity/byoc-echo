FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN make bin
RUN chmod +x bin/app

FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.8
COPY --from=builder /app/bin/app /usr/local/bin/app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD [ "app" ]
