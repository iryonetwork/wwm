BASIC_SERVICES = traefik postgres cloudSymmetric pgweb localMinio cloudMinio localNats natsStreamingExporter localPrometheusPushGateway localPrometheus cloudPrometheus cloudAuth localAuth localStorage cloudStorage localNats storageSync waitlist localStatusReporter cloudStatusReporter cloudDiscovery localDiscovery localSymmetric
BIN_CMD ?=

.PHONY: up run stop build

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
