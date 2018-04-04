CREATE TABLE patients (
    patient_id VARCHAR(36),
    PRIMARY KEY (patient_id)
);

ALTER TABLE patients OWNER TO clouddiscoveryservice;

CREATE TABLE connections (
    patient_id VARCHAR(36),
    key VARCHAR(64),
    value VARCHAR(128),
    PRIMARY KEY (patient_id, key)
);

ALTER TABLE connections OWNER TO clouddiscoveryservice;

CREATE INDEX locations_idx ON connections (value, key);

CREATE TABLE locations (
    patient_id VARCHAR(36) NOT NULL,
    location_id VARCHAR(36) NOT NULL,
    PRIMARY KEY(patient_id, location_id)
);

ALTER TABLE locations OWNER TO clouddiscoveryservice;

CREATE TABLE codes (
    category_id VARCHAR(64) NOT NULL,
    code_id VARCHAR(64) NOT NULL,
    parent_id VARCHAR(64),
    PRIMARY KEY (category_id, code_id)
);

CREATE INDEX codes_parent_idx ON codes (parent_id);

ALTER TABLE codes OWNER TO clouddiscoveryservice;

CREATE TABLE code_titles (
    category_id VARCHAR(64) NOT NULL,
    code_id VARCHAR(64) NOT NULL,
    locale VARCHAR(64) NOT NULL,
    title VARCHAR(255),
    PRIMARY KEY (category_id, code_id, locale),
    FOREIGN KEY (category_id, code_id) REFERENCES codes (category_id, code_id) ON DELETE CASCADE
);

CREATE INDEX code_titles_title_idx ON code_titles (title);

ALTER TABLE code_titles OWNER TO clouddiscoveryservice;
