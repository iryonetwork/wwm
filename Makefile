.PHONY: up run stop build specs

ALL: vendorSync

clear:
	docker-compose down
	rm -fr .bin

up: up/localAuth

up/%: build/% stop/%
	docker-compose up -d $*

run/%: build/% stop/%
	docker-compose up $*

stop/%:
	docker-compose stop $*

build/%:
	@mkdir -p .bin
	GOOS=linux GOARCH=amd64 go build -o ./.bin/$* ./cmd/$*

generate: clearGenerate
	go generate ./...

generate/%:
	go generate ./$*

clearGenerate:
	rm -f */mock/gen_*.go */*/mock/gen_*.go */*/*/mock/gen_*.go

test:
	go test ./...

test/%:
	go test ./$*

specs:
	$(MAKE) -C specs

vendorSync: vendor/vendor.json
	govendor sync

vendorUpdate:
	govendor fetch +missing
	govendor remove +unused
