# Authorization

Iryo WWM authorization setup is based on [Casbin](http://casbin.org) and custom authorization data storage. Data from authorization storage is reloaded to *casbin* on every update. The model used is Role-Based Access Control with domains (*tenants*).  

## Table of contents
* [Authorization data storage](#authorization-data-storage)
* [Casbin configuration](#casbin-configuration)
* [Authorization API](#authorization-api)


## Authorization data storage

Data is currently stored in file using Bolt DB.

### Entities

#### Users
User is an object defining user's authentication data and personal metadata. 

- id (*string*)
- username (*string*)
- email (*string*)
- password (*string, stored hashed*)
- personalData
    - firstName (*string*)
    - middleName (*string*)
    - lastName (*string*)
    - dateOfBirth (*date formatted string*)
    - specialisation (*string*) 
    - nationality (*ID of code in category 'countries'*)
    - residency (*ID of code in category 'countries'*)
    - passport:
        - number (*string*)
        - issuingCountry (*ID of code in category 'countries'*)
        - expiryDate
    - licenses:
        - array of codes in category 'licenses' (*e.g. code ID of code definining specfici category of driving license*)
    - languages:
        - array of codes in category 'languages'

#### Locations
Location is an object defining real world location at which WWM clinics are run, e.g. specific camp in specific town.

- id (*string*)
- name (*string*)
- country (*ID of code in category 'countries'*)
- city (*string*)
- capacity (*integer*)
- waterSupply (*boolean*)
- electricity (*boolean*)
- manager:
    - name (*string*)
    - email (*string*)
    - phoneNumber (*string*)
- clinics:
    - array of IDs of *clinics* at the location 

#### Organizations
Organization is an object defining real world organizations that are directly running WWM clinics.

- id (*string*)
- name (*string*)
- address: 
    - addressLine1 (*string*)
    - addressLine2 (*string*)
    - postCode (*string*)
    - city (*string*)
    - country(*string*)
- legalStatus (*string*)
- serviceType (*string*)
- representative:
    - name (*string*)
    - email (*string*)
    - phoneNumber (*string*)
- primaryContact:
    - name (*string*)
    - email (*string*)
    - phoneNumber (*string*)    
- clinics:
    - array of IDs of *clinics* run by organization 

#### Clinics
Clinic is an enitty defining real world WWM clinic as a pair of location where it is operated and organization that is running it. 

- id (*string*)
- name (*string*)
- location (*string, ID of location*)
- organization (*string, ID of organization*)

#### Rules
Rule is an object defining rule subject's access to performing specific actions on specific resource.

- id (*string*)
- subject (*string, either user ID or role ID*)
- resource (*string*)
- action (*integer*)
- deny (*boolean*)

#### Roles
Role is an object defining user's property that rule's can refer to as subjects.
- id (*string*)
- name (*string*)

#### User roles
User role is an object defining user as belonging to specific *role* within specific *domain*.
- id (*string*)
- userID (*string, user ID*)
- roleID (*string, role ID*)
- domainType (*string, one of: global, organization, clinic, location, user*)
- domainID (*string, either ID of organization/clinic/location/user or \* wildcard*)

### Additional information about auth storage 

* *Clinic* is tied to *organization* and *location*, if either is removed, the clinic will be removed as well. 
* *Organization*, *location*, *clinic* and *user* ale alongside *global* domain types - each entity is a domain
   in which user can have some role (defined by existence of *user role* entity). 
   Roles held at *clinic* are automatically valid for its *location* as well. 
* Assinging to user *role* at *organization*/*clinic* can be in intuitive way described as making him part of *organization*/*clinic*. 
* *User role* entity can assign *user* any *role* in any *domain* and it can be done using *User roles* section of dashboard. Nevertheless to make basic management more intuitive there is relationship between adding user to *organization* and to *clinic*. User needs to first belong to clinic's *organization* (*have a role in organization domain*) for *clinic* to be listed in adding *user* to *clinic* form.


## Casbin configuration

Data from authorization storage is reloaded to *casbin* on every update of the authorization storage data.
* The model used is Role-Based Access Control with domains (*tenants*). 
* Subjects of access rules defined in Casbin can be either roles or individual users. 
* Roles are constrained by *domain* for which they are assigned to user. Domain is constructed by combining *domain type* and *domain ID*. 
     * Domain types:
        - global [*global for whole platform*]
        - organization 
        - clinic
        - location
        - user

    Domain ID identifies specific entity of type domain type that domain refers to. It's possible to set domain ID as *\** wildcard which makes role valid for every domain of given type.
    All the rules for domain type *global* should be assinged with wildcard as this domain type on default refers to whole platform. 
* Matchers for resources apply to allow for using wildcards *\** and *{self}* keyword (replaced with subject ID from token upon validation) in the rule's resource definition. 
* Actions supported in rules are following:
    - Read
    - Write
    - Update
    - Delete
* One rule can be applied to multiple actions (through binary matching). 
* Deny-override: both allow and deny authorization rules are supported, deny overrides the allow.

### Casbin model definition
```
[request_definition]
r = sub, dom, obj, act
[dom actual location]

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = (g(r.sub, p.sub, r.dom) ||  g(r.sub, p.sub, "*")) && (wildcardMatch(r.obj, p.obj) || wildcardMatch(r.obj, selfMatch(p.obj, r.sub))) && binaryMatch(r.act, p.act)
```

### Examples

#### Example 1
New user A is added at the same time assigned with member role to organisation XYZ; following *user roles* are created in auth storage:
   - user A's everyoneRole on domain `domainType: global; domainID: *`
   - user A's authorRole on domain `domainType: user; domainID: A_ID`
   - user A's memberRole on domain `domainType: organization; domainID: XYZ_ID`
Later user A is added as a doctor to clinic ZYX run by organization XYZ in location YXZ and how we store this relationship is creating new *user roles* in the system:
   - user A's doctorRole on domain `domainType: user; domainID: ZYX_ID`

When Casbin policy is reloaded in rection for those changes following roles are loaded:
- `A_ID, everyoneRole.ID, global.*`
- `A_ID, authorRole.ID, user.A_ID`
- `A_ID, memberRole.ID, organization.XYZ_ID`
- `A_ID, doctorRole.ID, clinic.XYZ_ID`
- `A_ID, doctorRole.ID, location.YXZ_ID` [*inferred, all user roles from clinics apply to clinic's location*]


#### Example 2
New user B is setup that should be able to access & modify everything. Following *user roles* will be created:
   - user B's everyoneRole on domain `domainType: global; domainID: *`
   - user B's authorRole on domain `domainType: user; domainID: B_ID`
   - user B's superadminRole on domain `domainType: global; domainID: *` [*manually assinged*]

When Casbin policy is reloaded in rection for those changes following roles are loaded:
- `B_ID, everyoneRole.ID, global.*`
- `B_ID, authorRole.ID, user.A_ID`
- `B_ID, superadminRole.ID, global.*`

#### Example 3
New user C is a doctor working as supervisor for all the doctors and therefore should be able to access & modify all clinical data as he's taking responsiblity for correctness of it. Following *user roles* will be created:
   - user C's everyoneRole on domain `domainType: global; domainID: *`
   - user C's authorRole on domain `domainType: user; domainID: C_ID`
   - user C's authorRole on domain `domainType: user; domainID: *` [*manually assinged*] 
   - user C's doctorRole on domain `domainType: clinic; domainID: *` [*manually assinged*] 
   - user C's doctorRole on domain `domainType: location; domainID: *` [*manually assinged*] 

When Casbin policy is reloaded in rection for those changes following roles are loaded:
- `C_ID, everyoneRole.ID, global.*`
- `C_ID, authorRole.ID, user.C_ID`
- `C_ID, authorRole.ID, user.[...for each user existing in the system]`
- `C_ID, doctorRole.ID, clinic.[...for each clinic existing in the system]`
- `C_ID, doctorRole.ID, location.[...for each location existing in the system]`


## Authorization API

Swagger documentation of auth API is available [here](./api/auth.yml). 

### Authorization data management APIs
Biggest part of authorization API is REST (CRUD) API for management of auth storage entities. Local auth has only GET (*read*) endpoints exposed. 

### Authorization APIs

#### Token endpoint
* `POST /login` endpoint authenticates user and returns a token based on `username` and `password`. User will be authenticated succesfully only if the user belongs to *domain* that given *auth* service instance is configured for. 
`CloudAuth` runs configured with `domainType: global, domainID: *` so every user in the system can login. `LocalAuth` instances are meant to be run per clinic so they are configured with `domainType: clinic, domainID: {clinicID}`. 
* `POST /renew` renews authentication token.

#### Validation endpoint
* `POST /validate` endpoint allows Iryo WWM services to checks if the user has access to perform specific actions on a specific resource within specific domain. 
* The payload of *validate* call is an array of *validation pairs*. 
* Single *validation pair* contains following fields:
    - resource (*string*)
    - domainType (*string*)
    - domainID (*string*)
    - actions (*integer*)
* Service making *validation* call specifies both resource and domain in which the check should be done. 
* The response body is an array of *validation results*. Each *validation result* contains original *validation pair* under key `query` and boolean result under key `result`. 

#### Database sync endpoint
* `GET /database` endpoint allows local instances of *auth* service to get the whole database from `CloudAuth`. Sync is performed only one way as authorization storage can be modified only using `cloudAuth` API. 

### Handling services validation

* On top of validating user's token `POST /validate` endpoint of API allows also one service to verify validity of other Iryo WWM services calls, e.g. `cloudStorage` verifies that `storageSync` call is valid. Communication between services is handled through self-signed JWT tokens. Services are provisioned with auth API by specifying list of endpoints that given certificate is valid for. 
