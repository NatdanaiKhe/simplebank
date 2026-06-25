.PHONY: migrate-up migrate-down
MIGRATIONS_PATH := $(CURDIR)/db/migration
DB_URL := postgres://postgres:postgres@localhost:5432/bank?sslmode=disable

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

migrate-create:
	@which migrate > /dev/null || (echo "Installing migrate..." && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
	@if [ -z "$(name)" ]; then echo "Error: name is required. Usage: make migrate-create name=your_migration_name"; exit 1; fi
	PATH=$$PATH:$$(go env GOPATH)/bin migrate create -ext sql -dir db/migration -seq $(name)
sqlc:
	sqlc generate

test:
	go test -v ./...
test-coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
