#!/bin/bash

SYM_DIR=/opt/symmetric

# Create symmetric tables
# (We need to manually create them to be able to insert our own rules in
# the next step even though symmetric would create them automatically on
# first run)
$SYM_DIR/bin/symadmin --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
    create-sym-tables

# # Import the initial database structure
# # (reads the XML document and converts it to database specific SQL
# # statements)
# $SYM_DIR/bin/dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
#     --format XML --alter-case $SYM_DIR/samples/schema.xml

# Insert replication rules
# (Mind that this document also clears all preexisting configuration)
$SYM_DIR/bin/dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
    $SYM_DIR/samples/initial.sql

# Import initial data
# $SYM_DIR/bin/dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
#     --format XML --alter-case $SYM_DIR/samples/patients.xml
# $SYM_DIR/bin/dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
#     --format XML --alter-case $SYM_DIR/samples/patientLocations.xml
# $SYM_DIR/bin/dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
#     --format XML --alter-case $SYM_DIR/samples/connections.xml
# $SYM_DIR/bin/dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
#     --format XML --alter-case $SYM_DIR/samples/locations.xml

# Register the local node on the cloud node
$SYM_DIR/bin/symadmin open-registration \
    --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
    local 2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa

# Start the tool
$SYM_DIR/bin/sym --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326

