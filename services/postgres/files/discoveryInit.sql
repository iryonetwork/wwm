-- localDiscovery
CREATE ROLE localdiscoveryservice;
CREATE DATABASE localdiscovery;
GRANT ALL PRIVILEGES ON DATABASE localdiscovery TO localdiscoveryservice;

CREATE USER localsymmetric WITH PASSWORD 'symmetric';
GRANT localdiscoveryservice to localsymmetric;

CREATE USER localdiscovery WITH PASSWORD 'localdiscovery';
GRANT localdiscoveryservice to localdiscovery;

-- cloudDiscovery
CREATE ROLE clouddiscoveryservice;
CREATE DATABASE clouddiscovery;
GRANT ALL PRIVILEGES ON DATABASE clouddiscovery TO clouddiscoveryservice;

CREATE USER cloudsymmetric WITH PASSWORD 'symmetric';
GRANT clouddiscoveryservice to cloudsymmetric;

CREATE USER clouddiscovery WITH PASSWORD 'clouddiscovery';
GRANT clouddiscoveryservice to clouddiscovery;
