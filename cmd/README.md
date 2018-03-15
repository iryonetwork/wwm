# IRYO WWM Services

## Common configuration environment variables
Listed below are environment variables used to populate common part of configuration for each command. 
Not all values are used by all commands, documentation for each command contains list of all variables used.

Environment variable | Default value | Description
------------ | ------------- | -------------
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`METRICS_PORT` | `9090` | *Port under which service exposes its metrics HTTP server.*
`METRICS_NAMESPACE` | `""` | *Namespace/path under which service exposes its metrics HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
`STORAGE_HOST` | `localAuth` | *Hostname of adjacent Storage service API.*
`STORAGE_PATH` | `auth` | *Root path of adjacent Storage service API.*
`AUTH_HOST` | `localAuth` | *Hostname of adjacent auth service API*
`AUTH_PATH` | `auth` | *Root path of adjacent auth service API*