-include .env
export

DB_URL := clickhouse://$(CLICKHOUSE_MIGRATION_HOST):$(CLICKHOUSE_PORT)/$(CLICKHOUSE_DB)?username=$(CLICKHOUSE_USER)&password=$(CLICKHOUSE_PASSWORD)&database=$(CLICKHOUSE_DB)&x-multi-statement=true

.PHONY: migrate-up migrate-down migrate-version

migrate-up:
	migrate -path ./migrations/clickhouse -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations/clickhouse -database "$(DB_URL)" down 1

migrate-version:
	migrate -path ./migrations/clickhouse -database "$(DB_URL)" version

tools: install-migrate

install-migrate:
	go install -tags 'clickhouse' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: install-migrate
.PHONY: migrate-up migrate-down migrate-version
