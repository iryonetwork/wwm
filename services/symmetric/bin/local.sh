#!/bin/bash

SYM_DIR=/opt/symmetric

# Fill-in values in engine file with environment variables
cp /opt/symmetric/enginesTemplates/* /opt/symmetric/engines
sed -i -e "s#^engine.name=local-<LOCATION_ID>#engine.name=local-${LOCATION_ID}#" /opt/symmetric/engines/local.properties
sed -i -e "s#^external.id=<LOCATION_ID>#external.id=${LOCATION_ID}#" /opt/symmetric/engines/local.properties
sed -i -e "s#^db.url=<DB_URL>#db.url=${DB_URL}#" /opt/symmetric/engines/local.properties
sed -i -e "s#^db.user=<DB_USER>#db.user=${DB_USER}#" /opt/symmetric/engines/local.properties
sed -i -e "s#^db.password=<DB_PASSWORD>#db.password=${DB_PASSWORD}#" /opt/symmetric/engines/local.properties
sed -i -e "s#^registration.url=<REGISTRATION_URL>#registration.url=${REGISTRATION_URL}#" /opt/symmetric/engines/local.properties
if [ "$CLOUD_SYMMETRIC_BASIC_AUTH_ENABLED" = true ] ; then
    sed -i -e "s#^http.basic.auth.username=<CLOUD_SYMMETRIC_BASIC_AUTH_USERNAME>#http.basic.auth.username=${CLOUD_SYMMETRIC_BASIC_AUTH_USERNAME}#" /opt/symmetric/engines/local.properties
    sed -i -e "s#^http.basic.auth.password=<CLOUD_SYMMETRIC_BASIC_AUTH_PASSWORD>#http.basic.auth.password=${CLOUD_SYMMETRIC_BASIC_AUTH_PASSWORD}#" /opt/symmetric/engines/local.properties
else
    sed -i -e "s#^http.basic.auth.username=<CLOUD_SYMMETRIC_BASIC_AUTH_USERNAME>##" /opt/symmetric/engines/local.properties
    sed -i -e "s#^http.basic.auth.password=<CLOUD_SYMMETRIC_BASIC_AUTH_PASSWORD>##" /opt/symmetric/engines/local.properties
fi

# Create symmetric tables
# (We need to manually create them to be able to insert our own rules in
# the next step even though symmetric would create them automatically on
# first run)
$SYM_DIR/bin/symadmin --engine local-${LOCATION_ID} \
    create-sym-tables

$SYM_DIR/bin/dbimport --engine local-${LOCATION_ID} \
    $SYM_DIR/samples/clean.sql

# Start the tool
$SYM_DIR/bin/sym --engine local-${LOCATION_ID}
