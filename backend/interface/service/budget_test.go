package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"budget-app/domain"
	"budget-app/interface/service"
	mockdomain "budget-app/mocks/domain"
)

func TestBudgetService_GetAll(t *testing.T) {
	repo := mockdomain.NewMockBudgetRepository(t)
	svc := service.NewBudget(repo)
	ctx := context.Background()

	expected := []domain.Budget{{ID: 1, Category: "Food", MonthlyLimit: 500.0}}
	repo.EXPECT().GetAll(ctx).Return(expected, nil)

	got, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestBudgetService_Create(t *testing.T) {
	repo := mockdomain.NewMockBudgetRepository(t)
	svc := service.NewBudget(repo)
	ctx := context.Background()

	b := &domain.Budget{Category: "Entertainment", MonthlyLimit: 100.0}
	repo.EXPECT().Create(ctx, b).Return(nil)

	err := svc.Create(ctx, b)
	assert.NoError(t, err)
}

func TestBudgetService_Update(t *testing.T) {
	repo := mockdomain.NewMockBudgetRepository(t)
	svc := service.NewBudget(repo)
	ctx := context.Background()

	b := &domain.Budget{ID: 2, MonthlyLimit: 200.0}
	repo.EXPECT().Update(ctx, b).Return(nil)

	err := svc.Update(ctx, b)
	assert.NoError(t, err)
}

func TestBudgetService_Delete(t *testing.T) {
	repo := mockdomain.NewMockBudgetRepository(t)
	svc := service.NewBudget(repo)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, 3).Return(nil)

	err := svc.Delete(ctx, 3)
	assert.NoError(t, err)
}

func TestBudgetService_Delete_Error(t *testing.T) {
	repo := mockdomain.NewMockBudgetRepository(t)
	svc := service.NewBudget(repo)
	ctx := context.Background()

	repoErr := errors.New("delete failed")
	repo.EXPECT().Delete(ctx, 99).Return(repoErr)

	err := svc.Delete(ctx, 99)
	assert.ErrorIs(t, err, repoErr)
}
