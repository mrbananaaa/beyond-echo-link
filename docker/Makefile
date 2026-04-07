include .env
export POSTGRES_USER POSTGRES_DB

DEV_COMPOSE_FILE := compose.dev.yaml
DEBUG_COMPOSE_FILE := compose.debug.yaml
TEST_COMPOSE_FILE := compose.test.yaml

.PHONY: compose-dev-up
compose-dev-up:
	@docker compose -f ${DEV_COMPOSE_FILE} up -d

.PHONY: compose-dev-up-build
compose-dev-up-build:
	@docker compose -f ${DEV_COMPOSE_FILE} up --build -d

.PHONY: compose-dev-down
compose-dev-down:
	@docker compose -f ${DEV_COMPOSE_FILE} down

.PHONY: compose-dev-drop
compose-dev-drop:
	@docker compose -f ${DEV_COMPOSE_FILE} down -v

.PHONY: compose-dev-logs
compose-dev-logs:
	@docker compose -f ${DEV_COMPOSE_FILE} logs -f

.PHONY: psql
psql:
	@docker compose -f ${DEV_COMPOSE_FILE} exec -it postgres-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)
