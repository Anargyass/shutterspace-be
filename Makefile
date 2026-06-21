# Shutterspace Backend — Makefile

.PHONY: run build migrate migrate-down seed test

# Load .env variables
include .env
export

run:
	go run ./cmd/server/main.go

build:
	go build -o ./bin/shutterspace ./cmd/server/main.go

migrate:
	@echo "Running migrations..."
	psql "$(DB_DSN)" -f migrations/001_create_enums.sql
	psql "$(DB_DSN)" -f migrations/002_create_users.sql
	psql "$(DB_DSN)" -f migrations/003_create_studio_types.sql
	psql "$(DB_DSN)" -f migrations/004_create_studios.sql
	psql "$(DB_DSN)" -f migrations/005_create_availability_slots.sql
	psql "$(DB_DSN)" -f migrations/006_create_bookings.sql
	psql "$(DB_DSN)" -f migrations/007_create_payments.sql
	@echo "Migrations complete."

migrate-down:
	@echo "Rolling back..."
	psql "$(DB_DSN)" -c "DROP TABLE IF EXISTS payments, bookings, availability_slots, studios, studio_types, users CASCADE;"
	psql "$(DB_DSN)" -c "DROP TYPE IF EXISTS payment_method, payment_status, booking_status, day_of_week, surabaya_area, user_role CASCADE;"
	@echo "Rollback complete."

seed:
	@echo "Seeding data..."
	psql "$(DB_DSN)" -f migrations/seed/001_seed_types.sql
	psql "$(DB_DSN)" -f migrations/seed/002_seed_studios_and_slots.sql
	psql "$(DB_DSN)" -f migrations/seed/003_seed_users.sql
	@echo "Seed complete."

test:
	go test ./...
