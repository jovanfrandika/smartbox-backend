include .env

MIGRATE_DB=docker compose --profile tools run migrate -database mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@mongo:27017/${DB_NAME}?authSource=admin

build-main: 
	go build -o ./bin/main ./cmd/main/main.go

build-consumer: 
	go build -o ./bin/consumer ./cmd/consumer/main.go

migrate-up:
	@read -p "Enter Up: " up && ${MIGRATE_DB} up $$up

migrate-down:
	@read -p "Enter Down: " down && ${MIGRATE_DB} down $$down

migrate-force:
	@read -p "Enter Version: " version && ${MIGRATE_DB} force $$version

migrate-create:
	docker compose --profile tools run migrate create -ext json ${name}

up:
	docker compose up -d --build

stop:
	docker compose stop

rebuild:
	docker compose build --no-cache
