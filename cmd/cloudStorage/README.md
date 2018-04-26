# Cloud Storage

Cloud file storage service.

## Configuration environment variables
Environment variable | Default value | Description
------------ | ------------- | -------------
`DOMAIN_TYPE` | `global` | *Domain in which component is operating, normally it should be 'global' for all cloud components and 'clinic' for local components.*
`DOMAIN_ID` | `*` |  *Domain in which component is operating, normally it should be '*' for all cloud components and clinic ID for local components.*
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`S3_ENDPOINT` | `cloudMinio:9000` | *S3 object storage endpoint.*
`S3_ACCESS_KEY` | `cloud` | *S3 object storage access key.*
`S3_REGION` | `us-east-1` | *S3 object storage region.*
`S3_SECRET` | *none*, ***required*** | *S3 object storage secret.*
`STORAGE_ENCRYPTION_KEY` | *none*, ***required***  | *Base64-encoded storage encryption key.*
`AUTH_HOST` | `localAuth` | *Hostname of adjacent (cloud) Auth service API.*
`AUTH_PATH` | `auth` | *Root pathof adjacent (cloud) Auth service API.*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`METRICS_PORT` | `9090` | *Port under which service exposes its metrics HTTP server.*
`METRICS_NAMESPACE` | `""` | *Namespace/path under which service exposes its metrics HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
