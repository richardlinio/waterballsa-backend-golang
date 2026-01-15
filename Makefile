# Makefile for WaterBall SA Backend (Golang)

# Docker command variables
DOCKER_EXEC=docker exec backend

# Database Migration Commands (using goose)
# Environment variables are passed from docker-compose.yml
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=./migrations
GOOSE_DB_DSN=host=db port=5432 user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable
GOOSE_CMD=goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DB_DSN)"

.PHONY: fmt
fmt:
	$(DOCKER_EXEC) gofumpt -l -w .

.PHONY: lint
lint:
	$(DOCKER_EXEC) go build -o /dev/null ./...
	$(DOCKER_EXEC) golangci-lint run

.PHONY: test
test:
	$(DOCKER_EXEC) go test ./...

.PHONY: build
build:
	$(DOCKER_EXEC) go build -o /dev/null ./...

.PHONY: tidy
tidy:
	$(DOCKER_EXEC) go mod tidy

.PHONY: sqlc
sqlc:
	$(DOCKER_EXEC) sqlc generate

.PHONY: migrate-status
migrate-status:
	$(DOCKER_EXEC) sh -c '$(GOOSE_CMD) status'

.PHONY: migrate-up
migrate-up:
	$(DOCKER_EXEC) sh -c '$(GOOSE_CMD) up'

.PHONY: migrate-down
migrate-down:
	$(DOCKER_EXEC) sh -c '$(GOOSE_CMD) down'

.PHONY: migrate-reset
migrate-reset:
	$(DOCKER_EXEC) sh -c '$(GOOSE_CMD) reset'
	$(DOCKER_EXEC) sh -c '$(GOOSE_CMD) up'

.PHONY: test-bdd-isa
test-bdd-isa:
	$(DOCKER_EXEC) go test -v ./tests/bdd -run TestFeatures
