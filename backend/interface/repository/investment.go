package repository

import (
	"context"

	"budget-app/domain"

	"github.com/go-pg/pg/v10"
)

type investmentRepo struct{ db *pg.DB }

func NewInvestment(db *pg.DB) domain.InvestmentRepository {
	return &investmentRepo{db: db}
}

func (r *investmentRepo) GetAll(ctx context.Context) ([]domain.Investment, error) {
	var investments []domain.Investment
	err := r.db.WithContext(ctx).Model(&investments).OrderExpr("id").Select()
	return investments, err
}

func (r *investmentRepo) Create(ctx context.Context, inv *domain.Investment) error {
	_, err := r.db.WithContext(ctx).Model(inv).Returning("id, created_at").Insert()
	return err
}

func (r *investmentRepo) Update(ctx context.Context, inv *domain.Investment) error {
	_, err := r.db.WithContext(ctx).Model(inv).
		Column("name", "type", "amount", "frequency", "account_id", "category", "notes", "status").
		WherePK().Update()
	return err
}

func (r *investmentRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.WithContext(ctx).Model((*domain.Investment)(nil)).
		Where("id = ?", id).Delete()
	return err
}
