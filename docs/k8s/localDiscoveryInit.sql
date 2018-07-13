-- localDiscovery
CREATE ROLE localdiscoveryservice;
CREATE DATABASE localdiscovery;
GRANT ALL PRIVILEGES ON DATABASE localdiscovery TO localdiscoveryservice;

CREATE USER localsymmetric WITH PASSWORD 'symmetric';
GRANT localdiscoveryservice to localsymmetric;

CREATE USER localdiscovery WITH PASSWORD '434256y99m0ue5e46y77777769h3691v91049399';
GRANT localdiscoveryservice to localdiscovery;
