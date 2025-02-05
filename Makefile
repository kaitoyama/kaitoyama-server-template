.PHONY: build run test clean db-up db-down migrate generate-sqlc

# Build the application
build:
	go build -o main .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f main
	rm -rf tmp

# Start docker containers
db-up:
	docker compose up -d

# Stop docker containers
db-down:
	docker compose down

# Apply database migrations using sqldef
migrate-dry-run:
	mysqldef -h 127.0.0.1 -P 3306 -u user -ppassword template_db --dry-run < db/schema.sql

migrate:
	mysqldef -h 127.0.0.1 -P 3306 -u user -ppassword template_db < db/schema.sql

# Generate SQL code using sqlc
generate-sqlc:
	sqlc generate

# Wait for database to be ready
wait-db:
	@echo "Waiting for database to be ready..."
	@until mariadb -h127.0.0.1 -P3306 -uuser -ppassword -e "SELECT 1" >/dev/null 2>&1; do \
		sleep 1; \
	done

# Initialize development environment
init: db-up wait-db migrate generate-sqlc generate-oapi

# Stop and remove all containers
clean-docker:
	docker compose down -v

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

.PROXY: oapi
oapi: generate-server generate-models ## Generate the code from the openapi.yaml file

.PHONY: generate-server
generate-server: ## Generate the server code
	cd tools && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=server.cfg.yml ../openapi/swagger.yml

.PHONY: generate-models
generate-models: ## Generate the models code
	cd tools && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=models.cfg.yml ../openapi/swagger.yml
