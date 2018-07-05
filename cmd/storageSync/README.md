# Storage Sync

Service consuming sync messages from local Storage published via NATS streaming. It continuously syncs local Storage to cloud Storage.

## Configuration environment variables
Environment variable | Default value | Description
------------ | ------------- | -------------
`DOMAIN_TYPE` | `global` | *Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components.*
`DOMAIN_ID` | `*` |  *Domain in which component is operating, normally it should be '*' for all cloud components and clinic ID for local components.*
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`METRICS_PORT` | `9090` | *Port under which service exposes its metrics HTTP server.*
`METRICS_NAMESPACE` | `""` | *Namespace/path under which service exposes its metrics HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
`STORAGE_HOST` | `cloudStorage` | *Hostname of local Storage service API, used as source storage for sync.*
`STORAGE_PATH` | `storage` | *Root path of local Storage service API, used as source storage for sync.*
`CLOUD_STORAGE_HOST` | `cloudStorage` | *Hostname of cloud Storage service API, used as destination storage for sync.*
`CLOUD_STORAGE_PATH` | `storage` | *Root path of cloud Storage service API, used as destination storage for sync.*
`NATS_ADDR` | `localNats:4242` | *NATS server address.*
`NATS_USERNAME` | `nats` | *Username used to connect to NATS.*
`NATS_SECRET` | *none*, ***required*** | *Secret used to connect to NATS.*
`NATS_CONN_RETRIES` | `10` | *Number of attempts to connect to NATS.*
`NATS_CONN_WAIT` | `500ms` | *Initial wait time before reattempting to connect to NATS after failed attempt.*
`NATS_CONN_WAIT_FACTOR` | `3.0` | *Factor by which wait time increases after each consecutive failed retry.*
`NATS_CLUSTER_ID` | `localNats` | *NATS Streaming cluster ID*
`NATS_CLIENT_ID` | `storageSync` | *NATS Streaming client ID*