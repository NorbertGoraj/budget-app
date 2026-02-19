package repository

import (
	"context"

	"budget-app/domain"

	"github.com/go-pg/pg/v10"
)

type budgetRepo struct{ db *pg.DB }

func NewBudget(db *pg.DB) domain.BudgetRepository {
	return &budgetRepo{db: db}
}

func (r *budgetRepo) GetAll(ctx context.Context) ([]domain.Budget, error) {
	var budgets []domain.Budget
	err := r.db.WithContext(ctx).Model(&budgets).OrderExpr("id").Select()
	return budgets, err
}

func (r *budgetRepo) Create(ctx context.Context, b *domain.Budget) error {
	_, err := r.db.WithContext(ctx).Model(b).Returning("id, created_at").Insert()
	return err
}

func (r *budgetRepo) Update(ctx context.Context, b *domain.Budget) error {
	_, err := r.db.WithContext(ctx).Model(b).
		Column("category", "monthly_limit").
		WherePK().Update()
	return err
}

func (r *budgetRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.WithContext(ctx).Model((*domain.Budget)(nil)).
		Where("id = ?", id).Delete()
	return err
}
