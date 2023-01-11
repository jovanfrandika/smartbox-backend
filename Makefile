build-main: 
	go build -o ./bin/main ./cmd/main/main.go

build-consumer: 
	go build -o ./bin/consumer ./cmd/consumer/main.go

migrate-up:
	docker compose --profile tools run migrate

migrate-create:
	docker compose --profile tools run create-migration ${name}

up:
	docker compose up -d --build

stop:
	docker compose stop

rebuild:
	docker compose build --no-cache
