run:
	go run cmd\app\main.go

lint:
	golangci-lint run

swag:
	swag init -d ./internal/delivery/http -o ./docs/swagger -g handler.go
