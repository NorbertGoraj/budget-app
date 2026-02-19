package repository

import (
	"context"

	"budget-app/domain"

	"github.com/go-pg/pg/v10"
)

type transactionRepo struct{ db *pg.DB }

func NewTransaction(db *pg.DB) domain.TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) GetAll(ctx context.Context, f domain.TransactionFilter) ([]domain.Transaction, error) {
	var txns []domain.Transaction
	q := r.db.WithContext(ctx).Model(&txns).OrderExpr("date DESC, id DESC")
	if f.Month != "" {
		q = q.Where("TO_CHAR(date, 'YYYY-MM') = ?", f.Month)
	}
	if f.AccountID != "" {
		q = q.Where("account_id = ?", f.AccountID)
	}
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	return txns, q.Select()
}

func (r *transactionRepo) GetByID(ctx context.Context, id int) (*domain.Transaction, error) {
	t := &domain.Transaction{ID: id}
	err := r.db.WithContext(ctx).Model(t).WherePK().Select()
	return t, err
}

func (r *transactionRepo) Create(ctx context.Context, t *domain.Transaction) error {
	_, err := r.db.WithContext(ctx).Model(t).Returning("id, created_at").Insert()
	return err
}

func (r *transactionRepo) Update(ctx context.Context, t *domain.Transaction) error {
	_, err := r.db.WithContext(ctx).Model(t).
		Column("account_id", "amount", "description", "category", "type", "date").
		WherePK().Update()
	return err
}

func (r *transactionRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.WithContext(ctx).Model((*domain.Transaction)(nil)).
		Where("id = ?", id).Delete()
	return err
}

func (r *transactionRepo) Exists(ctx context.Context, date, description string, amount float64) (bool, error) {
	return r.db.WithContext(ctx).Model((*domain.Transaction)(nil)).
		Where("date = ? AND description = ? AND amount = ?", date, description, amount).
		Exists()
}

func (r *transactionRepo) MonthlySums(ctx context.Context, month string) (income, expenses float64, err error) {
	var result struct {
		Income   float64
		Expenses float64
	}
	_, err = r.db.WithContext(ctx).QueryOne(&result, `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expenses
		FROM transactions
		WHERE TO_CHAR(date, 'YYYY-MM') = ?`, month)
	return result.Income, result.Expenses, err
}

func (r *transactionRepo) SpentByCategory(ctx context.Context, category, month string) (float64, error) {
	var spent float64
	_, err := r.db.WithContext(ctx).QueryOne(pg.Scan(&spent), `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE category = ? AND type = 'expense' AND TO_CHAR(date, 'YYYY-MM') = ?`,
		category, month)
	return spent, err
}
