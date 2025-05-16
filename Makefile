run:
	@TZ=UTC go run main.go
	
clean:
	@rm -rf pb_data

upgrade:
	@go run main.go upgrade

downgrade:
	@go run main.go downgrade

new-migrate:
	@go run scripts/migration/main.go $(filter-out $@,$(MAKECMDGOALS))