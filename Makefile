APP=cmd/app/main.go

build:
	go build -o bin/app.out $(APP)

run:
	go run $(APP)

lint:
	go fmt ./...
	golangci-lint run

config:
	cp configs/config.yml.example configs/config.yml

swag:
	swag init --parseDependency -d ./internal/delivery/http -o ./docs/swagger -g handler.go

fmt:
	go fmt ./...

tidy:
	go mod tidy

init_db:
	. ./schema/bash/init.sh

down_db:
	. ./schema/bash/down.sh
