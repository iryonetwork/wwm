Write-Output "Generating certificate for local CA..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl /bin/bash createCert.sh

Write-Output "Generating certificate for localMinio..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh server localMinio

Write-Output "Generating certificate for localNats..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh server localNats

Write-Output "Generating certificate for localStatusReporter..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh server localStatusReporter

Write-Output "Generating certificate for postgres..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh server postgres

Write-Output "Generating certificate for localAuth..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer localAuth

Write-Output "Generating certificate for traefik..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer traefik

Write-Output "Generating certificate for localStorage..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer localStorage

Write-Output "Generating certificate for waitlist..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer waitlist

Write-Output "Generating certificate for storageSync..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer storageSync

Write-Output "Generating certificate for localDiscovery..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer localDiscovery

Write-Output "Generating certificate for localAuthSync..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh client localAuthSync

Write-Output "Generating certificate for localNatsStreaming..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer localNatsStreaming

Write-Output "Generating certificate for localPrometheus..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer localPrometheus

Write-Output "Generating certificate for batchStorageSync..."
docker run --rm -it -v ${env:IRYO_WWM_DIR}/bin/tls:/certs --entrypoint='' -w /certs cfssl/cfssl  /bin/bash createProfileCert.sh peer batchStorageSync

md ${env:IRYO_WWM_DIR}/bin/tls/certs/ -Force
cp ${env:IRYO_WWM_DIR}/bin/tls/*.pem ${env:IRYO_WWM_DIR}/bin/tls/certs/

Import-Certificate -Filepath ${env:IRYO_WWM_DIR}/bin/tls/ca.pem -CertStoreLocation cert:\CurrentUser\Root
