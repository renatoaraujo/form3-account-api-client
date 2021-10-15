dockerCompose = @docker compose

up: ## Spin up the containers
	$(dockerCompose)  up -d