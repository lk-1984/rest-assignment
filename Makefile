include make.properties

SHELL := /bin/sh

.PHONY: help setup build lint test

help:
	@echo "Available targets:"
	@if [ "$(shell uname)" = "Darwin" ]; then \
		awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_\/-]+:.*#/ {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST); \
	elif [ "$(shell uname)" = "Linux" ]; then \
		awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[/a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST); \
	else \
		echo "Unsupported OS"; \
		exit 1; \
	fi

build: ## Build API app (inside devcontainer)
	@mkdir -p build && \
	go build -o build/api cmd/api/main.go

run: ## Run API app (inside devcontainer, but assuming Postgres in Docker compose is running without api)
	PG_HOSTNAME=${PG_HOSTNAME} \
	PG_PORT=${PG_PORT} \
	PG_DATABASE=${PG_DATABASE} \
	PG_USERNAME=${PG_USERNAME} \
	PG_PASSWORD=${PG_PASSWORD} \
	build/api

docker/build: ## Build multi-arch Docker image of the API app (outside devcontainer)
	docker buildx build \
		--build-arg GOLANG_TAG=${GOLANG_TAG} \
		--platform linux/amd64,linux/arm64 \
		--no-cache \
		-t api:latest \
		.

docker/compose/network/create: ## Create Docker network
	docker network create api-network

docker/compose/network/delete: ## Delete Docker network
	docker network rm api-network

docker/compose/up: ## Run API app on Docker compose (outside devcontainer)
	docker compose -f compose.yml up -d

docker/compose/down: ## Stop API app on Docker compose (outside devcontainer)
	docker compose -f compose.yml down

docker/compose/postgres/up: ## Run Postgres on Docker compose (outside devcontainer)
	docker compose -f compose.yml up -d postgres

docker/compose/postgres/down: ## Stop Postgres on Docker compose (outside devcontainer)
	docker compose -f compose.yml down postgres

test/unit: ## Run API app's unit tests
	go test -v ./...