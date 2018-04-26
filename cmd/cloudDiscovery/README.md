# Cloud Discovery

Cloud patient discovery service.

## Configuration environment variables

Environment variable | Default value | Description
------------ | ------------- | -------------
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`DB_USERNAME` | *none*, ***required*** | *PostgreSQL DB username.*
`DB_PASSWORD` | *none*, ***required*** | *PostgreSQL DB password.*
`AUTH_HOST` | `cloudAuth` | *Hostname of adjacent (cloud) Auth service API.*
`AUTH_PATH` | `auth` | *Root path of adjacent (cloud) Auth service API.*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`METRICS_PORT` | `9090` | *Port under which service exposes its metrics HTTP server.*
`METRICS_NAMESPACE` | `""` | *Namespace/path under which service exposes its metrics HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
`VAULT_ADDRESS` | `""` | *Address to use when connection to vault server.*
`VAULT_TOKEN` | `""` | *Token to use when connecting to vault server.*
`VAULT_DB_ROLE` | `""` | *Name of the DB role provisioned in vault server.*
`POSTGRES_HOST` | `postgres` | *Hostname on which postgres is exposed on.*
`POSTGRES_DATABASE` | `clouddiscovery` | *Postgres database to connect to.*
`POSTGRES_ROLE` | `clouddiscoveryservice` | *Postgres role to assume once connected.*
