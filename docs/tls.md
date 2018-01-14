# TLS

To ensure patient data is transferred securely between components all communication in Iryo is done through secure channels.

All certificates are self-signed and generated with [Cloudflare's CFSSL](https://github.com/cloudflare/cfssl) tool.

There four types of certificates:

* root certficate
* server certificates (serving content)
* peer certificates (for identification of peer nodes)
* client certificates (for identification of clients)

All are created using a `Makefile` to ease generation in your development environment. Initially they will be created when you run `make` in your development environment. 

## Root certificate

To ease development the root certificate (`ca`) is added to the repository.

### Installing the root certificate

The root certificate has to be added to your computre for the developer to be able to use `curl` without the `-k (insecure)` flag or to open `https://iryo.local` in the browser without the `This page is not secure` warnings.

To install it navigate to repository's root folder and execute:

```bash
security add-trusted-cert -k $HOME/Library/Keychains/login.keychain bin/tls/ca.pem
```

## Adding a new component

1. components should be added to the `bin/tls/Makefile` under one of the three types `SERVERS`, `PEERS` or `CLIENTS`.
2. A `json` configuration file should be created in `bin/tls` (check `vault.json` as an example).
3. Run `make` to generate new certificates and to rebuild `traefik` docker image.