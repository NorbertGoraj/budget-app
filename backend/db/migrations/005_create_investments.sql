CREATE TABLE IF NOT EXISTS investments (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('recurring', 'one_time')),
    amount NUMERIC(12,2) NOT NULL,
    frequency TEXT NOT NULL DEFAULT '' CHECK (frequency IN ('', 'monthly', 'weekly', 'quarterly', 'yearly')),
    account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
    category TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'active', 'paused')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
