# Makefile for WaterBall SA Backend (Golang)

.PHONY: fmt
fmt:
	docker exec backend gofumpt -l -w .

.PHONY: lint
lint:
	docker exec backend go build -o /dev/null ./...
	docker exec backend golangci-lint run
