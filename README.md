# Budget App

A full-stack personal finance manager — track accounts, transactions, budgets, planned purchases, investments, and debts from a single dashboard.

---

## Architecture

```
awesomeProject7/
├── backend/          # Go API
└── frontend/         # Angular SPA
```

### Backend — Clean Architecture

Follows [Eminetto's Go Clean Architecture](https://eltonminetto.dev/en/post/2018-03-05-clean-architecture-using-go/) with four concentric rings:

```
domain/                   Enterprise business rules
  account.go              Structs + repository interfaces
  transaction.go
  budget.go
  purchase.go
  investment.go
  debt.go                 Debt + DebtPayment structs, DebtRepository + DebtPaymentRepository interfaces
  dashboard.go            Aggregated response types (incl. DebtSummary)

interface/
  controller/             HTTP handlers (Gin) — outermost ring
  service/                Use-case orchestration + unit tests
  repository/             go-pg implementations of domain interfaces

infrastructure/
  db.go                   go-pg connection
  httphandler/            Reusable HTTP server lifecycle (graceful shutdown)
  metrics/                Prometheus server + Gin middleware

mocks/domain/             mockery-generated mocks for all repository interfaces
storage/migrations/       Embedded SQL files + migration tracker
cmd/migrate/              Standalone migration binary (Cobra CLI)
appcontext/               Dependency wiring + signal handling
utils/                    Generic helpers (EnsureNotNil)
```

**Dependency rule:** `controller → service → domain ← repository`
No layer may import anything from a layer outside it.

### Frontend — Angular 19

Single-page app with Angular Material UI. Communicates with the backend via `/api/*` — Nginx proxies all such requests to the Go API.

```
src/app/
  dashboard/        Aggregated overview (balance, budgets, investments, debts, purchases)
  accounts/         Account CRUD
  transactions/     Transaction list + filters
  budgets/          Monthly budget limits
  purchases/        Planned purchase tracker
  investments/      Investment portfolio
  debts/            Debt management + payment history
  import/           CSV bulk import
  services/
    api.service.ts  Centralised HTTP client
    models.ts       Shared TypeScript interfaces
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| API framework | [Gin](https://github.com/gin-gonic/gin) v1.11 |
| Database | PostgreSQL 17 |
| ORM | [go-pg](https://github.com/go-pg/pg) v10 |
| Metrics | [Prometheus](https://github.com/prometheus/client_golang) + promhttp |
| Concurrency | [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| Frontend | Angular 19 + Angular Material |
| Container | Docker + Docker Compose |

---

## Running with Docker Compose

```bash
docker compose up --build
```

This starts four services in dependency order:

| Service | Description | Port |
|---|---|---|
| `db` | PostgreSQL 17 | `54329` (host) |
| `migrator` | Runs pending migrations then exits | — |
| `backend` | Go API (starts after migrator) | `8081` → `8080`, `9091` |
| `frontend` | Angular served via Nginx | `4201` → `80` |

URLs once running:

- **Frontend:** http://localhost:4201
- **API:** http://localhost:8081/api
- **Metrics:** http://localhost:9091/metrics

Stop everything:

```bash
docker compose down
```

Wipe the database volume too:

```bash
docker compose down -v
```

---

## Running Locally

### Prerequisites

- Go 1.25+
- Node 22+
- PostgreSQL running on `localhost:5432`

### Backend

```bash
cd backend

# Set database URL (or export it in your shell)
export DATABASE_URL=postgres://budget:budget@localhost:5432/budget?sslmode=disable

# Run migrations
make migrate

# Start the API
make run
```

### Frontend

```bash
cd frontend
npm install
npm start        # serves on http://localhost:4200
```

---

## Migrations

Migrations live in `backend/storage/migrations/` as numbered SQL file pairs (`*.up.sql` / `*.down.sql`). The `gopg_migrations` table tracks which have been applied. Each migration runs in a transaction — it either fully applies or fully rolls back.

### Makefile targets

```bash
make migrate            # apply all pending migrations
make migrate-status     # show applied / pending state
make migrate-version    # print count of applied migrations
make migrate-help       # print CLI usage
```

Pass `WAIT=1` to sleep 10 s before connecting (useful when the DB container is still starting):

```bash
make migrate WAIT=1
```

### Migrate binary directly

```bash
cd backend
go build -o bin/migrate ./cmd/migrate

./bin/migrate up              # apply pending
./bin/migrate status          # show state
./bin/migrate version         # print count
./bin/migrate --wait up       # wait 10s then apply
./bin/migrate --help          # full usage
```

### Applied migrations

| # | Migration |
|---|---|
| 1 | `create_accounts` |
| 2 | `create_transactions` |
| 3 | `create_budgets` |
| 4 | `create_planned_purchases` |
| 5 | `create_investments` |
| 6 | `create_debts` |
| 7 | `create_debt_payments` |

### Adding a new migration

Create two files in `backend/storage/migrations/` following the naming pattern:

```
8_your_description.up.sql
8_your_description.down.sql
```

The `init()` function in `migrations.go` discovers and registers all `*.up.sql` / `*.down.sql` pairs automatically — no Go code changes needed. The next `migrate up` run will detect and apply it.

---

## Metrics

The backend exposes Prometheus metrics on `:9091/metrics`.

| Metric | Type | Description |
|---|---|---|
| `http_requests_total` | Counter | Request count by method, path, status |
| `http_request_duration_seconds` | Histogram | Latency by method and path |

---

## Backend Makefile Reference

```bash
make build           # compile budget-api binary  → bin/budget-api
make build-migrate   # compile migrate binary      → bin/migrate
make build-all       # both of the above
make run             # build + run the API
make migrate         # apply pending migrations
make migrate-down    # roll back one migration
make migrate-reset   # roll back all migrations
make migrate-status  # migration state
make migrate-version # applied count
make migrate-help    # CLI help
make test            # run all unit tests
make mock            # regenerate mockery mocks
make fmt             # format code with goimports
make clean           # remove bin/
```

---

## API Reference

All endpoints are prefixed with `/api`. The server listens on `:8080`.

### Accounts
| Method | Path | Description |
|---|---|---|
| GET | `/api/accounts` | List all accounts |
| POST | `/api/accounts` | Create account |
| PUT | `/api/accounts/:id` | Update account |
| DELETE | `/api/accounts/:id` | Delete account |

### Transactions
| Method | Path | Description |
|---|---|---|
| GET | `/api/transactions` | List transactions (`?month=`, `?account_id=`, `?category=`) |
| POST | `/api/transactions` | Create transaction — updates account balance |
| PUT | `/api/transactions/:id` | Update transaction |
| DELETE | `/api/transactions/:id` | Delete transaction — reverses balance change |

### Budgets
| Method | Path | Description |
|---|---|---|
| GET | `/api/budgets` | List all budgets |
| POST | `/api/budgets` | Create budget |
| PUT | `/api/budgets/:id` | Update budget |
| DELETE | `/api/budgets/:id` | Delete budget |

### Planned Purchases
| Method | Path | Description |
|---|---|---|
| GET | `/api/purchases` | List all planned purchases |
| POST | `/api/purchases` | Create planned purchase |
| PUT | `/api/purchases/:id` | Update planned purchase |
| DELETE | `/api/purchases/:id` | Delete planned purchase |

### Investments
| Method | Path | Description |
|---|---|---|
| GET | `/api/investments` | List all investments |
| POST | `/api/investments` | Create investment |
| PUT | `/api/investments/:id` | Update investment |
| DELETE | `/api/investments/:id` | Delete investment |

### Debts
| Method | Path | Description |
|---|---|---|
| GET | `/api/debts` | List all debts |
| POST | `/api/debts` | Create debt |
| PUT | `/api/debts/:id` | Update debt |
| DELETE | `/api/debts/:id` | Delete debt (cascades to payments) |
| GET | `/api/debts/:id/payments` | List payments for a debt |
| POST | `/api/debts/:id/payments` | Record payment — reduces `current_balance`, sets `paid_off` at zero |
| DELETE | `/api/debts/payments/:id` | Delete payment — reverses balance, reactivates `paid_off` debt |

### Dashboard
| Method | Path | Description |
|---|---|---|
| GET | `/api/dashboard` | Aggregated overview: balances, monthly summary, budget status, purchase affordability, investments, debt summary |

### Import
| Method | Path | Description |
|---|---|---|
| POST | `/api/import/csv` | Bulk import transactions from CSV (`multipart/form-data`: `account_id`, `file`) |

### Metrics
| Method | Path | Description |
|---|---|---|
| GET | `/metrics` | Prometheus metrics (port `9091`) |

---

## Testing

Service-layer unit tests use mockery-generated mocks. No database required to run them.

```bash
cd backend && make test
```

Mocks live in `mocks/domain/` and are generated from the interfaces in `domain/`. After changing any repository interface run `make mock` to regenerate.

Covered services: `AccountService`, `BudgetService`, `TransactionService`, `PurchaseService`, `InvestmentService`, `DebtService`, `DashboardService`.

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://budget:budget@localhost:5432/budget?sslmode=disable` | PostgreSQL connection string |
