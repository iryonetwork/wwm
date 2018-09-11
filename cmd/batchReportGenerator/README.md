# Batch Report Generator

Command for scheduled generation of CSV reports from EHR files data exported previously to DB.

## Configuration environment variables

| Environment variable              | Default value                                                     | Description                                                                                                                            |
| --------------------------------- | ----------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| `DOMAIN_TYPE`                     | `global`                                                          | _Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components._    |
| `DOMAIN_ID`                       | `*`                                                               | _Domain in which component is operating, normally it should be '_' for all cloud components and clinic ID for local components.\*      |
| `KEY_PATH`                        | _none_, **_required_**                                            | _Path to service's private key (PEM-formatted file)._                                                                                  |
| `CERT_PATH`                       | _none_, **_required_**                                            | _Path to service's public key (PEM-formatted file)._                                                                                   |
| `REPORT_SPECS_FILEPATHS`          | `assets/encountersReportSpec.json,assets/patientsReportSpec.json` | _*Path to JSON files defining CSV reports. If starts with "assets/" the file is assumed to be a bundled in the binary as an asset. *._ |
| `REPORTS_BUCKET_UUID`             | `c8220891-c582-41a3-893d-19e211985db5`                            | _*UUID of storage bucket in which reports are stored.*._                                                                               |
| `BOLT_DB_FILEPATH`                | `/data/batchDataExporter.db`                                      | _Path to Bolt DB file in which command saves file UUIDs for report types._                                                             |
| `STORAGE_HOST`                    | `cloudStorage`                                                    | _Hostname of storage API where reports CSV files are stored._                                                                          |
| `STORAGE_PATH`                    | `storage`                                                         | _Root path of storage API where reports CSV files are stored._                                                                         |  |  |
| `PROMETHEUS_PUSH_GATEWAY_ADDRESS` | `http://localPrometheusPushGateway:9091`                          | _Full address of Prometheus Push Gateway to push metrics from a single run of the command._                                            |
| `DB_USERNAME`                     | _none_, **_required_**                                            | _PostgreSQL DB username._                                                                                                              |
| `DB_PASSWORD`                     | _none_, **_required_**                                            | _PostgreSQL DB password._                                                                                                              |
| `POSTGRES_HOST`                   | `postgres`                                                        | _Hostname on which postgres is exposed on._                                                                                            |
| `POSTGRES_DATABASE`               | `reports`                                                         | _Postgres database to connect to._                                                                                                     |
| `POSTGRES_ROLE`                   | `reportsservice`                                                  | _Postgres role to assume once connected._                                                                                              |
| `DB_DETAILED_LOG`                 | `false`                                                           | _Allows to enable detailed DB statements log, otherwise only errors are printed._                                                      |

## Report specification files

Reports specification files define contents of CSV file report built from data in DB by report generator. Single specification file is a JSON file containing an object with following fields:

*   **type** specifies report type, it's used as a name of the report and defines bucket to which the report file will be uploaded.
*   **fileCategory** specifies EHR category to be considered for creating the report. Files data stored in DB need to have field _/category_ set to this value to be used for generation.
*   **groupByPatientID** is a boolean value specifying if report should treat multiple matching files with the same patient ID as source of data for single report row.
*   **columns** is an array of report columns names (_labels_) in order in which they should be included in the report.
*   **columnsSpecs** is a key-value object in which key is a column's name (_label_) and the value is column value spec.

Column value specs specify what value should be put under specific column. Column value specs objects have following properties:

*   **type** specifies field type:
    *   _fileMeta_ - used when the column is supposed to conain value taken from file's metadata;
    *   _value_;
    *   _fixedValue_;
    *   _boolean_;
    *   _quantity_;
    *   _integer_;
    *   _fixedValue_;
    *   _array_.
*   **metaField** is used only for values of type _fileMeta_. Possible values of the _metaField_ property are:
    *   _patientID_;
    *   _fileID_;
    *   _createdAt_;
    *   _updatedAt_.
*   **ehrPath** specifies path to value of the field in EHR file. It's used for all values except those of type _fileMeta_.
*   **unit** value specific for _quantity_ type values, specifies as a string unit in which the quantity was measured.
*   **properties** is an array of properties of items of an _array_ type field (_each array item has at least one property_). As such _properties_ property is valid only for _array_ type fields. Objects inside _properties_ are using the same schema as root-level fields and they themselves can contain _array_ type fields.
*   **format** specifies how multiple values that consitiute single item _array_ type field should be combined into single string (strings representing consecutive items of array are always seperated with `,`) _Format_ property follows _printf format string_ convention but only strings (`%s`) are in usage. It should contain the same number of `%s` parameters as number of fields listed in _properties_. The values of properties are inserted into final string in the same order as they are listed in _properties_.
*   **includeItems** is _array_ type specific object specifying which items of array should be included in the column value:
    *   _start_ - index of first item to be included, 0 or -1 if items starting with first should be included
    *   _end_ - index of first item NOT to be included in column value, -1 if column value should contain all items till the end.
