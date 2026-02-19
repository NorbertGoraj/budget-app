package service

import (
	"context"

	"budget-app/domain"
)

type PurchaseService struct {
	repo domain.PurchaseRepository
}

func NewPurchase(repo domain.PurchaseRepository) *PurchaseService {
	return &PurchaseService{repo: repo}
}

func (s *PurchaseService) GetAll(ctx context.Context) ([]domain.PlannedPurchase, error) {
	return s.repo.GetAll(ctx)
}

func (s *PurchaseService) Create(ctx context.Context, p *domain.PlannedPurchase) error {
	return s.repo.Create(ctx, p)
}

func (s *PurchaseService) Update(ctx context.Context, p *domain.PlannedPurchase) error {
	return s.repo.Update(ctx, p)
}

func (s *PurchaseService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
