package repository

import (
	"context"

	"budget-app/domain"

	"github.com/go-pg/pg/v10"
)

type debtPaymentRepo struct{ db *pg.DB }

func NewDebtPayment(db *pg.DB) domain.DebtPaymentRepository {
	return &debtPaymentRepo{db: db}
}

func (r *debtPaymentRepo) GetByDebtID(ctx context.Context, debtID int) ([]domain.DebtPayment, error) {
	var payments []domain.DebtPayment
	err := r.db.WithContext(ctx).Model(&payments).
		Where("debt_id = ?", debtID).
		OrderExpr("paid_at DESC, id DESC").
		Select()
	return payments, err
}

func (r *debtPaymentRepo) GetByID(ctx context.Context, id int) (*domain.DebtPayment, error) {
	p := &domain.DebtPayment{ID: id}
	err := r.db.WithContext(ctx).Model(p).WherePK().Select()
	return p, err
}

func (r *debtPaymentRepo) Create(ctx context.Context, p *domain.DebtPayment) error {
	_, err := r.db.WithContext(ctx).Model(p).Returning("id, created_at").Insert()
	return err
}

func (r *debtPaymentRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.WithContext(ctx).Model((*domain.DebtPayment)(nil)).
		Where("id = ?", id).Delete()
	return err
}
