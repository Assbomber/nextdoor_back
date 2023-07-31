# TODO: Add go gingansta
dev:
	@go run cmd/myzone/main.go

start:
	@go build -o build/myzone ./cmd/myzone
	@./build/myzone

doc:
	@which swag || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g cmd/myzone/main.go

sqlc:
	@docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

migrateup:
	@./migrate -database ${POSTGRES_URL} -path migrations up

migratedown:
	@./migrate -database ${POSTGRES_URL} -path migrations down

migrate_download:
	@which ./migrate || curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz