# Local Auth

Local auhtentication service

## Initial data

On initialization database is pulled from **cloudAuth**. Information about initial data in **cloudAuth** can be found [here](../cloudAuth/README.md).


## Configuration environment variables
Environment variable | Default value | Description
------------ | ------------- | -------------
`DOMAIN_TYPE` | `global` | *Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components.*
`DOMAIN_ID` | `*` |  *Domain in which component is operating, normally it should be '*' for all cloud components and clinic ID for local components.*
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`STORAGE_ENCRYPTION_KEY` | *none*, ***required*** | *Base64-encoded storage encryption key.*
`AUTH_SYNC_KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file) used for auth data sync.*
`AUTH_SYNC_CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file) used for auth data sync.*
`BOLT_DB_FILEPATH` | `/data/localAuth.db` | *Path to Bolt DB file in which auhtentication data are stored.*
`CLOUD_AUTH_HOST` | `cloudAuth` | *Hostname of cloud Auth service API, used as a source for auth data sync.*
`CLOUD_AUTH_PATH` | `auth` | *Root path of cloud Auth service API, used as a source for auth data sync.*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
