#!/bin/bash

VAULT_TOKEN=root
VAULT_ADDR=https://vault.iryo.local:8200

# enable database secret engine
vault secrets enable database

## LOCAL DISCOVERY

# enable secret management on localdiscovery database
vault write database/config/localdiscovery \
    plugin_name=postgresql-database-plugin \
    allowed_roles="localDiscoveryService" \
    connection_url="postgresql://root:root@postgres:5432/"

# create a role for the localDiscovery service
vault write database/roles/localDiscoveryService \
    db_name=localdiscovery \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT localdiscoveryservice TO \"{{name}}\";" \
    default_ttl="1h" \
    max_ttl="720h"

# create localDiscoveryPolicy
vault policy write localDiscoveryService ./services/vault/policies/localDiscoveryService.hcl

# create token for localDiscoveryService
vault token create -id=LOCAL-DISCOVERY-TOKEN -policy=localDiscoveryService -ttl=720h

## CLOUD DISCOVERY

# enable secret management on clouddiscovery database
vault write database/config/clouddiscovery \
    plugin_name=postgresql-database-plugin \
    allowed_roles="cloudDiscoveryService" \
    connection_url="postgresql://root:root@postgres:5432/"

# create a role for the cloudDiscovery service
vault write database/roles/cloudDiscoveryService \
    db_name=clouddiscovery \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT clouddiscoveryservice TO \"{{name}}\";" \
    default_ttl="1h" \
    max_ttl="720h"

# create cloudDiscoveryPolicy
vault policy write cloudDiscoveryService ./services/vault/policies/cloudDiscoveryService.hcl

# create token for cloudDiscoveryService
vault token create -id=CLOUD-DISCOVERY-TOKEN -policy=cloudDiscoveryService -ttl=720h

# example of reading a token; will be done by the service using the API
# vault read database/creds/localDiscoveryService

# example renew lease; will be done by the service periodically using the API
# vault lease renew database/creds/localDiscoveryService/ea75acdf-3c4f-fb12-a06d-07e7d254ea8a
