-- reports
CREATE ROLE dataexportservice;
CREATE ROLE reportgenerationservice;
CREATE DATABASE reports;

GRANT ALL PRIVILEGES ON DATABASE reports TO dataexportservice;
GRANT ALL PRIVILEGES ON DATABASE reports TO reportgenerationservice;

CREATE USER batchdataexporter WITH PASSWORD 'vdj4532ejuf270774460687e0043825466861116';
GRANT dataexportservice to batchdataexporter;

CREATE USER reportgenerator WITH PASSWORD 'mo99c7f3plq6690298807931203675165588007e';
GRANT reportgenerationservice to reportgenerator;
