# Batch Storage Sync

Command for scheduled local->cloud storage sync batch recheck (performing actual files sync from local to cloud if needed).

## Configuration environment variables

| Environment variable              | Default value                            | Description                                                                                                                         |
| --------------------------------- | ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `DOMAIN_TYPE`                     | `global`                                 | _Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components._ |
| `DOMAIN_ID`                       | `*`                                      | _Domain in which component is operating, normally it should be '_' for all cloud components and clinic ID for local components.\*   |
| `KEY_PATH`                        | _none_, **_required_**                   | _Path to service's private key (PEM-formatted file)._                                                                               |
| `CERT_PATH`                       | _none_, **_required_**                   | _Path to service's public key (PEM-formatted file)._                                                                                |
| `BUCKETS_RATE_LIMIT`              | _2_                                      | _Specifies maximum number of buckets that can be synced in parallel._                                                               |
| `FILES_PER_BUCKET_RATE_LIMIT`     | _3_                                      | _Specifies maximum number of files per bucket that can be synced in parallel._                                                      |
| `STORAGE_HOST`                    | `localStorage`                           | _Hostname of local Storage API, used as source storage for sync._                                                                   |
| `STORAGE_PATH`                    | `storage`                                | _Root path of local Storage API, used as source storage for sync._                                                                  |
| `CLOUD_STORAGE_HOST`              | `cloudStorage`                           | _Hostname of cloud Storage API, used as destination storage for sync._                                                              |
| `CLOUD_STORAGE_PATH`              | `storage`                                | _Root path of cloud Storage API, used as destination storage for sync._                                                             |
| `BOLT_DB_FILEPATH`                | `/data/batchStorageSync.db`              | _Path to Bolt DB file in which command saves datetime of last succesful run._                                                       |
| `PROMETHEUS_PUSH_GATEWAY_ADDRESS` | `http://localPrometheusPushGateway:9091` | _Full address of Prometheus Push Gateway to push metrics from a single run of the command._                                         |
