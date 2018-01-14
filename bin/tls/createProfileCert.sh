#!/bin/bash

cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json \
    -profile=$1 $2.json | cfssljson -bare $2
