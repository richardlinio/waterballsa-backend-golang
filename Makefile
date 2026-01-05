# Makefile for WaterBall SA Backend (Golang)

.PHONY: fmt
fmt:
	docker exec backend gofumpt -l -w .

.PHONY: lint
lint:
	docker exec backend go build -o /dev/null ./...
	docker exec backend golangci-lint run

.PHONY: test
test:
	docker exec backend go test ./...

# Database Migration Commands (using goose)
# Environment variables are passed from docker-compose.yml
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=./migrations

.PHONY: migrate-status
migrate-status:
	docker exec backend sh -c 'goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "host=db port=5432 user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" status'

.PHONY: migrate-up
migrate-up:
	docker exec backend sh -c 'goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "host=db port=5432 user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" up'

.PHONY: migrate-down
migrate-down:
	docker exec backend sh -c 'goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "host=db port=5432 user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" down'

.PHONY: migrate-reset
migrate-reset:
	docker exec backend sh -c 'goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "host=db port=5432 user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" reset'
	docker exec backend sh -c 'goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "host=db port=5432 user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" up'
