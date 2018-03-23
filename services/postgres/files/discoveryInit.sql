-- localDiscovery
CREATE ROLE localdiscoveryservice;
CREATE DATABASE localdiscovery;
GRANT ALL PRIVILEGES ON DATABASE localdiscovery TO localdiscoveryservice;

CREATE USER localsymmetric WITH PASSWORD 'symmetric';
GRANT localdiscoveryservice to localsymmetric;

-- cloudDiscovery
CREATE ROLE clouddiscoveryservice;
CREATE DATABASE clouddiscovery;
GRANT ALL PRIVILEGES ON DATABASE clouddiscovery TO clouddiscoveryservice;

CREATE USER cloudsymmetric WITH PASSWORD 'symmetric';
GRANT clouddiscoveryservice to cloudsymmetric;
