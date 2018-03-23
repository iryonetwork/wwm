#!/bin/bash

SYM_DIR=/opt/symmetric

# Create symmetric tables
# (We need to manually create them to be able to insert our own rules in
# the next step even though symmetric would create them automatically on
# first run)
$SYM_DIR/bin/symadmin --engine local-2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa \
    create-sym-tables

# Register the local node on the cloud node
$SYM_DIR/bin/symadmin open-registration \
    --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
    local 2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa

# Trigger the full reload of the local node
$SYM_DIR/bin/symadmin reload-node \
    --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 \
    2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa

# Start the tool
$SYM_DIR/bin/sym --engine local-2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa
