# Authorization

Iryo WWM authorization setup is based on [Casbin](http://casbin.org) and custom authorization data storage. Data from authorization storage is reloaded to _casbin_ on every update. The model used is Role-Based Access Control with domains (_tenants_).

## Table of contents

* [Authorization data storage](#authorization-data-storage)
* [Casbin configuration](#casbin-configuration)
* [Authorization API](#authorization-api)

## Authorization data storage

Data is currently stored in file using Bolt DB.

### Entities

#### Users

User is an object defining user's authentication data and personal metadata.

* id (_string_)
* username (_string_)
* email (_string_)
* password (_string, stored hashed_)
* personalData
  * firstName (_string_)
  * middleName (_string_)
  * lastName (_string_)
  * dateOfBirth (_date formatted string_)
  * specialisation (_string_)
  * nationality (_ID of code in category 'countries'_)
  * residency (_ID of code in category 'countries'_)
  * passport:
    * number (_string_)
    * issuingCountry (_ID of code in category 'countries'_)
    * expiryDate
  * licenses:
    * array of codes in category 'licenses' (_e.g. code ID of code definining specfici category of driving license_)
  * languages:
    * array of codes in category 'languages'

#### Locations

Location is an object defining real world location at which WWM clinics are run, e.g. specific camp in specific town.

* id (_string_)
* name (_string_)
* country (_ID of code in category 'countries'_)
* city (_string_)
* capacity (_integer_)
* waterSupply (_boolean_)
* electricity (_boolean_)
* manager:
  * name (_string_)
  * email (_string_)
  * phoneNumber (_string_)
* clinics:
  * array of IDs of _clinics_ at the location

#### Organizations

Organization is an object defining real world organizations that are directly running WWM clinics.

* id (_string_)
* name (_string_)
* address:
  * addressLine1 (_string_)
  * addressLine2 (_string_)
  * postCode (_string_)
  * city (_string_)
  * country(_string_)
* legalStatus (_string_)
* serviceType (_string_)
* representative:
  * name (_string_)
  * email (_string_)
  * phoneNumber (_string_)
* primaryContact:
  * name (_string_)
  * email (_string_)
  * phoneNumber (_string_)
* clinics:
  * array of IDs of _clinics_ run by organization

#### Clinics

Clinic is an enitty defining real world WWM clinic as a pair of location where it is operated and organization that is running it.

* id (_string_)
* name (_string_)
* location (_string, ID of location_)
* organization (_string, ID of organization_)

#### Rules

Rule is an object defining rule subject's access to performing specific actions on specific resource.

* id (_string_)
* subject (_string, either user ID or role ID_)
* resource (_string_)
* action (_integer_)
* deny (_boolean_)

#### Roles

Role is an object defining user's property that rule's can refer to as subjects.

* id (_string_)
* name (_string_)

#### User roles

User role is an object defining user as belonging to specific _role_ within specific _domain_.

* id (_string_)
* userID (_string, user ID_)
* roleID (_string, role ID_)
* domainType (_string, one of: global, cloud, organization, clinic, location, user_)
* domainID (_string, either ID of organization/clinic/location/user or \* wildcard_)

### Additional information about auth storage

* _Clinic_ is tied to _organization_ and _location_, if either is removed, the clinic will be removed as well.
* _Organization_, _location_, _clinic_ and _user_ are, alongside _global_, domain types - each entity of a certain type is a domain
  in which user can have some role (defined by existence of _user role_ entity).
  Roles held at _clinic_ are automatically valid for its _location_ as well.
* Assinging to user _role_ at _organization_/_clinic_ can be in intuitive way described as making him part of _organization_/_clinic_.
* _User role_ entity can assign _user_ any _role_ in any _domain_ and it can be done using _User roles_ section of dashboard. Nevertheless to make basic management more intuitive there is relationship between adding user to _organization_ and to _clinic_. User needs to first belong to clinic's _organization_ (_have a role in organization domain_) for _clinic_ to be listed in adding _user_ to _clinic_ form.

## Casbin configuration

Data from authorization storage is reloaded to _casbin_ on every update of the authorization storage data.

* The model used is Role-Based Access Control with domains (_tenants_).
* Subjects of access rules defined in Casbin can be either roles or individual users.
* Roles are constrained by _domain_ for which they are assigned to user. Domain is constructed by combining _domain type_ and _domain ID_.

  * Domain types:
    * global [*global for whole platform*]
    * cloud [*for cloud part of platform*]
    * organization
    * clinic
    * location
    * user

  Domain ID identifies specific entity of type domain type that domain refers to. It's possible to set domain ID as _\*_ wildcard which makes role valid for every domain of given type.
  All the rules for domain types _global_ and _cloud_ should be assinged with wildcard as this domain type on default refers to whole system.

* Matchers for resources allow for using wildcards _\*_ and _{self}_ keyword (replaced with subject ID from token upon validation) in the ACL rule's resource definition.
* Actions supported in rules are following:
  * Read (_binary 1_)
  * Write (_binary 2_)
  * Delete (_binary 4_)
  * Update (_binary 8_)
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

New user A is added and at the same time assigned with a `member role` to organisation XYZ; following _user roles_ are created in auth storage:

* user A's everyoneRole on domain `domainType: global; domainID: *` [*automatically when user is created*]
* user A's memberRole on domain `domainType: cloud; domainID: *` [*automatically when user is created*]
* user A's authorRole on domain `domainType: user; domainID: A_ID`
* user A's memberRole on domain `domainType: organization; domainID: XYZ_ID`
  Later user A is added as a doctor to clinic ZYX run by organization XYZ in location YXZ. The relationship is stored by creating new _user roles_ in the system:
* user A's doctorRole on domain `domainType: clinic; domainID: ZYX_ID`

When Casbin policy is reloaded in rection for those changes following roles are loaded:

* `A_ID, everyoneRole.ID, *` [*global*]
* `A_ID, memberRole.ID, cloud.*`
* `A_ID, authorRole.ID, user.A_ID`
* `A_ID, memberRole.ID, organization.XYZ_ID`
* `A_ID, doctorRole.ID, clinic.XYZ_ID`
* `A_ID, doctorRole.ID, location.YXZ_ID` [*inferred, all user roles from clinics apply to clinic's location*]

#### Example 2

New user B is set up. User B should be able to access & modify everything. Following _user roles_ will be created:

* user B's everyoneRole on domain `domainType: global; domainID: *` [*automatically when user is created*]
* user B's memberRole on domain `domainType: cloud; domainID: *` [*automatically when user is created*]
* user B's authorRole on domain `domainType: user; domainID: B_ID`
* user B's superadminRole on domain `domainType: global; domainID: *` [*manually assinged*]

When Casbin policy is reloaded in rection for those changes following roles are loaded:

* `B_ID, everyoneRole.ID, *`
* `B_ID, cloudRole.ID, cloud.*`
* `B_ID, authorRole.ID, user.A_ID`
* `B_ID, superadminRole.ID, global.*`

#### Example 3

New user C is a doctor working as a supervisor for all the doctors and therefore should be able to access & modify all clinical data as s/he's taking responsiblity for the data being correct. Following _user roles_ will be created:

* user C's everyoneRole on domain `domainType: global; domainID: *` [*automatically when user is created*]
* user C's memberRole on domain `domainType: cloud; domainID: *` [*automatically when user is created*]
* user C's authorRole on domain `domainType: user; domainID: C_ID`
* user C's authorRole on domain `domainType: user; domainID: *` [*manually assinged*]
* user C's doctorRole on domain `domainType: clinic; domainID: *` [*manually assinged*]
* user C's doctorRole on domain `domainType: location; domainID: *` [*manually assinged*]

When Casbin policy is reloaded following roles are loaded:

* `C_ID, everyoneRole.ID, global.*`
* `C_ID, memberRole.ID, cloud.*`
* `C_ID, authorRole.ID, user.C_ID`
* `C_ID, authorRole.ID, user.[...for each user existing in the system]`
* `C_ID, doctorRole.ID, clinic.[...for each clinic existing in the system]`
* `C_ID, doctorRole.ID, location.[...for each location existing in the system]`

## Authorization API

Swagger documentation of auth API is available [here](./api/auth.yml).

### Authorization data management APIs

Biggest part of authorization API is REST (CRUD) API for management of auth storage entities. Local auth has only GET (_read_) endpoints exposed.

### Authorization APIs

#### Token endpoint

* `POST /login` endpoint authenticates user and returns a token based on `username` and `password`. User will be authenticated succesfully only if the user belongs to _domain_ that given _auth_ service instance is configured for.
  `CloudAuth` runs configured with `domainType: cloud, domainID: *` so every user in the system can login. `LocalAuth` instances are meant to be run per clinic so they are configured with `domainType: clinic, domainID: {clinicID}`.
* `POST /renew` renews authentication token.

#### Validation endpoint

* `POST /validate` endpoint allows Iryo WWM services to checks if the user has access to perform specific actions on a specific resource within specific domain.
* The payload of _validate_ call is an array of _validation pairs_.
* Single _validation pair_ contains following fields:
  * resource (_string_)
  * domainType (_string_)
  * domainID (_string_)
  * actions (_integer_)
* Service making _validation_ call specifies both resource and domain in which the check should be done.
* The response body is an array of _validation results_. Each _validation result_ contains original _validation pair_ under key `query` and boolean result under key `result`.

#### Database sync endpoint

* `GET /database` endpoint allows local instances of _auth_ service to get the whole database from `CloudAuth`. Sync is performed only one way as authorization storage can be modified only using `cloudAuth` API.

### Handling services validation

* On top of validating user's token `POST /validate` endpoint of API allows also one service to verify validity of other Iryo WWM services calls, e.g. `cloudStorage` verifies that `storageSync` call is valid. Communication between services is handled through self-signed JWT tokens. Services are provisioned with auth API by specifying list of endpoints that given certificate is valid for.
