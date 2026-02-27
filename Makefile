# Makefile

# Variables
name ?= example_table_name
DATABASE_URL ?= postgresql://postgres:mysecretpassword@127.0.0.1:5433/postgres?sslmode=disable
JWT_KEY_DIR=keys
JWT_PRIVATE_KEY=$(JWT_KEY_DIR)/jwt_private.pem
JWT_PUBLIC_KEY=$(JWT_KEY_DIR)/jwt_public.pem

# Command to generate RSA 2048-bit keys
generate-jwt-keys:
	mkdir -p $(JWT_KEY_DIR)
	openssl genrsa -out $(JWT_PRIVATE_KEY) 2048
	openssl rsa -in $(JWT_PRIVATE_KEY) -pubout -out $(JWT_PUBLIC_KEY)
	@echo "✅ RSA keys generated in $(JWT_KEY_DIR)"

# Target: Create a new migration
migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

# Target: Run migrations
migrate:
	migrate -path db/migrations -database "${DATABASE_URL}" up

# Target: Rollback migrations
migrate-down:
	migrate -path db/migrations -database "${DATABASE_URL}" down -all

# Target: Drop all tables in the database
migrate-dev:
	make migrate-down
	make migrate

# Target: Generate code with SQLC
sqlc:
	sqlc generate

# Target: Build the Go project
build:
	go build -o bin/app main.go

# Target: Run the Go project
run:
	go run main.go

# Target: Run the Go project in development mode (with migration down first)
run-dev:
	make migrate-down
	make run

# Target: generate swagger documentation
swagger:
	swag init
	go run scripts/openapi/convert/main.go

# Target: Run all
all: sqlc build run swagger
