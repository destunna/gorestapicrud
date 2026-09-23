.PHONY: up down

run:
	go run cmd/main.go

up:
	docker compose up --build -d

down:
	docker compose down -v

