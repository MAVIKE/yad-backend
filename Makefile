APP=cmd/app/main.go

docker_build:
	docker-compose build app

docker_run:
	docker-compose up app

build:
	go build -o bin/app.out $(APP)

run:
	go run $(APP)

test:
	go test -v ./tests/

e2e_test:
	go test -tags=e2e -v ./tests/

lint:
	go fmt ./...
	golangci-lint run

config:
	cp configs/config.yml.example configs/config.yml

stage_config:
	cp configs/stage-config.yml.example configs/config.yml

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
