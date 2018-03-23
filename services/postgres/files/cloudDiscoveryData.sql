-- patients
INSERT INTO patients (patient_id) VALUES ('554ec907-9f86-4abf-a1fb-7ec211c22611');
INSERT INTO patients (patient_id) VALUES ('ea925eba-98c1-4708-a6b6-c337dfa45900');

-- patient connections
INSERT INTO connections (patient_id, key, value) VALUES ('554ec907-9f86-4abf-a1fb-7ec211c22611', 'name', 'Local User');
INSERT INTO connections (patient_id, key, value) VALUES ('ea925eba-98c1-4708-a6b6-c337dfa45900', 'name', 'Cloud User');

-- patient locations
INSERT INTO locations (patient_id, location_id) VALUES ('554ec907-9f86-4abf-a1fb-7ec211c22611', '2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa');
INSERT INTO locations (patient_id, location_id) VALUES ('ea925eba-98c1-4708-a6b6-c337dfa45900', 'f7e41e48-ec79-4c78-9db6-37c0c4f78326');
