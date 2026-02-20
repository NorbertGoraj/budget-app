APP     = bin/budget-api
MIGRATE = bin/migrate

# Pass WAIT=1 to sleep 10s before connecting, e.g. make migrate WAIT=1
WAIT_FLAG = $(if $(WAIT),--wait,)

.PHONY: build build-migrate build-all run \
        migrate migrate-down migrate-reset migrate-status migrate-version migrate-help \
        mock test \
        fmt \
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

migrate-down: build-migrate
	./$(MIGRATE) $(WAIT_FLAG) down

migrate-reset: build-migrate
	./$(MIGRATE) $(WAIT_FLAG) reset

migrate-status: build-migrate
	./$(MIGRATE) $(WAIT_FLAG) status

migrate-version: build-migrate
	./$(MIGRATE) $(WAIT_FLAG) version

migrate-help: build-migrate
	./$(MIGRATE) --help

## ── Mocks & Tests ─────────────────────────────────────────────────────────────

mock:
	mockery

test:
	go test ./...

## ── Format ────────────────────────────────────────────────────────────────────

fmt:
	goimports -w .

## ── Misc ──────────────────────────────────────────────────────────────────────

clean:
	rm -rf bin/
