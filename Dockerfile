FROM alpine:latest

ADD ./bin/tls/certs /usr/local/share/ca-certificates/

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates