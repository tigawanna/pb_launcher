run:
	@TZ=UTC go run *.go
	
clean:
	@rm -rf pb_data

upgrade:
	@go run *.go upgrade

downgrade:
	@go run *.go downgrade

new-migrate:
	@go run scripts/migration/main.go $(filter-out $@,$(MAKECMDGOALS))