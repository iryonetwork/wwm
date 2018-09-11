-- reports
CREATE ROLE reportsservice;
CREATE DATABASE reports;

GRANT ALL PRIVILEGES ON DATABASE reports TO reportsservice;

CREATE USER batchdataexporter WITH PASSWORD 'vdj4532ejuf270774460687e0043825466861116';
GRANT reportsservice to batchdataexporter;

CREATE USER reportgenerator WITH PASSWORD 'mo99c7f3plq6690298807931203675165588007e';
GRANT reportsservice to reportgenerator;
