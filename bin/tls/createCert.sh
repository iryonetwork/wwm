#!/bin/bash
cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
