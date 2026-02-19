APP     = bin/budget-api
MIGRATE = bin/migrate

# Pass WAIT=1 to sleep 10s before connecting, e.g. make migrate WAIT=1
WAIT_FLAG = $(if $(WAIT),--wait,)

.PHONY: build build-migrate build-all run \
        migrate migrate-status migrate-version migrate-help \
        clean

## ── Build ─────────────────────────────────────────────────────────────────────

build:
	go build -o $(APP) .

build-migrate:
	go build -o $(MIGRATE) ./cmd/migrate

build-all: build build-migrate

## ── Run ───────────────────────────────────────────────────────────────────────

run: build
	./$(APP)

## ── Migrations ────────────────────────────────────────────────────────────────

migrate: build-migrate
	./$(MIGRATE) $(WAIT_FLAG) up

migrate-status: build-migrate
	./$(MIGRATE) $(WAIT_FLAG) status

migrate-version: build-migrate
	./$(MIGRATE) $(WAIT_FLAG) version

migrate-help: build-migrate
	./$(MIGRATE) --help

## ── Misc ──────────────────────────────────────────────────────────────────────

clean:
	rm -rf bin/
