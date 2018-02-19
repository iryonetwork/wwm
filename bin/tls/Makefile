CDIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
DOCKER = docker run --rm -it -v $(CDIR):/certs --entrypoint='' -w /certs cfssl/cfssl
SERVERS = vault localMinio cloudMinio
PEERS = localAuth cloudAuth traefik localStorage cloudStorage
CLIENTS = localAuthSync
TYPE := test

.PRECIOUS: %.pem certs/%.pem

ALL: ca.pem $(addprefix server/,$(SERVERS)) $(addprefix peer/,$(PEERS)) $(addprefix client/,$(CLIENTS))

clean:
	rm *.csr *.pem

ca.pem:
	$(DOCKER) ./createCert.sh

%.pem: %.json
	$(DOCKER) ./createProfileCert.sh $(TYPE) $*

certs/%.pem:
	mkdir -p $(CDIR)/certs
	cp $(CDIR)/$*.pem $(CDIR)/certs/

server/%: TYPE=server
server/%: %.pem certs/%.pem ;

peer/%: TYPE=peer
peer/%: %.pem certs/%.pem ;

client/%: TYPE=client
client/%: %.pem certs/%.pem ;
