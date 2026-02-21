package repository

import (
	"context"

	"budget-app/domain"

	"github.com/go-pg/pg/v10"
)

type debtRepo struct{ db *pg.DB }

func NewDebt(db *pg.DB) domain.DebtRepository {
	return &debtRepo{db: db}
}

func (r *debtRepo) GetAll(ctx context.Context) ([]domain.Debt, error) {
	var debts []domain.Debt
	err := r.db.WithContext(ctx).Model(&debts).OrderExpr("id").Select()
	return debts, err
}

func (r *debtRepo) GetByID(ctx context.Context, id int) (*domain.Debt, error) {
	d := &domain.Debt{ID: id}
	err := r.db.WithContext(ctx).Model(d).WherePK().Select()
	return d, err
}

func (r *debtRepo) Create(ctx context.Context, d *domain.Debt) error {
	_, err := r.db.WithContext(ctx).Model(d).Returning("id, created_at").Insert()
	return err
}

func (r *debtRepo) Update(ctx context.Context, d *domain.Debt) error {
	_, err := r.db.WithContext(ctx).Model(d).
		Column("name", "type", "original_amount", "current_balance", "interest_rate", "minimum_payment", "due_day", "status", "notes").
		WherePK().Update()
	return err
}

func (r *debtRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.WithContext(ctx).Model((*domain.Debt)(nil)).
		Where("id = ?", id).Delete()
	return err
}
