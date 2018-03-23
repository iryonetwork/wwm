#!/bin/sh

# Create discovery databases and roles
psql --file /docker-entrypoint-initdb.d/files/discoveryInit.sql

# Import localdiscovery schema and initial data
psql --dbname=localdiscovery --file=/docker-entrypoint-initdb.d/files/localDiscoverySchema.sql

# Import clouddiscovery schema and initial data
psql --dbname=clouddiscovery --file=/docker-entrypoint-initdb.d/files/cloudDiscoverySchema.sql

# Import clouddiscovery data
psql --dbname=clouddiscovery --file=/docker-entrypoint-initdb.d/files/cloudDiscoveryData.sql

# Import localdiscovery data
# psql --dbname=localdiscovery --file=/docker-entrypoint-initdb.d/files/localDiscoveryData.sql
