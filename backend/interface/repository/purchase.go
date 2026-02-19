package repository

import (
	"context"

	"budget-app/domain"

	"github.com/go-pg/pg/v10"
)

type purchaseRepo struct{ db *pg.DB }

func NewPurchase(db *pg.DB) domain.PurchaseRepository {
	return &purchaseRepo{db: db}
}

func (r *purchaseRepo) GetAll(ctx context.Context) ([]domain.PlannedPurchase, error) {
	var purchases []domain.PlannedPurchase
	err := r.db.WithContext(ctx).Model(&purchases).
		OrderExpr("target_month, priority, id").Select()
	return purchases, err
}

func (r *purchaseRepo) Create(ctx context.Context, p *domain.PlannedPurchase) error {
	_, err := r.db.WithContext(ctx).Model(p).Returning("id, created_at").Insert()
	return err
}

func (r *purchaseRepo) Update(ctx context.Context, p *domain.PlannedPurchase) error {
	_, err := r.db.WithContext(ctx).Model(p).
		Column("name", "estimated_cost", "category", "priority", "target_month", "notes", "status").
		WherePK().Update()
	return err
}

func (r *purchaseRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.WithContext(ctx).Model((*domain.PlannedPurchase)(nil)).
		Where("id = ?", id).Delete()
	return err
}
