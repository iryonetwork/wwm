# Batch Data Exporter

Command for scheduled export of data from files in storage to DB to allow for reports generation.

## Configuration environment variables

| Environment variable              | Default value                            | Description                                                                                                                                                                     |
| --------------------------------- | ---------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `DOMAIN_TYPE`                     | `global`                                 | _Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components._                                             |
| `DOMAIN_ID`                       | `*`                                      | _Domain in which component is operating, normally it should be '_' for all cloud components and clinic ID for local components.\*                                               |
| `KEY_PATH`                        | _none_, **_required_**                   | _Path to service's private key (PEM-formatted file)._                                                                                                                           |
| `CERT_PATH`                       | _none_, **_required_**                   | _Path to service's public key (PEM-formatted file)._                                                                                                                            |
| `BUCKETS_RATE_LIMIT`              | `3`                                      | _Specifies maximum number of buckets that can be synced in parallel._                                                                                                           |
| `BUCKETS_TO_SKIP`                 | `c8220891-c582-41a3-893d-19e211985db5`   | _Comma-separated list of bucket IDs from which files data are not to be exported._                                                                                              |
| `LABELS_TO_SKIP`                  | `filesCollection`                        | _Comma-separated list of labels to skip. Data from files containing any of those level is not to be exported._                                                                  |
| `EXPORT_PERIOD`                   | `336h`                                   | _Time period of data to be exported counting from last successful run. Valid units are: `ns`, `us`, `ms`, `s`, `m` and `h`. On default it's set to 336h which equals 2 weeks. _ |
| `DATA_ENCRYPTION_KEY`             | _none_, **_required_**                   | _Base64-encoded data encryption key for sanitizer._                                                                                                                             |  |
| `SANITIZER_CONFIG_FILEPATH`       | _/sanitizerConfig.json_                  | _*Path to JSON file with configuration of fields to sanitize for data sanitizer*._                                                                                              |
| `BOLT_DB_FILEPATH`                | `/data/batchDataExporter.db`             | _Path to Bolt DB file in which command saves datetime of last succesful run._                                                                                                   |
| `STORAGE_HOST`                    | `cloudStorage`                           | _Hostname of source Storage API, used as source storage for sync._                                                                                                              |
| `STORAGE_PATH`                    | `storage`                                | _Root path of source Storage API, used as source storage for sync._                                                                                                             |  |  |
| `PROMETHEUS_PUSH_GATEWAY_ADDRESS` | `http://localPrometheusPushGateway:9091` | _Full address of Prometheus Push Gateway to push metrics from a single run of the command._                                                                                     |
| `DB_USERNAME`                     | _none_, **_required_**                   | _PostgreSQL DB username._                                                                                                                                                       |
| `DB_PASSWORD`                     | _none_, **_required_**                   | _PostgreSQL DB password._                                                                                                                                                       |
| `POSTGRES_HOST`                   | `postgres`                               | _Hostname on which postgres is exposed on._                                                                                                                                     |
| `POSTGRES_DATABASE`               | `reports`                                | _Postgres database to connect to._                                                                                                                                              |
| `POSTGRES_ROLE`                   | `reportsservice`                         | _Postgres role to assume once connected._                                                                                                                                       |
| `DB_DETAILED_LOG`                 | `false`                                  | _Allows to enable detailed DB statements log, otherwise only errors are printed._                                                                                               |

## Sanitizer configuration

Files data sanitizer configuration is a configuration for encryption, removal or substring extraction of specific files' data that should not end up in their original form in the database. It's a JSON file containing array of objects defining fields of files' data to _sanitize_. Field to sanitize object may contain following properties:

*   **description** describes field to sanitize for better readability. It has no role in sanitization process.
*   **type** specifies field type:
    *   _value_;
    *   _fixedValue_;
    *   _boolean_;
    *   _quantity_;
    *   _integer_;
    *   _code_;
    *   _array_.
*   **ehrPath** specifies path to value of the field in EHR file.
*   **transformation** specifies transformation to be applied to the field:
    *   _encrypt_ - value of the field will get encrypted before exporting of file's data to DB;
    *   _remove_ - value of the field will get removed before exporting of file's data to DB;
    *   _substring_ - substring specified using _transformationParameters_ will be extracted from field value and only this substring will included in data exported to DB.
*   **transformationParameters** is transformation-specific object specifying parameters of transformation. Currently only _substring_ transformation supports _transformationParameters_:
    *   _start_ - index of first character to be included in substring, 0 or -1 if substring should start from the beginning;
    *   _end_ - index of first character NOT to be included in substring, -1 if substring should contain last character.
*   **properties** is an array of properties of items of an _array_ type field (_each array item has at least one property_). As such _properties_ property is valid only for _array_ type fields. Objects inside _properties_ are using the same schema as root-level fields and they themselves can contain _array_ type fields.
