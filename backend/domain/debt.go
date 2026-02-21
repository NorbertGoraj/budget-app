package domain

import (
	"context"
	"time"
)

type Debt struct {
	tableName      struct{}  `pg:"debts"`
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"` // credit_card | loan | mortgage | student_loan | car_loan | other
	OriginalAmount float64   `json:"original_amount"`
	CurrentBalance float64   `json:"current_balance"`
	InterestRate   float64   `json:"interest_rate"` // APR %
	MinimumPayment float64   `json:"minimum_payment"`
	DueDay         int       `json:"due_day"` // 1-28: day of month payment is due
	Status         string    `json:"status"`  // active | paid_off
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}

type DebtPayment struct {
	tableName struct{}  `pg:"debt_payments"`
	ID        int       `json:"id"`
	DebtID    int       `json:"debt_id"`
	Amount    float64   `json:"amount"`
	PaidAt    string    `json:"paid_at"` // DATE as "YYYY-MM-DD"
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type DebtRepository interface {
	GetAll(ctx context.Context) ([]Debt, error)
	GetByID(ctx context.Context, id int) (*Debt, error)
	Create(ctx context.Context, d *Debt) error
	Update(ctx context.Context, d *Debt) error
	Delete(ctx context.Context, id int) error
}

type DebtPaymentRepository interface {
	GetByDebtID(ctx context.Context, debtID int) ([]DebtPayment, error)
	GetByID(ctx context.Context, id int) (*DebtPayment, error)
	Create(ctx context.Context, p *DebtPayment) error
	Delete(ctx context.Context, id int) error
}
