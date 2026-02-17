package models

import (
	"context"
	"time"

	"budget-app/db"
)

type Investment struct {
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

func GetAllInvestments() ([]Investment, error) {
	rows, err := db.Pool.Query(context.Background(),
		"SELECT id, name, type, amount, frequency, account_id, category, notes, status, created_at FROM investments ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []Investment
	for rows.Next() {
		var inv Investment
		if err := rows.Scan(&inv.ID, &inv.Name, &inv.Type, &inv.Amount, &inv.Frequency, &inv.AccountID, &inv.Category, &inv.Notes, &inv.Status, &inv.CreatedAt); err != nil {
			return nil, err
		}
		investments = append(investments, inv)
	}
	return investments, nil
}

func CreateInvestment(inv *Investment) error {
	return db.Pool.QueryRow(context.Background(),
		`INSERT INTO investments (name, type, amount, frequency, account_id, category, notes, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`,
		inv.Name, inv.Type, inv.Amount, inv.Frequency, inv.AccountID, inv.Category, inv.Notes, inv.Status,
	).Scan(&inv.ID, &inv.CreatedAt)
}

func UpdateInvestment(inv *Investment) error {
	_, err := db.Pool.Exec(context.Background(),
		`UPDATE investments SET name=$1, type=$2, amount=$3, frequency=$4, account_id=$5, category=$6, notes=$7, status=$8
		 WHERE id=$9`,
		inv.Name, inv.Type, inv.Amount, inv.Frequency, inv.AccountID, inv.Category, inv.Notes, inv.Status, inv.ID)
	return err
}

func DeleteInvestment(id int) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM investments WHERE id=$1", id)
	return err
}

func GetMonthlyInvestmentTotal() (float64, error) {
	var total float64
	err := db.Pool.QueryRow(context.Background(), `
		SELECT COALESCE(SUM(
			CASE
				WHEN type = 'one_time' THEN 0
				WHEN frequency = 'weekly' THEN amount * 4.33
				WHEN frequency = 'monthly' THEN amount
				WHEN frequency = 'quarterly' THEN amount / 3
				WHEN frequency = 'yearly' THEN amount / 12
				ELSE 0
			END
		), 0) FROM investments WHERE status = 'active'
	`).Scan(&total)
	return total, err
}
