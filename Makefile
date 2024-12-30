include .env

MIGRATIONS_PATH = ./cmd/migrate/migrations

hello:
	@echo "Hello, World!"

migrate-create:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))


migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATION_ADDR) up

migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATION_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

seed:
	go run cmd/migrate/seed/main.go