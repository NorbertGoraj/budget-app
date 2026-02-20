package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"budget-app/domain"
	"budget-app/interface/service"
	mockdomain "budget-app/mocks/domain"
)

func TestPurchaseService_GetAll(t *testing.T) {
	repo := mockdomain.NewMockPurchaseRepository(t)
	svc := service.NewPurchase(repo)
	ctx := context.Background()

	expected := []domain.PlannedPurchase{{ID: 1, Name: "Laptop", EstimatedCost: 1500.0}}
	repo.EXPECT().GetAll(ctx).Return(expected, nil)

	got, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestPurchaseService_Create(t *testing.T) {
	repo := mockdomain.NewMockPurchaseRepository(t)
	svc := service.NewPurchase(repo)
	ctx := context.Background()

	p := &domain.PlannedPurchase{Name: "Phone", EstimatedCost: 800.0, Status: "pending"}
	repo.EXPECT().Create(ctx, p).Return(nil)

	err := svc.Create(ctx, p)
	assert.NoError(t, err)
}

func TestPurchaseService_Update(t *testing.T) {
	repo := mockdomain.NewMockPurchaseRepository(t)
	svc := service.NewPurchase(repo)
	ctx := context.Background()

	p := &domain.PlannedPurchase{ID: 3, Status: "purchased"}
	repo.EXPECT().Update(ctx, p).Return(nil)

	err := svc.Update(ctx, p)
	assert.NoError(t, err)
}

func TestPurchaseService_Delete(t *testing.T) {
	repo := mockdomain.NewMockPurchaseRepository(t)
	svc := service.NewPurchase(repo)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, 4).Return(nil)

	err := svc.Delete(ctx, 4)
	assert.NoError(t, err)
}
