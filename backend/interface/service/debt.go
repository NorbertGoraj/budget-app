package service

import (
	"context"

	"budget-app/domain"
)

type DebtService struct {
	debtRepo        domain.DebtRepository
	debtPaymentRepo domain.DebtPaymentRepository
}

func NewDebt(debtRepo domain.DebtRepository, debtPaymentRepo domain.DebtPaymentRepository) *DebtService {
	return &DebtService{
		debtRepo:        debtRepo,
		debtPaymentRepo: debtPaymentRepo,
	}
}

func (s *DebtService) GetAll(ctx context.Context) ([]domain.Debt, error) {
	return s.debtRepo.GetAll(ctx)
}

func (s *DebtService) GetByID(ctx context.Context, id int) (*domain.Debt, error) {
	return s.debtRepo.GetByID(ctx, id)
}

func (s *DebtService) Create(ctx context.Context, d *domain.Debt) error {
	d.Status = "active"
	return s.debtRepo.Create(ctx, d)
}

func (s *DebtService) Update(ctx context.Context, d *domain.Debt) error {
	return s.debtRepo.Update(ctx, d)
}

func (s *DebtService) Delete(ctx context.Context, id int) error {
	return s.debtRepo.Delete(ctx, id)
}

func (s *DebtService) GetPayments(ctx context.Context, debtID int) ([]domain.DebtPayment, error) {
	return s.debtPaymentRepo.GetByDebtID(ctx, debtID)
}

func (s *DebtService) RecordPayment(ctx context.Context, debtID int, p *domain.DebtPayment) error {
	debt, err := s.debtRepo.GetByID(ctx, debtID)
	if err != nil {
		return err
	}

	p.DebtID = debtID
	if err := s.debtPaymentRepo.Create(ctx, p); err != nil {
		return err
	}

	newBalance := debt.CurrentBalance - p.Amount
	if newBalance <= 0 {
		newBalance = 0
		debt.Status = "paid_off"
	}
	debt.CurrentBalance = newBalance

	return s.debtRepo.Update(ctx, debt)
}

func (s *DebtService) DeletePayment(ctx context.Context, paymentID int) error {
	payment, err := s.debtPaymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return err
	}

	debt, err := s.debtRepo.GetByID(ctx, payment.DebtID)
	if err != nil {
		return err
	}

	debt.CurrentBalance += payment.Amount
	if debt.Status == "paid_off" {
		debt.Status = "active"
	}

	if err := s.debtRepo.Update(ctx, debt); err != nil {
		return err
	}

	return s.debtPaymentRepo.Delete(ctx, paymentID)
}
