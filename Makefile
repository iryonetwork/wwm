.PHONY: up run stop build specs

ALL: vendorSync

clean:
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

specs:
	$(MAKE) -C specs

vendorSync:
	govendor sync

vendorUpdate:
	govendor fetch +missing
	govendor remove +unused
