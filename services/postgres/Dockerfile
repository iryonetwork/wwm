FROM postgres:10.1-alpine

ADD ./bin/tls/ca.pem ./bin/tls/postgres.pem ./bin/tls/postgres-key.pem /ownCerts/
RUN chown -R postgres:postgres /ownCerts && chmod 0600 /ownCerts/* && chown -R postgres:postgres /ownCerts/*
CMD ["postgres", "--ssl=true", "--ssl_ca_file=/ownCerts/ca.pem", "--ssl_cert_file=/ownCerts/postgres.pem", "--ssl_key_file=/ownCerts/postgres-key.pem"]
