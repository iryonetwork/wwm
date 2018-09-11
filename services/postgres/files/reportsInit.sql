-- reports
CREATE ROLE reportsservice;
CREATE DATABASE reports;

GRANT ALL PRIVILEGES ON DATABASE reports TO reportsservice;

CREATE USER batchdataexporter WITH PASSWORD 'batchdataexporter';
GRANT reportsservice to batchdataexporter;

CREATE USER reportgenerator WITH PASSWORD 'reportgenerator';
GRANT reportsservice to reportgenerator;
