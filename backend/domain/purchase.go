package domain

import (
	"context"
	"time"
)

type PlannedPurchase struct {
	tableName     struct{}  `pg:"planned_purchases"`
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	EstimatedCost float64   `json:"estimated_cost"`
	Category      string    `json:"category"`
	Priority      string    `json:"priority"`
	TargetMonth   string    `json:"target_month"`
	Notes         string    `json:"notes"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type PurchaseRepository interface {
	GetAll(ctx context.Context) ([]PlannedPurchase, error)
	Create(ctx context.Context, p *PlannedPurchase) error
	Update(ctx context.Context, p *PlannedPurchase) error
	Delete(ctx context.Context, id int) error
}
