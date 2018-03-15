# Batch Storage Sync

Command for scheduled local->cloud storage sync batch recheck (performing actual files sync from local to cloud if needed). 

## Configuration environment variables
Environment variable | Default value | Description
------------ | ------------- | -------------
`DOMAIN_TYPE` | `global` | *Domain in which component is operating, normally it should be 'global' for all cloud components and 'clinic' for local components.*
`DOMAIN_ID` | `*` |  *Domain in which component is operating, normally it should be '*' for all cloud components and clinic ID for local components.*
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`STORAGE_HOST` | `cloudStorage` | *Hostname of local Storage API, used as source storage for sync.*
`STORAGE_PATH` | `storage` | *Root path of local Storage API, used as source storage for sync.*
`CLOUD_STORAGE_HOST` | `cloudStorage` | *Hostname of cloud Storage API, used as destination storage for sync.*
`CLOUD_STORAGE_PATH` | `storage` | *Root path of cloud Storage API, used as destination storage for sync.*
`BOLT_DB_FILEPATH` | `/data/batchStorageSync.db` | *Path to Bolt DB file in which command saves datetime of last succesful run.*
`PROMETHEUS_PUSH_GATEWAY_ADDRESS` | `http://localPrometheusPushGateway:9091` | *Full address of Prometheus Push Gateway to push metrics from a single run of the command.*
