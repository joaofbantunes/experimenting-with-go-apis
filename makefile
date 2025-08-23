run-std:
	cd std && APP_ENV=development go run ./cmd/main.go

run-deps:
	docker compose up -d