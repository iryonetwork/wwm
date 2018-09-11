CREATE TABLE files (
    file_id VARCHAR(36),
    version VARCHAR(36),
    patient_id VARCHAR(36),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    data jsonb,
    PRIMARY KEY (file_id)
);

ALTER TABLE files OWNER TO reportsservice;

CREATE INDEX created_idx ON files (created_at ASC);
CREATE INDEX patient_idx ON files (patient_id);

CREATE INDEX category_idx ON files ((data->>'/category'));
CREATE INDEX clinic_idx ON files ((data->>'/context/health_care_facility|identifier'));
CREATE INDEX author_idx ON files ((data->>'/composer|identifier'));

CREATE INDEX data_idx ON files USING gin (data);
