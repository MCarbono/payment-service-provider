test-integration:
	go test ./tests/integration -v
test-unit:
	go test ./domain/entity  -v

test:
	go test ./... -v

.PHONY: migrate-create
migrate-create:
	@read -p "Enter the migration name: " name; \
	migrate create -ext sql -dir infra/db/migration -seq $$name

.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate