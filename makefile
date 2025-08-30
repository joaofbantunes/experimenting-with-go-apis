run-std:
	cd std && \
	OTEL_RESOURCE_ATTRIBUTES="service.name=std-api,service.version=0.0.1" \
	OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4317" \
	APP_ENV=development \
	go run ./cmd/main.go

run-deps:
	docker compose up -d