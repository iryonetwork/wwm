disable_mlock = true

listener "tcp" {
    address = "0.0.0.0:8200"
    tls_cert_file = "/certs/vault.pem"
    tls_key_file = "/certs/vault-key.pem"
}
