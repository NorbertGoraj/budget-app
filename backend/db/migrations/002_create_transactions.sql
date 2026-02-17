CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    amount NUMERIC(12,2) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category TEXT NOT NULL DEFAULT 'Uncategorized',
    type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    imported BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
