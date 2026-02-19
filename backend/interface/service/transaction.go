package service

import (
	"context"

	"budget-app/domain"
)

type TransactionService struct {
	txRepo      domain.TransactionRepository
	accountRepo domain.AccountRepository
}

func NewTransaction(txRepo domain.TransactionRepository, accountRepo domain.AccountRepository) *TransactionService {
	return &TransactionService{txRepo: txRepo, accountRepo: accountRepo}
}

func (s *TransactionService) GetAll(ctx context.Context, f domain.TransactionFilter) ([]domain.Transaction, error) {
	return s.txRepo.GetAll(ctx, f)
}

func (s *TransactionService) Create(ctx context.Context, t *domain.Transaction) error {
	if err := s.txRepo.Create(ctx, t); err != nil {
		return err
	}
	delta := t.Amount
	if t.Type == "expense" {
		delta = -delta
	}
	return s.accountRepo.UpdateBalance(ctx, t.AccountID, delta)
}

func (s *TransactionService) Update(ctx context.Context, t *domain.Transaction) error {
	return s.txRepo.Update(ctx, t)
}

func (s *TransactionService) Delete(ctx context.Context, id int) error {
	t, err := s.txRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	delta := -t.Amount
	if t.Type == "expense" {
		delta = t.Amount
	}
	if err := s.accountRepo.UpdateBalance(ctx, t.AccountID, delta); err != nil {
		return err
	}

	return s.txRepo.Delete(ctx, id)
}

func (s *TransactionService) Exists(ctx context.Context, date, description string, amount float64) (bool, error) {
	return s.txRepo.Exists(ctx, date, description, amount)
}
