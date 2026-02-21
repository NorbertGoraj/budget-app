CREATE TABLE IF NOT EXISTS debt_payments (
    id         SERIAL PRIMARY KEY,
    debt_id    INTEGER       NOT NULL REFERENCES debts(id) ON DELETE CASCADE,
    amount     NUMERIC(12,2) NOT NULL,
    paid_at    DATE          NOT NULL DEFAULT CURRENT_DATE,
    notes      TEXT          NOT NULL DEFAULT '',
    created_at TIMESTAMP     NOT NULL DEFAULT NOW()
);
