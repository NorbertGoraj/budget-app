package domain

import (
	"context"
	"time"
)

type Account struct {
	tableName struct{}  `pg:"accounts"`
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountRepository interface {
	GetAll(ctx context.Context) ([]Account, error)
	Create(ctx context.Context, a *Account) error
	Update(ctx context.Context, a *Account) error
	Delete(ctx context.Context, id int) error
	UpdateBalance(ctx context.Context, accountID int, delta float64) error
}
