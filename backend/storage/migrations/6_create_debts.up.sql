CREATE TABLE IF NOT EXISTS debts (
    id              SERIAL PRIMARY KEY,
    name            TEXT          NOT NULL,
    type            TEXT          NOT NULL CHECK (type IN ('credit_card','loan','mortgage','student_loan','car_loan','other')),
    original_amount NUMERIC(12,2) NOT NULL,
    current_balance NUMERIC(12,2) NOT NULL,
    interest_rate   NUMERIC(5,2)  NOT NULL DEFAULT 0,
    minimum_payment NUMERIC(12,2) NOT NULL DEFAULT 0,
    due_day         SMALLINT      NOT NULL DEFAULT 1 CHECK (due_day BETWEEN 1 AND 28),
    status          TEXT          NOT NULL DEFAULT 'active' CHECK (status IN ('active','paid_off')),
    notes           TEXT          NOT NULL DEFAULT '',
    created_at      TIMESTAMP     NOT NULL DEFAULT NOW()
);
