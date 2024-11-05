GO=go

APP_NAME=unimates
BUILD_DIR=$(CURDIR)/build
CONFIG_PATH=./config/config.yaml

POSTGRES_DB=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432

POSTGRES_DSN=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable


build: clean
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

run: build
	$(BUILD_DIR)/$(APP_NAME) --config="$(CONFIG_PATH)"

clean: clean-build

clean-build:
	rm -rf $(BUILD_DIR)/

# --------------------------
# Database startup in Docker
# --------------------------

# compose-up:
# 	docker-compose up -d

compose-up:
	docker-compose up --build -d 

compose-down:
	docker-compose down

compose-start:
	docker-compose start postgres

compose-stop:
	docker-compose stop postgres

compose-ps:
	docker-compose ps postgres


# ---------------------
# Run migrations: Goose
# ---------------------

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" create rename_me sql

goose-up:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" up

goose-down:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" down

goose-status:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" status


