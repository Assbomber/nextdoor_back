# TODO: Add go gingansta
dev:
	@go run .

start:
	@go build -o build/myzone .
	@./build/myzone

migrateup:
	@./migrate -database ${POSTGRES_URL} -path db/migrations up

migratedown:
	@./migrate -database ${POSTGRES_URL} -path db/migrations down

migrate_download:
	@which ./migrate || curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz