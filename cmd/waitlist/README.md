# Waitlist

Local waitlist service.

## Configuration environment variables

| Environment variable     | Default value                          | Description                                                           |
| ------------------------ | -------------------------------------- | --------------------------------------------------------------------- |
| `BOLT_DB_FILEPATH`       | `/data/waitlist.db`                    | _Path to Bolt DB file in which waitlist data are stored._             |
| `DEFAULT_LIST_ID`        | `22afd921-0630-49f4-89a8-d1ad7639ee83` | _ID of default waitlist that is ensured to always exist._             |
| `DEFAULT_LIST_NAME`      | `default`                              | _Name of default waitlist that is ensured to always exist._           |
| `KEY_PATH`               | _none_, **_required_**                 | _Path to service's private key (PEM-formatted file)._                 |
| `CERT_PATH`              | _none_, **_required_**                 | _Path to service's public key (PEM-formatted file)._                  |
| `STORAGE_ENCRYPTION_KEY` | _none_, **_required_**                 | _Base64-encoded storage encryption key._                              |
| `SERVER_HOST`            | `0.0.0.0`                              | _Hostname under which service exposes its HTTP servers._              |
| `SERVER_PORT`            | `443`                                  | _Port under which service exposes its main HTTP server._              |
| `METRICS_PORT`           | `9090`                                 | _Port under which service exposes its metrics HTTP server._           |
| `METRICS_NAMESPACE`      | `""`                                   | _Namespace/path under which service exposes its metrics HTTP server._ |
| `STATUS_PORT`            | `4433`                                 | _Port under which service exposes its metrics HTTP server._           |
| `STATUS_NAMESPACE`       | `""`                                   | _Namespace/path under which service exposes its status HTTP server._  |
| `STORAGE_HOST`           | `localAuth`                            | _Hostname of adjacent Storage service API._                           |
| `STORAGE_PATH`           | `auth`                                 | _Root path of adjacent Storage service API._                          |
| `AUTH_HOST`              | `localAuth`                            | _Hostname of adjacent auth service API_                               |
| `AUTH_PATH`              | `auth`                                 | _Root path of adjacent auth service API_                              |
