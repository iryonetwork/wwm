#!/bin/bash

function howto {
    echo "./mockgen.sh $$PACKAGE $$INTERAFACES $$GOFILE"
    echo "Generates mocks for comma separated list of interfaces located in a go package."
    echo ""
}

# Check number of variables
if [ "$#" != "3" ]; then
    echo "Invalid number of arguments"
    echo ""
    howto
    exit 1
fi

# ensure mock folder exists
mkdir -p mock

# call the mockgen function
mockgen -destination=mock/gen_$3 -package mock  github.com/iryonetwork/wwm/$1 $2

# replace possible imports from the vendor folder
sed -e 's/github.com\/iryonetwork\/wwm\/vendor\///' -i '' mock/gen_$3
