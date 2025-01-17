# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s | tr A-Z a-z)
export ARCH := $(shell uname -m)



# --- Tooling & Variables ----------------------------------------------------------------
include ./scripts/make/tools.Makefile
include ./scripts/make/help.Makefile

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

up: dev-env dev-air             ## Startup / Spinup Docker Compose and air
down: docker-stop               ## Stop Docker
destroy: docker-teardown clean-artifacts  ## Teardown (removes volumes, tmp files, etc...)

install-deps: migrate air gotestsum tparse mockery ## Install Development Dependencies (localy).
deps: $(MIGRATE) $(AIR) $(GOTESTSUM) $(TPARSE) $(MOCKERY) $(GOLANGCI) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

dev-env: ## Bootstrap Environment (with a Docker-Compose help).
	@ docker-compose up -d --build db loki promtail grafana redis

dev-env-test: dev-env ## Run application (within a Docker-Compose help)
	@ $(MAKE) image-build
	docker-compose up web

dev-air: $(AIR) ## Starts AIR ( Continuous Development app).
	air

docker-stop:
	@ docker-compose down

docker-teardown:
	@ docker-compose down --remove-orphans -v

# ~~~ Docker Build ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.ONESHELL:
image-build:
	@ echo "Docker Build"
	@ DOCKER_BUILDKIT=0 docker build \
		--file Dockerfile \
		--tag tuskar-go-fiber \
		--build-arg SERVER_PORT=$(SERVER_PORT) \
			.

# # ~~~ Database Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

POSTGRE_DSN := "postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable"

migrate-up: $(MIGRATE) ## Apply all (or N up) migrations.
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate  -database $(POSTGRE_DSN) -path=scripts/migrations up ${NN}

.PHONY: migrate-down
migrate-down: $(MIGRATE) ## Apply all (or N down) migrations.
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate  -database $(POSTGRE_DSN) -path=scripts/migrations down ${NN}

.PHONY: migrate-dirty
migrate-dirty: $(MIGRATE) ## Clean a dirty migration by force
	@ read -p "Fixed migration error and continue the migration process on version: " N; \
	migrate -database $(POSTGRE_DSN) -path=scripts/migrations force $${N}

.PHONY: migrate-create
migrate-create: $(MIGRATE) ## Create a set of up/down migrations with a specified name.
	@ read -p "Please provide name for the migration: " Name; \
	migrate create -ext sql -dir scripts/migrations -seq $${Name}
