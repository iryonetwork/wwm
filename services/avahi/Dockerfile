FROM alpine:3.8

ARG AVAHI_VERSION=0.7
ARG AVAHI_RELEASE=r1

RUN env && apk add --no-cache \
    avahi=$AVAHI_VERSION-$AVAHI_RELEASE \
    avahi-tools=$AVAHI_VERSION-$AVAHI_RELEASE 

ADD avahi-daemon.conf /etc/avahi/avahi-daemon.conf
ADD entrypoint.sh /opt/entrypoint.sh

ENTRYPOINT /opt/entrypoint.sh
