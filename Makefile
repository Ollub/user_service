ifneq ("$(wildcard .env)","")
    include .env
    export
endif

CUR_DIR=$(shell pwd)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: build
build:
	@echo "-- building binary"
	go build \
		-ldflags "-X main.buildHash=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o ./bin/user_service \
		./cmd

.PHONY: rocknroll
rocknroll: ## Run the service with all the dependencies and apply migrations
	docker-compose down
	docker-compose up -d
	@echo "service is running on port 8080"

.PHONY: build-migrate
build-migrate: ## Build image of goose for migration
	docker build -f ./build/migration/Dockerfile -t migrate:latest .

command=up
.PHONY: migrate
migrate: ## Run migration, use built migrate image `make build-migrate`. Default command up. Example create new sql migration: make migrate command='create add_fg sql'
	docker run --rm --network host -v $(CUR_DIR)/migrations:/migrations migrate:latest goose -v "$(GOOSE_DRIVER)" "$(GOOSE_DBSTRING)" $(command)

.PHONY: e2e
e2e: build-e2e run-e2e

.PHONY: build-e2e
build-e2e: ## Build image of goose for migration
	docker build -f ./build/e2e/Dockerfile -t e2e:latest .


.PHONY: run-e2e
run-e2e: ## Run migration, use built migrate image `make build-migrate`. Default command up. Example create new sql migration: make migrate command='create add_fg sql'
	docker run --rm --network host -v $(CUR_DIR)/e2e:/tests e2e:latest pytest -v ./
