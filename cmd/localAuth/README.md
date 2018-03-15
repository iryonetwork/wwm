# Local Auth

Local auhtentication service

## Configuration environment variables
Environment variable | Default value | Description
------------ | ------------- | -------------
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`AUTH_SYNC_KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file) used for auth data sync.*
`AUTH_SYNC_CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file) used for auth data sync.*
`BOLT_DB_FILEPATH` | `/data/localAuth.db` | *Path to Bolt DB file in which auhtentication data are stored.*
`CLOUD_AUTH_HOST` | `cloudAuth` | *Hostname of cloud Auth service API, used as a source for auth data sync.*
`CLOUD_AUTH_PATH` | `auth` | *Root path of cloud Auth service API, used as a source for auth data sync.*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
