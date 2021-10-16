.PHONY: tests

dockerCompose = @docker compose

up: ## Spin up the containers
	$(dockerCompose)  up -d

build: ## Build the containers
	$(dockerCompose) build

reset: ## Resets the containers removing the images
	$(dockerCompose) down --rmi all

tests: reset build ## Run the tests
	$(dockerCompose) up tests