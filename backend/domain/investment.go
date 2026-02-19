package domain

import (
	"context"
	"time"
)

type Investment struct {
	tableName struct{}  `pg:"investments"`
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Frequency string    `json:"frequency"`
	AccountID *int      `json:"account_id"`
	Category  string    `json:"category"`
	Notes     string    `json:"notes"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type InvestmentRepository interface {
	GetAll(ctx context.Context) ([]Investment, error)
	Create(ctx context.Context, inv *Investment) error
	Update(ctx context.Context, inv *Investment) error
	Delete(ctx context.Context, id int) error
	MonthlyTotal(ctx context.Context) (float64, error)
}
