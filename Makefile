run:
	@env TZ=UTC go run *.go -c config.yml

build-ui:
	@cd ui && rm -rf dist && npm run build-embed

build: build-ui
	@go build \
	-ldflags "-X main.commit=$(shell git rev-parse --short HEAD)" \
	-o build/pb_launcher *.go

clean:
	@rm -rf pb_data

gen-config:
	@go run *.go gen-config

print-version:
	@go run *.go version

upgrade:
	@go run *.go upgrade

downgrade:
	@go run *.go downgrade

new-migrate:
	@go run scripts/migration/main.go $(filter-out $@,$(MAKECMDGOALS))