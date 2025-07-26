FROM --platform=linux/amd64 golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN make bin
RUN chmod +x bin/app

FROM --platform=linux/amd64 alpine:latest as certs
RUN apk --update add ca-certificates

FROM --platform=linux/amd64 registry.access.redhat.com/ubi8/ubi-minimal:8.8
COPY --from=builder /app/bin/app /usr/local/bin/app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD [ "app" ]
