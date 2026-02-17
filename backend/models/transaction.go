package models

import (
	"context"
	"fmt"
	"time"

	"budget-app/db"
)

type Transaction struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"account_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Type        string    `json:"type"`
	Date        string    `json:"date"`
	Imported    bool      `json:"imported"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetTransactions(month, accountID, category string) ([]Transaction, error) {
	query := "SELECT id, account_id, amount, description, category, type, date, imported, created_at FROM transactions WHERE 1=1"
	args := []any{}
	argIdx := 1

	if month != "" {
		query += " AND TO_CHAR(date, 'YYYY-MM') = $" + itoa(argIdx)
		args = append(args, month)
		argIdx++
	}
	if accountID != "" {
		query += " AND account_id = $" + itoa(argIdx)
		args = append(args, accountID)
		argIdx++
	}
	if category != "" {
		query += " AND category = $" + itoa(argIdx)
		args = append(args, category)
		argIdx++
	}
	query += " ORDER BY date DESC, id DESC"

	rows, err := db.Pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.Amount, &t.Description, &t.Category, &t.Type, &t.Date, &t.Imported, &t.CreatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, nil
}

func CreateTransaction(t *Transaction) error {
	err := db.Pool.QueryRow(context.Background(),
		`INSERT INTO transactions (account_id, amount, description, category, type, date, imported)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`,
		t.AccountID, t.Amount, t.Description, t.Category, t.Type, t.Date, t.Imported,
	).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return err
	}

	delta := t.Amount
	if t.Type == "expense" {
		delta = -delta
	}
	return UpdateAccountBalance(t.AccountID, delta)
}

func UpdateTransaction(t *Transaction) error {
	_, err := db.Pool.Exec(context.Background(),
		`UPDATE transactions SET account_id=$1, amount=$2, description=$3, category=$4, type=$5, date=$6
		 WHERE id=$7`,
		t.AccountID, t.Amount, t.Description, t.Category, t.Type, t.Date, t.ID)
	return err
}

func DeleteTransaction(id int) error {
	var accountID int
	var amount float64
	var txType string
	err := db.Pool.QueryRow(context.Background(),
		"SELECT account_id, amount, type FROM transactions WHERE id=$1", id,
	).Scan(&accountID, &amount, &txType)
	if err != nil {
		return err
	}

	delta := -amount
	if txType == "expense" {
		delta = amount
	}
	if err := UpdateAccountBalance(accountID, delta); err != nil {
		return err
	}

	_, err = db.Pool.Exec(context.Background(), "DELETE FROM transactions WHERE id=$1", id)
	return err
}

func TransactionExists(date, description string, amount float64) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM transactions WHERE date=$1 AND description=$2 AND amount=$3)",
		date, description, amount,
	).Scan(&exists)
	return exists, err
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
