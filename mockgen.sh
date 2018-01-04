#!/bin/bash

function howto {
    echo "./mockgen.sh $GOFILE"
    echo "Generates mocks for all public interfaces located in a go source file."
    echo ""
}

# Check number of variables
if [ "$#" != "1" ]; then
    echo "Invalid number of arguments"
    echo ""
    howto
    exit 1
fi

# ensure mock folder exists
mkdir -p mock

# call the mockgen function
mockgen -source=$1 -destination=mock/gen_$1 -package mock
