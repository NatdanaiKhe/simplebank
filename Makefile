.PHONY: migrate-up migrate-down migrate-create mock-service

MIGRATIONS_PATH := $(CURDIR)/db/migration
DB_URL := postgres://postgres:postgres@localhost:5432/bank?sslmode=disable
SERVICE_PKG := github.com/NatdanaiKhe/simplebank/service
MOCK_DIR := ./api/mock
DB_MOCK_DIR := ./db/mock

dev:
	air

build:
	go build -o simplebank main.go

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" -verbose down 1

migrate-create: check-name
	migrate create -ext sql -dir db/migration -seq $(name)

check-name:
ifndef name
	$(error Usage: make migrate-create name=create_users_table)
endif

sqlc:
	sqlc generate

test:
	go test -v ./...

test-coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

mock-service: check-service
	mockgen \
		-package=mock \
		-destination=$(MOCK_DIR)/$(service).go \
		$(SERVICE_PKG) \
		$(service)


check-service:
ifndef service
	$(error Usage: make mock-service service=ServiceName)
endif

mock-store:
	mockgen \
		-package mockdb \
		-destination $(DB_MOCK_DIR)/store.go \
		github.com/NatdanaiKhe/simplebank/db/sqlc \
		Store

install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/air-verse/air@latest
