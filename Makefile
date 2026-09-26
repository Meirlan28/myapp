include .env
export

env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down

env-cleanup:
	@docker compose down -v

ACTION ?= up

migrate:
	@docker compose run --rm migrate \
		-path=/migrations \
		-database="$(DATABASE_URL)" \
		$(ACTION) $(N)

migrate-create:
	@docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq \
		$(NAME)

myapp-run:
	@go run cmd/myapp/main.go