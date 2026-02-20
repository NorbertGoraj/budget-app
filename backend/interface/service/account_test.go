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

func TestAccountService_GetAll(t *testing.T) {
	repo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewAccount(repo)
	ctx := context.Background()

	expected := []domain.Account{{ID: 1, Name: "Checking", Balance: 1000.0}}
	repo.EXPECT().GetAll(ctx).Return(expected, nil)

	got, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestAccountService_GetAll_Error(t *testing.T) {
	repo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewAccount(repo)
	ctx := context.Background()

	repoErr := errors.New("db error")
	repo.EXPECT().GetAll(ctx).Return(nil, repoErr)

	_, err := svc.GetAll(ctx)
	assert.ErrorIs(t, err, repoErr)
}

func TestAccountService_Create(t *testing.T) {
	repo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewAccount(repo)
	ctx := context.Background()

	acc := &domain.Account{Name: "Savings", Type: "savings"}
	repo.EXPECT().Create(ctx, acc).Return(nil)

	err := svc.Create(ctx, acc)
	assert.NoError(t, err)
}

func TestAccountService_Update(t *testing.T) {
	repo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewAccount(repo)
	ctx := context.Background()

	acc := &domain.Account{ID: 1, Name: "Updated"}
	repo.EXPECT().Update(ctx, acc).Return(nil)

	err := svc.Update(ctx, acc)
	assert.NoError(t, err)
}

func TestAccountService_Delete(t *testing.T) {
	repo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewAccount(repo)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, 5).Return(nil)

	err := svc.Delete(ctx, 5)
	assert.NoError(t, err)
}
