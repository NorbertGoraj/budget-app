package domain

import (
	"context"
	"time"
)

type Budget struct {
	tableName    struct{}  `pg:"budgets"`
	ID           int       `json:"id"`
	Category     string    `json:"category"`
	MonthlyLimit float64   `json:"monthly_limit"`
	CreatedAt    time.Time `json:"created_at"`
}

type BudgetRepository interface {
	GetAll(ctx context.Context) ([]Budget, error)
	Create(ctx context.Context, b *Budget) error
	Update(ctx context.Context, b *Budget) error
	Delete(ctx context.Context, id int) error
}
