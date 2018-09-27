#!/bin/sh
set -e

# replace host variables
AVAHI_HOST=${AVAHI_HOST:-"iryo"}
AVAHI_DOMAIN=${AVAHI_DOMAIN:-"local"}
sed -i -e "s/^host-name=.*/host-name=${AVAHI_HOST}/" /etc/avahi/avahi-daemon.conf
sed -i -e "s/^domain-name=.*/domain-name=${AVAHI_DOMAIN}/" /etc/avahi/avahi-daemon.conf

# remove junk from previous run
rm -r /var/run/avahi-daemon

avahi-daemon
