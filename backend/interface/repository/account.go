package repository

import (
	"context"

	"budget-app/domain"

	"github.com/go-pg/pg/v10"
)

type accountRepo struct{ db *pg.DB }

func NewAccount(db *pg.DB) domain.AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) GetAll(ctx context.Context) ([]domain.Account, error) {
	var accounts []domain.Account
	err := r.db.WithContext(ctx).Model(&accounts).OrderExpr("id").Select()
	return accounts, err
}

func (r *accountRepo) Create(ctx context.Context, a *domain.Account) error {
	_, err := r.db.WithContext(ctx).Model(a).Returning("id, created_at").Insert()
	return err
}

func (r *accountRepo) Update(ctx context.Context, a *domain.Account) error {
	_, err := r.db.WithContext(ctx).Model(a).
		Column("name", "type", "balance").
		WherePK().Update()
	return err
}

func (r *accountRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.WithContext(ctx).Model((*domain.Account)(nil)).
		Where("id = ?", id).Delete()
	return err
}

func (r *accountRepo) UpdateBalance(ctx context.Context, accountID int, delta float64) error {
	_, err := r.db.WithContext(ctx).Model((*domain.Account)(nil)).
		Set("balance = balance + ?", delta).
		Where("id = ?", accountID).
		Update()
	return err
}
