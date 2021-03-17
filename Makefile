run:
	go run cmd/app/main.go

lint:
	go fmt ./...
	golangci-lint run

config:
	cp configs/config.yml.example configs/config.yml

swag:
	swag init -d ./internal/delivery/http -o ./docs/swagger -g handler.go

fmt:
	go fmt ./...
