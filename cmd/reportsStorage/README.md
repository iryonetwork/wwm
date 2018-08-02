# Reports storage

Reports files storage service.

## Configuration environment variables

| Environment variable     | Default value          | Description                                                                                                                         |
| ------------------------ | ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `DOMAIN_TYPE`            | `global`               | _Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components._ |
| `DOMAIN_ID`              | `*`                    | _Domain in which component is operating, normally it should be '_' for all cloud components and clinic ID for local components.\*   |
| `KEY_PATH`               | _none_, **_required_** | _Path to service's private key (PEM-formatted file)._                                                                               |
| `CERT_PATH`              | _none_, **_required_** | _Path to service's public key (PEM-formatted file)._                                                                                |
| `S3_ENDPOINT`            | `cloudMinio:9000`      | _S3 object storage endpoint._                                                                                                       |
| `S3_ACCESS_KEY`          | `cloud`                | _S3 object storage access key._                                                                                                     |
| `S3_REGION`              | `us-east-1`            | _S3 object storage region._                                                                                                         |
| `S3_SECRET`              | _none_, **_required_** | _S3 object storage secret._                                                                                                         |
| `STORAGE_ENCRYPTION_KEY` | _none_, **_required_** | _Base64-encoded storage encryption key._                                                                                            |
| `AUTH_HOST`              | `localAuth`            | _Hostname of adjacent (cloud) Auth service API._                                                                                    |
| `AUTH_PATH`              | `auth`                 | _Root pathof adjacent (cloud) Auth service API._                                                                                    |
| `SERVER_HOST`            | `0.0.0.0`              | _Hostname under which service exposes its HTTP servers._                                                                            |
| `SERVER_PORT`            | `443`                  | _Port under which service exposes its main HTTP server._                                                                            |
| `METRICS_PORT`           | `9090`                 | _Port under which service exposes its metrics HTTP server._                                                                         |
| `METRICS_NAMESPACE`      | `""`                   | _Namespace/path under which service exposes its metrics HTTP server._                                                               |
| `STATUS_PORT`            | `4433`                 | _Port under which service exposes its metrics HTTP server._                                                                         |
| `STATUS_NAMESPACE`       | `""`                   | _Namespace/path under which service exposes its status HTTP server._                                                                |
