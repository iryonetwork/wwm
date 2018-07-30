# Batch Data Exporter

Command for scheduled export of data from files in storage to DB to allow for reports generation.

## Configuration environment variables

| Environment variable              | Default value                            | Description                                                                                                                         |
| --------------------------------- | ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `DOMAIN_TYPE`                     | `global`                                 | _Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components._ |
| `DOMAIN_ID`                       | `*`                                      | _Domain in which component is operating, normally it should be '_' for all cloud components and clinic ID for local components.\*   |
| `KEY_PATH`                        | _none_, **_required_**                   | _Path to service's private key (PEM-formatted file)._                                                                               |
| `CERT_PATH`                       | _none_, **_required_**                   | _Path to service's public key (PEM-formatted file)._                                                                                |
| `BUCKETS_RATE_LIMIT`              | _3_                                      | _Specifies maximum number of buckets that can be synced in parallel._                                                               |
| `DATA_ENCRYPTION_KEY`             | _none_, **_required_**                   | _Base64-encoded data encryption key for sanitizer._                                                                                 |  |
| `SANITIZER_CONFIG_FILEPATH`       | _/sanitizerConfig.json_                  | _*Path to JSON file with configuration of fields to sanitize for data sanitizer.*_                                                  |
| `BOLT_DB_FILEPATH`                | `/data/batchDataExporter.db`             | _Path to Bolt DB file in which command saves datetime of last succesful run._                                                       |
| `STORAGE_HOST`                    | `cloudStorage`                           | _Hostname of source Storage API, used as source storage for sync._                                                                  |
| `STORAGE_PATH`                    | `storage`                                | _Root path of source Storage API, used as source storage for sync._                                                                 |  |  |
| `PROMETHEUS_PUSH_GATEWAY_ADDRESS` | `http://localPrometheusPushGateway:9091` | _Full address of Prometheus Push Gateway to push metrics from a single run of the command._                                         |
| `DB_USERNAME`                     | _none_, **_required_**                   | _PostgreSQL DB username._                                                                                                           |
| `DB_PASSWORD`                     | _none_, **_required_**                   | _PostgreSQL DB password._                                                                                                           |
| `POSTGRES_HOST`                   | `postgres`                               | _Hostname on which postgres is exposed on._                                                                                         |
| `POSTGRES_DATABASE`               | `reports`                                | _Postgres database to connect to._                                                                                                  |
| `POSTGRES_ROLE`                   | `dataexportservice`                      | _Postgres role to assume once connected._                                                                                           |
| `DB_DETAILED_LOG`                 | `false`                                  | _Allows to enable detailed DB statements log, otherwise only errors are printed._                                                   |
