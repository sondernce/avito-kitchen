.PHONY: up down build test clean

up:
	docker-compose up --build

down:
	docker-compose down

clean:
	docker-compose down -v

build:
	docker-compose build

test:
	cd main-service && go test ./...
	cd restaurant-service && go test ./...

lint:
	golangci-lint run ./main-service/...
	golangci-lint run ./restaurant-service/...