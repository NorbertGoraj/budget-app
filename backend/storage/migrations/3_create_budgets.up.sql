CREATE TABLE IF NOT EXISTS budgets (
    id            SERIAL PRIMARY KEY,
    category      TEXT          NOT NULL UNIQUE,
    monthly_limit NUMERIC(12,2) NOT NULL,
    created_at    TIMESTAMP     NOT NULL DEFAULT NOW()
);
