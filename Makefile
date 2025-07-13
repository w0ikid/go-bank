# Подключение переменных окружения
include .env
export $(shell sed 's/=.*//' .env)

# Константы
POSTGRES_URL := postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=$(PG_SSLMODE)
CONTAINER_NAME := postgres17
POSTGRES_IMAGE := postgres:17.5-alpine3.22

# PostgreSQL контейнер
postgres:
	docker run --name $(CONTAINER_NAME) -p $(PG_PORT):5432 \
		-e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASSWORD) \
		-d $(POSTGRES_IMAGE)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' started successfully."

stop-postgres:
	docker stop $(CONTAINER_NAME)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' stopped successfully."

start-postgres:
	docker start $(CONTAINER_NAME)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' started successfully."

restart-postgres: stop-postgres start-postgres
	@echo "PostgreSQL container '$(CONTAINER_NAME)' restarted successfully."

remove-postgres:
	docker rm -f $(CONTAINER_NAME)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' removed successfully."

# Управление базой данных
createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=$(PG_USER) --owner=$(PG_USER) $(PG_DATABASE)
	@echo "Database '$(PG_DATABASE)' created successfully."

dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb --username=$(PG_USER) --if-exists $(PG_DATABASE)
	@echo "Database '$(PG_DATABASE)' dropped successfully."

# Миграции
migrate-postgres:
	@echo "Running PostgreSQL migrations..."
	migrate -path db/migrations -database "$(POSTGRES_URL)" up

rollback-postgres:
	@echo "Rolling back PostgreSQL migrations..."
	migrate -path db/migrations -database "$(POSTGRES_URL)" down 1

# Генерация
sqlc:
	@echo "Generating SQLC code..."
	sqlc generate

mock:
	@echo "Generating mock code..."
	go generate ./...

.PHONY: createdb dropdb postgres stop-postgres start-postgres restart-postgres remove-postgres migrate-postgres rollback-postgres sqlc mock
