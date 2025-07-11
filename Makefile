include .env

# POSTGRES_URL := postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=${PG_SSLMODE}

postgres:
	docker run --name postgres17 -p 5433:5432 -e POSTGRES_USER=doniback -e POSTGRES_PASSWORD=secret -d postgres:17.5-alpine3.22
	@echo "PostgreSQL container 'postgres17' started successfully."

stop-postgres:
	docker stop postgres17
	@echo "PostgreSQL container 'postgres17' stopped successfully."

start-postgres:
	docker start postgres17
	@echo "PostgreSQL container 'postgres17' started successfully."

restart-postgres: stop-postgres start-postgres
	@echo "PostgreSQL container 'postgres17' restarted successfully."

remove-postgres:
	docker rm -f postgres17
	@echo "PostgreSQL container 'postgres17' removed successfully."	



# Database management commands
createdb:
	docker exec -it postgres17 createdb --username=doniback --owner=doniback go_bank
	@echo "Database 'go_bank' created successfully."

dropdb:
	docker exec -it postgres17 dropdb --username=doniback --if-exists go_bank
	@echo "Database 'go_bank' dropped successfully."

migrate-postgres:
	@echo "Running PostgreSQL migrations..."
	migrate -path db/migrations -database "postgres://doniback:secret@localhost:5433/go_bank?sslmode=disable" up

rollback-postgres:
	@echo "Rolling back PostgreSQL migrations..."
	migrate -path db/migrations -database "postgres://doniback:secret@localhost:5433/go_bank?sslmode=disable" down 1

sqlc:
	@echo "Generating SQLC code..."
	sqlc generate

mock:
	@echo "Generating mock code..."
	go generate ./...

.PHONY: createdb dropdb postgres stop-postgres start-postgres restart-postgres remove-postgres migrate-postgres rollback-postgres sqlc mock
