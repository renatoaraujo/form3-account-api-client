.PHONY: tests

dockerCompose = @docker compose

up: ## Spin up the containers
	$(dockerCompose)  up -d

build: ## Build the containers
	$(dockerCompose) build

reset: ## Resets the containers removing the built image
	$(dockerCompose) down --rmi local

tests: reset build ## Run the tests
	$(dockerCompose) up tests