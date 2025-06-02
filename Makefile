run:
	@env TZ=UTC go run *.go -c config.yml

build-ui:
	@cd ui && npm run build

build: build-ui
	@go build -o build/pb_launcher *.go

clean:
	@rm -rf pb_data

upgrade:
	@go run *.go upgrade

downgrade:
	@go run *.go downgrade

new-migrate:
	@go run scripts/migration/main.go $(filter-out $@,$(MAKECMDGOALS))