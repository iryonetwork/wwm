BASIC_SERVICES = traefik vault localMinio cloudMinio localNats cloudAuth localAuth localStorage cloudStorage localNats storageSync

.PHONY: up run stop build specs

ifeq ($(CI),)
ALL: vendorSync generate specs yarnInstall certs rebuildDocker
else
ALL: vendorSync specs generate
endif

clear: clearGenerate clearSpecs ## clears artifacts
	docker-compose down
	rm -fr .bin
	rm -fr .data
	rm -fr vendor/*/

up: $(addprefix up/,$(BASIC_SERVICES)) ## start all basic services

up/%: build/% stop/% ## start a service in background
	docker-compose up -d $*

run/%: build/% stop/% ## run a service in foreground
	docker-compose up $*

stop: ## stop all services in docker-compose
	docker-compose stop

stop/%: ## stop a service in docker-compose
	docker-compose stop $*

build/%: ## builds a specific project
	@mkdir -p .bin
	if [ -a cmd/$*/main.go ]; then \
		GOOS=linux GOARCH=amd64 go build -o ./.bin/$* ./cmd/$* ; \
	fi

logs: ## shows docker compose logs
	docker-compose logs -f --tail=0 $*

generate: clearGenerate ## run generate on all projects
	go generate -v ./...

generate/%: ## runs generate for a specific project
	go generate ./$*

clearGenerate: ## clears artifacts created by generate
	rm -f */mock/gen_*.go */*/mock/gen_*.go */*/*/mock/gen_*.go
	rm -rf gen

test: test/unit ## run all tests

test/unit: ## run all unit tests
	go test ./...

test/unit/%: ## run unit tests for a specific project
	go test ./$*

specs: ## rebuild specs
	$(MAKE) -C specs

clearSpecs:
	$(MAKE) -C specs clear

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
	docker-compose build traefik localAuth cloudAuth

yarnInstall: ## installs requrements for frontend
	cd frontend/cloud && yarn install
