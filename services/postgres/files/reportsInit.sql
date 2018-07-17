-- reports
CREATE ROLE dataexportservice;
CREATE ROLE reportgenerationservice;
CREATE DATABASE reports;

GRANT ALL PRIVILEGES ON DATABASE reports TO dataexportservice;
GRANT ALL PRIVILEGES ON DATABASE reports TO reportgenerationservice;

CREATE USER batchdataexporter WITH PASSWORD 'batchdataexporter';
GRANT dataexportservice to batchdataexporter;

CREATE USER reportgenerator WITH PASSWORD 'reportgenerator';
GRANT reportgenerationservice to reportgenerator;
