APP := avito-masspost
CONFIG := ./config.toml
COMPOSE := docker compose

POSTGRES_SERVICE := postgres
POSTGRES_DB := avito_masspost
POSTGRES_USER := postgres

.PHONY: up down restart logs ps wait-db psql migrate feedgen tokencheck

up:
	$(COMPOSE) up -d $(POSTGRES_SERVICE)
	$(MAKE) wait-db

down:
	$(COMPOSE) down

restart: down up

logs:
	$(COMPOSE) logs -f $(POSTGRES_SERVICE)

ps:
	$(COMPOSE) ps

wait-db:
	@until $(COMPOSE) exec -T $(POSTGRES_SERVICE) pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) >/dev/null 2>&1; do \
		echo "waiting for postgres..."; \
		sleep 1; \
	done

psql:
	$(COMPOSE) exec $(POSTGRES_SERVICE) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

migrate: wait-db
	go run ./cmd/migrate -config $(CONFIG)

feedgen:
	go run ./cmd/feedgen -config $(CONFIG)

tokencheck:
	go run ./cmd/tokencheck -config $(CONFIG)
