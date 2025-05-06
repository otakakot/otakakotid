SHELL := /bin/bash
include .env
export
export APP_NAME := $(basename $(notdir $(shell pwd)))

.PHONY: gen
gen: ## Generate code.
	@go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types -package api api/openapi.yaml > pkg/api/types.gen.go
	@find pkg/schema -type f -not -name "*.sql" -exec rm -rf {} \;
	@go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate
	@go mod tidy

.PHONY: module
module: ## go modules
	@go get -u -t ./...
	@go mod tidy

.PHONY: oidc-dev
oidc-dev: ## oidc dev
	@make -C worker/oidc dev

.PHONY: authn-dev
authn-dev: ## authn dev
	@make -C worker/authn dev

.PHONY: oidc-deploy
oidc-deploy: ## oidc deploy
	@make -C worker/oidc deploy

.PHONY: authn-deploy
authn-deploy: ## authn deploy
	@make -C worker/authn deploy

.PHONY: migrate
migrate: ## migrate
	@wrangler d1 execute otakakot --remote --file=./schema/schema.sql
