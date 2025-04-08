APP_NAME=musannif

run:
	go run cmd/$(APP_NAME)/main.go

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go
