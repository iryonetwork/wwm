CREATE TABLE patients (
    patient_id VARCHAR(36),
    PRIMARY KEY (patient_id)
);

ALTER TABLE patients OWNER TO localdiscoveryservice;

CREATE TABLE connections (
    patient_id VARCHAR(36),
    key VARCHAR(64),
    value VARCHAR(128),
    PRIMARY KEY (patient_id, key)
);

ALTER TABLE connections OWNER TO localdiscoveryservice;

CREATE INDEX locations_idx ON connections (value, key);

CREATE TABLE locations (
    patient_id VARCHAR(36) NOT NULL,
    location_id VARCHAR(36) NOT NULL,
    PRIMARY KEY(patient_id, location_id)
);

ALTER TABLE locations OWNER TO localdiscoveryservice;
