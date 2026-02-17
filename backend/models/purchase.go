package models

import (
	"context"
	"time"

	"budget-app/db"
)

type PlannedPurchase struct {
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

func GetAllPurchases() ([]PlannedPurchase, error) {
	rows, err := db.Pool.Query(context.Background(),
		"SELECT id, name, estimated_cost, category, priority, target_month, notes, status, created_at FROM planned_purchases ORDER BY target_month, priority, id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []PlannedPurchase
	for rows.Next() {
		var p PlannedPurchase
		if err := rows.Scan(&p.ID, &p.Name, &p.EstimatedCost, &p.Category, &p.Priority, &p.TargetMonth, &p.Notes, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}

func CreatePurchase(p *PlannedPurchase) error {
	return db.Pool.QueryRow(context.Background(),
		`INSERT INTO planned_purchases (name, estimated_cost, category, priority, target_month, notes, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`,
		p.Name, p.EstimatedCost, p.Category, p.Priority, p.TargetMonth, p.Notes, p.Status,
	).Scan(&p.ID, &p.CreatedAt)
}

func UpdatePurchase(p *PlannedPurchase) error {
	_, err := db.Pool.Exec(context.Background(),
		`UPDATE planned_purchases SET name=$1, estimated_cost=$2, category=$3, priority=$4, target_month=$5, notes=$6, status=$7
		 WHERE id=$8`,
		p.Name, p.EstimatedCost, p.Category, p.Priority, p.TargetMonth, p.Notes, p.Status, p.ID)
	return err
}

func DeletePurchase(id int) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM planned_purchases WHERE id=$1", id)
	return err
}
