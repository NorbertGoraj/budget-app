package models

import (
	"context"
	"time"

	"budget-app/db"
)

type Budget struct {
	ID           int       `json:"id"`
	Category     string    `json:"category"`
	MonthlyLimit float64   `json:"monthly_limit"`
	CreatedAt    time.Time `json:"created_at"`
}

func GetAllBudgets() ([]Budget, error) {
	rows, err := db.Pool.Query(context.Background(),
		"SELECT id, category, monthly_limit, created_at FROM budgets ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.Category, &b.MonthlyLimit, &b.CreatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func CreateBudget(b *Budget) error {
	return db.Pool.QueryRow(context.Background(),
		"INSERT INTO budgets (category, monthly_limit) VALUES ($1, $2) RETURNING id, created_at",
		b.Category, b.MonthlyLimit).Scan(&b.ID, &b.CreatedAt)
}

func UpdateBudget(b *Budget) error {
	_, err := db.Pool.Exec(context.Background(),
		"UPDATE budgets SET category=$1, monthly_limit=$2 WHERE id=$3",
		b.Category, b.MonthlyLimit, b.ID)
	return err
}

func DeleteBudget(id int) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM budgets WHERE id=$1", id)
	return err
}
