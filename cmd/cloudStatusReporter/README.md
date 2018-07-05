# Cloud Status Reporter

Service for reporting status of services and internet in general.

## Configuration environment variables
Environment variable | Default value | Description
------------ | ------------- | -------------
`DOMAIN_TYPE` | `global` | *Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components.*
`DOMAIN_ID` | `*` |  *Domain in which component is operating, normally it should be '*' for all cloud components and clinic ID for local components.*
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`COMPONENTS_FILEPATH` | `/components.yml` | *Path to YAML file listing components to be used for status reporting together with configuration for their status polling.*
`DEFAULT_TIMEOUT` | `1s` | *Default timeout for component status call, value used if not configured per component.*
`COUNT_THRESHOLD` | `3` | *Numnber of times given status has to be returned to be regarded as current by status reporter, value used if not configured per component.* **Note:** If status changes inconsitently, the most frequent one is regarded as current status by status reporter.
`INTERVAL` | `3s` | *Interval between status calls made by status reporter to components, value used if not configured per component.*
`STATUS_VALIDITY` | `20s` | *Time duration for which indiviudal status response from component is deemed valid and used to determine current component status, value used if not configured per component.*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`METRICS_PORT` | `9090` | *Port under which service exposes its metrics HTTP server.*
`METRICS_NAMESPACE` | `""` | *Namespace/path under which service exposes its metrics HTTP server.*
