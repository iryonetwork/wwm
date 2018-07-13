-- cloudDiscovery
CREATE ROLE clouddiscoveryservice;
CREATE DATABASE clouddiscovery;
GRANT ALL PRIVILEGES ON DATABASE clouddiscovery TO clouddiscoveryservice;

CREATE USER cloudsymmetric WITH PASSWORD 'symmetric';
GRANT clouddiscoveryservice to cloudsymmetric;

CREATE USER clouddiscovery WITH PASSWORD 'i62r16x64973ue33m3m4x042339092674140213w';
GRANT clouddiscoveryservice to clouddiscovery;
