FROM postgres:10.1-alpine

ADD ./bin/tls/ca.pem ./bin/tls/postgres.pem ./bin/tls/postgres-key.pem /certs/
RUN chown -R postgres:postgres /certs && chmod 0600 /certs/*
CMD ["postgres", "--ssl=true", "--ssl_ca_file=/certs/ca.pem", "--ssl_cert_file=/certs/postgres.pem", "--ssl_key_file=/certs/postgres-key.pem"]
