# DOCKER_COMPOSE_DIR := ./deploy

# db_down:
# 	cd $(DOCKER_COMPOSE_DIR) && docker compose down 

# db_up:
# 	cd $(DOCKER_COMPOSE_DIR) && docker compose up -d 

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
	docker compose down 

db_up:
	docker compose up -d 

run-local:
	go run main.go

build:
	docker-compose -f docker-compose.yml build

run_prod:
	docker-compose -f docker-compose.yml up -d