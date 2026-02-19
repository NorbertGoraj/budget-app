package service

import (
	"context"

	"budget-app/domain"
)

type AccountService struct {
	repo domain.AccountRepository
}

func NewAccount(repo domain.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) GetAll(ctx context.Context) ([]domain.Account, error) {
	return s.repo.GetAll(ctx)
}

func (s *AccountService) Create(ctx context.Context, a *domain.Account) error {
	return s.repo.Create(ctx, a)
}

func (s *AccountService) Update(ctx context.Context, a *domain.Account) error {
	return s.repo.Update(ctx, a)
}

func (s *AccountService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
