# Batch Data Exporter

Command for scheduled export of data from files in storage to DB to allow for reports generation.

## Configuration environment variables

| Environment variable              | Default value                            | Description                                                                                                                         |
| --------------------------------- | ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `DOMAIN_TYPE`                     | `cloud`                                  | _Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components._ |
| `DOMAIN_ID`                       | `*`                                      | _Domain in which component is operating, normally it should be '_' for all cloud components and clinic ID for local components.\*   |
| `KEY_PATH`                        | _none_, **_required_**                   | _Path to service's private key (PEM-formatted file)._                                                                               |
| `CERT_PATH`                       | _none_, **_required_**                   | _Path to service's public key (PEM-formatted file)._                                                                                |
| `STORAGE_HOST`                    | `cloudStorage`                           | _Hostname of local Storage API, used as source storage for sync._                                                                   |
| `STORAGE_PATH`                    | `storage`                                | _Root path of local Storage API, used as source storage for sync._                                                                  |
| `BOLT_DB_FILEPATH`                | `/data/batchStorageSync.db`              | _Path to Bolt DB file in which command saves datetime of last succesful run._                                                       |
| `PROMETHEUS_PUSH_GATEWAY_ADDRESS` | `http://localPrometheusPushGateway:9091` | _Full address of Prometheus Push Gateway to push metrics from a single run of the command._                                         |
