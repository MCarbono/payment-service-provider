DOCKER_COMPOSE_DIR := ./deploy

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

db_down:
	cd $(DOCKER_COMPOSE_DIR) && docker compose down 

db_up:
	cd $(DOCKER_COMPOSE_DIR) && docker compose up -d 

run-local:
	go run main.go