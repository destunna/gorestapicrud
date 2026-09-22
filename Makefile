.PHONY: up down

run:
	go run main.go

up:
	docker compose up --build -d

down:
	docker compose down -v

