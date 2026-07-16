.PHONY: build test run fmt lint frontend-css docker-build docker-up docker-down docker-logs docker-restart clean

BACKEND_DIR := backend
BIN := $(BACKEND_DIR)/bin/server

## Local Go build/test/run (no Docker)

build:
	cd $(BACKEND_DIR) && go build -o bin/server ./cmd/server

test:
	cd $(BACKEND_DIR) && go test ./...

run:
	cd $(BACKEND_DIR) && go run ./cmd/server

fmt:
	cd $(BACKEND_DIR) && gofmt -l .

lint:
	cd $(BACKEND_DIR) && go vet ./...

frontend-css:
	cd frontend && ./scripts/build-css.sh

## Docker Compose (backend + frontend + local Postgres, all containers on this machine)

docker-build:
	docker compose build

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-restart: docker-down docker-up

clean:
	rm -f $(BIN)
	rm -f $(BACKEND_DIR)/taskflow.db
