CREATE TABLE IF NOT EXISTS planned_purchases (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    estimated_cost NUMERIC(12,2) NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('high', 'medium', 'low')),
    target_month TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'purchased', 'cancelled', 'deferred')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
