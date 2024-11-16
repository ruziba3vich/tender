SHELL := /bin/bash

include .env
export $(shell sed 's/=.*//' .env)

generate-env:
	@./generate_db_url.sh > .generated_env
	@echo "Generated .generated_env:"
	@cat .generated_env
	@source .generated_env

migrate-create: generate-env ### create new migration
	@source .generated_env && migrate create -ext sql -dir migrations 'insert_to_tables'
.PHONY: migrate-create

migrate-up: generate-env ### migration up
	@source .generated_env && migrate -path migrations -database "$$DB_URL" up
.PHONY: migrate-up

migrate-down: generate-env ### migration down
	@source .generated_env && migrate -path migrations -database "$$DB_URL" down
.PHONY: migrate-down

migrate-force: generate-env ### force migration to version 1
	@source .generated_env && migrate -path migrations -database "$$DB_URL" force 1
.PHONY: migrate-force

migrate-file: generate-env ### create a new migration file
	@source .generated_env && migrate create -ext sql -dir migrations/ -seq init_tables
.PHONY: migrate-file

give-permissions:
	- chmod +x additional.sh
	- chmod +x generate_db_url.sh

PORT ?= 8080

swag-gen:
	swag init -g ./internal/http/app/app.go -o ./internal/http/app/docs
kill-port:
	@echo "Checking if port $(PORT) is in use..."
	@PID=$$(lsof -ti :$(PORT) -sTCP:LISTEN); \
	if [ -n "$$PID" ]; then \
	  echo "Port $(PORT) is in use by process $$PID. Killing it..."; \
	  kill -9 $$PID; \
	else \
	  echo "Port $(PORT) is free."; \
	fi

run: kill-port
	@echo "Starting Go application..."
	@go run cmd/main.go
