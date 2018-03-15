# Cloud Auth

Cloud authentication service.

## Configuration environment variables
Environment variable | Default value | Description
------------ | ------------- | -------------
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`BOLT_DB_FILEPATH` | `/data/cloudAuth.db` | *Path to Bolt DB file in which auhtentication data are stored.*
`SERVICES_FILEPATH` | `/serviceCertsAndPaths.yml` | *Path to YAML file listing services certificates and API paths that they are allowed to access.*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
