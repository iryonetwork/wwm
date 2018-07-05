# Cloud Auth

Cloud authentication service.

## Initial data

1. On initialization basic roles (*everyone role* & *admin role*) and rules are setup.
2. Additionally init data yaml file can be specified which will result with provisioning certain entities if they do not exists yet. 

### Development environment init data

In development environment [init data yaml file](./storageInitData.yml) is used to setup database. 
It currently contains following entities:

#### Users
**ID**: 45163e62-8fcc-436e-8b21-2db615284d1d  
**Username**: admin  
**Password**: admin  
**Email**: admin@iryo.io  
**Roles**: *admin* (3720198b-74ed-40de-a45e-8756f22e67d2)     

**ID**: 7e3ac78a-c164-4167-a3c4-be90c438217f
**Username**: user  
**Password**: user  
**Email**: user@iryo.io  
**Roles**: *everyone* (338fae76-9859-4803-8441-c5c441319cfd)   
**Organizations**: *Test organization* (12ae0bd2-c77f-47e8-a612-9adfe2f70571)  
**Clinics**: *Test clinic* (e4ebb41b-7c62-4db7-9e1c-f47058b96dd0)  

#### Locations
**Name**: Cloud  
**ID**: f7e41e48-ec79-4c78-9db6-37c0c4f78326   

**Name**: Local  
**ID**: 2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa  

#### Organizations
**Name**: Local
**ID**: a422f7f5-291b-4454-ae61-3d98c6091c3e   

#### Clinics
**Name**: Local  
**ID**: e4ebb41b-7c62-4db7-9e1c-f47058b96dd0  
**Location**: 2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa *(Local)*  
**Organization**: a422f7f5-291b-4454-ae61-3d98c6091c3e *(Test organization)*  

#### User roles
- *admin* has *everyoneRole* on *global* domain
- *admin* has *superadminRole* on *global* domain
- *user* has *everyoneRole* on *global* domain
- *user* has *memberRole* on *organization Local* domain
- *user* has *memberRole* on *clinic Local* domain

## Configuration environment variables
Environment variable | Default value | Description
-----------| ------------| -------------
`DOMAIN_TYPE` | `global` | *Domain in which component is operating, normally it should be 'cloud' for all cloud components and 'clinic' for local components.*
`DOMAIN_ID` | `*` |  *Domain in which component is operating, normally it should be '*' for all cloud components and clinic ID for local components.*
`KEY_PATH` | *none*, ***required*** | *Path to service's private key (PEM-formatted file).*
`CERT_PATH` | *none*, ***required*** | *Path to service's public key (PEM-formatted file).*
`STORAGE_ENCRYPTION_KEY` |  *none*, ***required*** | *Base64-encoded storage encryption key.*
`BOLT_DB_FILEPATH` | `/data/cloudAuth.db` | *Path to Bolt DB file in which authentication data are stored.*
`SERVICES_FILEPATH` | `/serviceCertsAndPaths.yml` | *Path to YAML file listing services certificates and API paths that they are allowed to access.*
`STORAGE_INIT_DATA_FILEPATHS` | `/rolesAndRules.yml` | *Comma-separated list of paths to YAML files containing data to be initialized in database.*
`SERVER_HOST` | `0.0.0.0` | *Hostname under which service exposes its HTTP servers.*
`SERVER_PORT` | `443` | *Port under which service exposes its main HTTP server.*
`STATUS_PORT` | `4433` | *Port under which service exposes its metrics HTTP server.*
`STATUS_NAMESPACE` | `""` | *Namespace/path under which service exposes its status HTTP server.*
