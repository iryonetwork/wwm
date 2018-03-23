#!/bin/bash

function howto {
    echo "./listGoFiles.sh $$PACKAGE"
    echo "Lists all source files used by the package. It excludes all internal go packages and packages listed in the vendor folder."
    echo ""
}

# Check number of variables
if [ "$#" != "1" ]; then
    echo "Invalid number of arguments"
    echo ""
    howto
    exit 1
fi

# collect all relevant packages
ownPackage=github.com/iryonetwork/wwm
packages=$(go list -f '{{ join .Imports "\n" }}' ./$1 | grep ${ownPackage} | grep -v '/vendor/')
packages=$(echo "$1 ${packages}" | sed -e "s/github.com\/iryonetwork\/wwm\///g")

# list files used inside all the packages
for p in ${packages}; do
    files=$(go list -f "${p}/{{ join .GoFiles \"\n${p}/\" }}" ./$p)
    echo "$files"
done
