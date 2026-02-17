package models

import (
	"context"
	"time"

	"budget-app/db"
)

type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllAccounts() ([]Account, error) {
	rows, err := db.Pool.Query(context.Background(),
		"SELECT id, name, type, balance, created_at FROM accounts ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func CreateAccount(a *Account) error {
	return db.Pool.QueryRow(context.Background(),
		"INSERT INTO accounts (name, type, balance) VALUES ($1, $2, $3) RETURNING id, created_at",
		a.Name, a.Type, a.Balance).Scan(&a.ID, &a.CreatedAt)
}

func UpdateAccount(a *Account) error {
	_, err := db.Pool.Exec(context.Background(),
		"UPDATE accounts SET name=$1, type=$2, balance=$3 WHERE id=$4",
		a.Name, a.Type, a.Balance, a.ID)
	return err
}

func DeleteAccount(id int) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM accounts WHERE id=$1", id)
	return err
}

func UpdateAccountBalance(accountID int, delta float64) error {
	_, err := db.Pool.Exec(context.Background(),
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2", delta, accountID)
	return err
}
