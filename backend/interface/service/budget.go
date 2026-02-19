package service

import (
	"context"

	"budget-app/domain"
)

type BudgetService struct {
	repo domain.BudgetRepository
}

func NewBudget(repo domain.BudgetRepository) *BudgetService {
	return &BudgetService{repo: repo}
}

func (s *BudgetService) GetAll(ctx context.Context) ([]domain.Budget, error) {
	return s.repo.GetAll(ctx)
}

func (s *BudgetService) Create(ctx context.Context, b *domain.Budget) error {
	return s.repo.Create(ctx, b)
}

func (s *BudgetService) Update(ctx context.Context, b *domain.Budget) error {
	return s.repo.Update(ctx, b)
}

func (s *BudgetService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
