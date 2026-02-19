package domain

import (
	"context"
	"time"
)

type Transaction struct {
	tableName   struct{}  `pg:"transactions"`
	ID          int       `json:"id"`
	AccountID   int       `json:"account_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Type        string    `json:"type"`
	Date        string    `json:"date"`
	Imported    bool      `json:"imported"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransactionFilter struct {
	Month     string
	AccountID string
	Category  string
}

type TransactionRepository interface {
	GetAll(ctx context.Context, f TransactionFilter) ([]Transaction, error)
	GetByID(ctx context.Context, id int) (*Transaction, error)
	Create(ctx context.Context, t *Transaction) error
	Update(ctx context.Context, t *Transaction) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, date, description string, amount float64) (bool, error)
	MonthlySums(ctx context.Context, month string) (income, expenses float64, err error)
	SpentByCategory(ctx context.Context, category, month string) (float64, error)
}
