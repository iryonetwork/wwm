lc = $(subst A,a,$(subst B,b,$(subst C,c,$(subst D,d,$(subst E,e,$(subst F,f,$(subst G,g,$(subst H,h,$(subst I,i,$(subst J,j,$(subst K,k,$(subst L,l,$(subst M,m,$(subst N,n,$(subst O,o,$(subst P,p,$(subst Q,q,$(subst R,r,$(subst S,s,$(subst T,t,$(subst U,u,$(subst V,v,$(subst W,w,$(subst X,x,$(subst Y,y,$(subst Z,z,$1))))))))))))))))))))))))))
BASIC_SERVICES = traefik postgres cloudSymmetric pgweb localMinio cloudMinio localNats natsStreamingExporter localPrometheusPushGateway localPrometheus cloudPrometheus cloudAuth localAuth localStorage cloudStorage localNats storageSync waitlist localStatusReporter cloudDiscovery localDiscovery localSymmetric
BIN_CMD ?=
DOCKER_TAG ?= $(shell git rev-parse --short HEAD)
DOCKER_REGISTRY ?= localhost:5000/
COMMANDS ?= $(filter-out README.md,$(patsubst cmd/%,%,$(wildcard cmd/*)))

.PHONY: up run stop build
.PRECIOUS: .bin/%

ifeq ($(CI),)
ALL: vendorSync generate yarnInstall certs rebuildDocker
else
ALL: vendorSync generate
endif

clear: clearGenerate  ## clears artifacts
	docker-compose down
	rm -fr .bin
	rm -fr .data
	rm -fr vendor/*/
	$(MAKE) -C bin/tls clean

up: $(addprefix up/,$(BASIC_SERVICES)) ## start all basic services

up/%: .bin/% stop/% ## start a service in background
	docker-compose up -d $*

up/cloudSymmetric: stop/cloudSymmetric ensure/postgres waitFor/postgres-5432
	docker-compose up -d cloudSymmetric

up/localSymmetric: stop/localSymmetric ensure/postgres waitFor/postgres-5432,cloudSymmetric-31415
	docker-compose up -d localSymmetric

ensure/%: .bin/% ## starts the service it's not already running
	docker-compose up -d $*

run/%: .bin/% stop/% ## run a service in foreground
	docker-compose up $*

stop: ## stop all services in docker-compose
	docker-compose stop

stop/%: ## stop a service in docker-compose
	docker-compose stop $*

ifeq ($(BIN_CMD),)
.bin/%: .FORCE ## builds a specific command line app
	@mkdir -p .bin
	@if [ -a cmd/$*/main.go ]; then \
		BIN_CMD=$* $(MAKE) .bin/$*; \
	fi
else
SOURCE_FILES := $(shell ./bin/listGoFiles.sh cmd/$(BIN_CMD))
.bin/%: $(SOURCE_FILES)
	GOOS=linux GOARCH=amd64 go build -o ./.bin/$* ./cmd/$*
endif

logs: ## shows docker compose logs
	docker-compose logs -f --tail=0 $*

buildFrontend/%:
	yarn --cwd frontend/$* install
	yarn --cwd frontend/$* build
	rm -fr .bin/$*Frontend
	cp -r frontend/$*/build .bin/$*Frontend

package: $(addprefix package/,$(COMMANDS)) localFrontend cloudFrontend

package/cloudAuth: INCLUDE_FILES = cmd/cloudAuth/rolesAndRules.yml
package/localFrontend: buildFrontend/local
package/localFrontend: DOCKERFILE = frontend/Dockerfile
package/localFrontend: INCLUDE_FILES = frontend/Caddyfile
package/cloudFrontend: buildFrontend/cloud
package/cloudFrontend: DOCKERFILE = frontend/Dockerfile
package/cloudFrontend: INCLUDE_FILES = frontend/Caddyfile

package/%: DOCKERFILE = Dockerfile.package
package/%: .bin/%
	echo packaging $*
	rm -fr .bin/$*-data
	mkdir -p .bin/$*-data
	cp -r .bin/$* .bin/$*-data/

	cp $(DOCKERFILE) .bin/Dockerfile
	$(if $(INCLUDE_FILES), cp -r $(INCLUDE_FILES) .bin/$*-data/,)

	docker build --build-arg BIN=$* --tag iryo/$(call lc,$*) .bin

publish: $(addprefix publish/,$(COMMANDS) localFrontend cloudFrontend)

publish/%: package/%
	docker tag iryo/$(call lc,$*) $(DOCKER_REGISTRY)iryo/$(call lc,$*):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)iryo/$(call lc,$*):$(DOCKER_TAG)

generate: clearGenerate ## run generate on all projects
	go generate -v ./...

generate/%: ## runs generate for a specific project
	go generate ./$*

clearGenerate:  ## clears artifacts created by generate
	rm -f */mock/gen_*.go */*/mock/gen_*.go */*/*/mock/gen_*.go
	rm -fr ./gen

test: test/unit ## run all tests

test/unit: ## run all unit tests
	go test -short ./...

test/unit/%: ## run unit tests for a specific project
	go test ./$*

vendorSync: vendor/vendor.json ## syncs the vendor folder to match vendor.json
	govendor sync

vendorUpdate: ## updates the vendor folder
	govendor fetch +missing +external
	govendor remove +unused

help: ## displays this message
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

certs: ## build missing certificates
	$(MAKE) -C bin/tls

watch/%:
	watchmedo shell-command -i "./.git/*;./.data/*;./.bin/*;./gen/*;gen_*.go" --recursive --ignore-directories --wait --command "$(MAKE) $*"

rebuildDocker: ## rebuilds docker container to include all certificates
	docker-compose build traefik localAuth cloudAuth waitlist

yarnInstall: ## installs requrements for frontend
	cd frontend/cloud && yarn install

waitFor/%: ## wait for a specific service to become available, use - instead of :
	docker-compose run waiter -c $(subst -,:,$*) -t 10

.FORCE:
