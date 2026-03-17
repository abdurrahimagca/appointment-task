DB_PATH := services/db/appointments.db
DB_URL := sqlite3://$(DB_PATH)
MIGRATIONS_DIR := services/db/migrations

.PHONY: tools-install migrate-new migrate-up migrate-down migrate-force sqlc-generate seed

tools-install:
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

migrate-new:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down 1

migrate-force:
	@read -p "Version: " version; \
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) force $$version

sqlc-generate:
	sqlc generate

seed:
	sqlite3 $(DB_PATH) < services/db/seed.sql
