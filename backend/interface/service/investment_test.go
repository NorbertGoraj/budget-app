package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"budget-app/domain"
	"budget-app/interface/service"
	mockdomain "budget-app/mocks/domain"
)

func TestInvestmentService_GetAll(t *testing.T) {
	repo := mockdomain.NewMockInvestmentRepository(t)
	svc := service.NewInvestment(repo)
	ctx := context.Background()

	expected := []domain.Investment{{ID: 1, Name: "ETF", Status: "active", Amount: 500.0}}
	repo.EXPECT().GetAll(ctx).Return(expected, nil)

	got, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestInvestmentService_Create(t *testing.T) {
	repo := mockdomain.NewMockInvestmentRepository(t)
	svc := service.NewInvestment(repo)
	ctx := context.Background()

	inv := &domain.Investment{Name: "Bonds", Amount: 1000.0, Status: "active"}
	repo.EXPECT().Create(ctx, inv).Return(nil)

	err := svc.Create(ctx, inv)
	assert.NoError(t, err)
}

func TestInvestmentService_Update(t *testing.T) {
	repo := mockdomain.NewMockInvestmentRepository(t)
	svc := service.NewInvestment(repo)
	ctx := context.Background()

	inv := &domain.Investment{ID: 2, Amount: 1500.0}
	repo.EXPECT().Update(ctx, inv).Return(nil)

	err := svc.Update(ctx, inv)
	assert.NoError(t, err)
}

func TestInvestmentService_Delete(t *testing.T) {
	repo := mockdomain.NewMockInvestmentRepository(t)
	svc := service.NewInvestment(repo)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, 7).Return(nil)

	err := svc.Delete(ctx, 7)
	assert.NoError(t, err)
}
