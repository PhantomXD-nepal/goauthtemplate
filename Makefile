APP_NAME := api
DB_URL := "mysql://root:phantom0627@tcp(localhost:3306)/testdb"

.PHONY: build run test sqlc migrate-up migrate-down tidy

build:
	@echo "🔨 building..."
	@go build -o bin/$(APP_NAME) ./cmd/api

run:
	@echo "🚀 running..."
	@go run ./cmd/api

test:
	@echo "🧪 testing..."
	@go test ./... -v

tidy:
	@go mod tidy

sqlc:
	@echo "🧬 generating sqlc code..."
	@sqlc generate -f db/sqlc.yaml

migrate-up:
	@echo "⬆️ applying migrations..."
	@migrate -path db/migrations -database $(DB_URL) up

migrate-down:
	@echo "⬇️ reverting last migration..."
	@migrate -path db/migrations -database $(DB_URL) down 1
