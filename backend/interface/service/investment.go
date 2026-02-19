package service

import (
	"context"

	"budget-app/domain"
)

type InvestmentService struct {
	repo domain.InvestmentRepository
}

func NewInvestment(repo domain.InvestmentRepository) *InvestmentService {
	return &InvestmentService{repo: repo}
}

func (s *InvestmentService) GetAll(ctx context.Context) ([]domain.Investment, error) {
	return s.repo.GetAll(ctx)
}

func (s *InvestmentService) Create(ctx context.Context, inv *domain.Investment) error {
	return s.repo.Create(ctx, inv)
}

func (s *InvestmentService) Update(ctx context.Context, inv *domain.Investment) error {
	return s.repo.Update(ctx, inv)
}

func (s *InvestmentService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
