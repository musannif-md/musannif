APP_NAME = musannif
ENV = 

default: run

run:
	go run cmd/$(APP_NAME)/main.go

build:
	$(ENV) go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

build-linux-amd64: ENV = GOOS=linux GOARCH=amd64
build-linux-amd64: build

clean:
	rm -rf bin/ notes/ || true
	rm musannif.db || true
	rm *.log || true
