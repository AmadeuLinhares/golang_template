## start migration
include .env
export

DATABASE=goose mysql "${db_user}:${db_password}@tcp(${db_host}:${db_port})/${db_name}?${db_parameters}"
COMAND_MIGRATE=cd database/migrations && ${DATABASE}
COMAND_SEED=cd database/seeders && ${DATABASE}

## end 

VERSION = $(shell git branch --show-current)

help:  ## show this help
	@echo "usage: make [target]"
	@echo ""
	@egrep "^(.+)\:\ .*##\ (.+)" ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

run: ## run it will instance server 
	VERSION=$(VERSION) go run main.go

run-watch: ## run-watch it will instance server with reload
	VERSION=$(VERSION) nodemon --exec go run main.go --signal SIGTERM

test: ## runing unit tests with covarage
	go test -coverprofile=coverage.out ./...

mock: ## mock is command to generate mock using mockgen
	rm -rf ./mocks

	mockgen -source=./store/health/health.go -destination=./mocks/health_mock.go -package=mocks -mock_names=Store=MockHealthStore
	mockgen -source=./util/cache/cache.go -destination=./mocks/cache_mock.go -package=mocks

docs: ## docs is a command to generate doc with swagger
	swag init



# ------------------------------------ migration ------------------------------------

migrate-create: ## Creates new migration or seeders file with the current timestamp
	@read -p  "Do you want to create a migration or a seeder? [m/s] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Mm] ]]; then \
		@cd ./; read -p  "What is the migration name? " NAME; \
		${COMAND_MIGRATE} create $$NAME+_Table sql ;\
	else \
		if [[ $$REPLY =~ ^[Ss] ]]; then \
			@cd .; read -p  "What is the seeder name? " NAME; \
			${COMAND_SEED} create $$NAME+_Seed sql ;\
		fi \
	fi

# ------------------------------------ comands migrations ------------------------------------
migrate-status: ## Dump the migration status for the current DB
	${COMAND_MIGRATE} status

migrate-version: ## Print the current version of the database
	${COMAND_MIGRATE} version

migrate-up: ## Migrate the DB to the most recent version available
	${COMAND_MIGRATE} up

migrate-up-by-one: ## Migrate the DB up by 1
	${COMAND_MIGRATE} up-by-one

migrate-down: ## Roll back the version by 1
	@read -p  "Are you sure you want to turn down the last migration? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		${COMAND_MIGRATE} down; \
	fi

migrate-down-to: ## Roll back to a specific VERSION
	@read -p  "Are you sure you want to turn down the migration? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		@cd ./; read -p  "What is the migration version? " VERSION; \
		${COMAND_MIGRATE} down-to $$VERSION; \
	fi

migrate-reset: ## Roll back all migrations
	@read -p  "Are you sure you want to reset the migrations? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		${COMAND_MIGRATE} reset; \
	fi

migrate-redo: ## Re-run the latest migration
	@read -p  "Are you sure you want to reapply the last migration? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		${COMAND_MIGRATE} redo; \
	fi

# ------------------------------------ comands seeders ------------------------------------
seed-status: ## Dump the seed status for the current DB
	${COMAND_SEED} status

seed-version: ## Print the current version of the database
	${COMAND_SEED} version

seed-up: ## Seed the DB to the most recent version available
	${COMAND_SEED} up

seed-up-by-one: ## Seed the DB up by 1
	${COMAND_SEED} up-by-one

seed-down: ## Roll back the version by 1
	@read -p  "Are you sure you want to turn down the last migration? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		${COMAND_SEED} down; \
	fi

seed-down-to: ## Roll back to a specific VERSION
	@read -p  "Are you sure you want to turn down the seed? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		@cd ./; read -p  "What is the seed version?" VERSION; \
		${COMAND_MIGRATE} down-to $$VERSION; \
	fi

seed-reset: ## Roll back all seeders
	@read -p  "Are you sure you want to reset the migrations? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		${COMAND_SEED} reset; \
	fi

seed-redo: ## Re-run the latest migration
	@read -p  "Are you sure you want to reapply the last migration? [y/n] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy] ]]; then \
		${COMAND_SEED} redo; \
	fi
