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

func TestTransactionService_Create_Income(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)

	tx := &domain.Transaction{AccountID: 1, Amount: 200.0, Type: "income"}
	ctx := context.Background()

	txRepo.EXPECT().Create(ctx, tx).Return(nil)
	accRepo.EXPECT().UpdateBalance(ctx, 1, 200.0).Return(nil)

	err := svc.Create(ctx, tx)
	assert.NoError(t, err)
}

func TestTransactionService_Create_Expense(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)

	tx := &domain.Transaction{AccountID: 2, Amount: 50.0, Type: "expense"}
	ctx := context.Background()

	txRepo.EXPECT().Create(ctx, tx).Return(nil)
	accRepo.EXPECT().UpdateBalance(ctx, 2, -50.0).Return(nil)

	err := svc.Create(ctx, tx)
	assert.NoError(t, err)
}

func TestTransactionService_Create_RepoError(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)

	tx := &domain.Transaction{AccountID: 1, Amount: 100.0, Type: "income"}
	ctx := context.Background()

	repoErr := errors.New("db error")
	txRepo.EXPECT().Create(ctx, tx).Return(repoErr)

	err := svc.Create(ctx, tx)
	assert.ErrorIs(t, err, repoErr)
}

func TestTransactionService_Delete_Income(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	existing := &domain.Transaction{ID: 10, AccountID: 3, Amount: 300.0, Type: "income"}
	txRepo.EXPECT().GetByID(ctx, 10).Return(existing, nil)
	// income reversal: subtract the amount that was added
	accRepo.EXPECT().UpdateBalance(ctx, 3, -300.0).Return(nil)
	txRepo.EXPECT().Delete(ctx, 10).Return(nil)

	err := svc.Delete(ctx, 10)
	assert.NoError(t, err)
}

func TestTransactionService_Delete_Expense(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	existing := &domain.Transaction{ID: 11, AccountID: 4, Amount: 75.0, Type: "expense"}
	txRepo.EXPECT().GetByID(ctx, 11).Return(existing, nil)
	// expense reversal: add back the amount that was subtracted
	accRepo.EXPECT().UpdateBalance(ctx, 4, 75.0).Return(nil)
	txRepo.EXPECT().Delete(ctx, 11).Return(nil)

	err := svc.Delete(ctx, 11)
	assert.NoError(t, err)
}

func TestTransactionService_Delete_GetByIDError(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	repoErr := errors.New("not found")
	txRepo.EXPECT().GetByID(ctx, 99).Return(nil, repoErr)

	err := svc.Delete(ctx, 99)
	assert.ErrorIs(t, err, repoErr)
}

func TestTransactionService_Delete_UpdateBalanceError(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	existing := &domain.Transaction{ID: 12, AccountID: 5, Amount: 100.0, Type: "income"}
	txRepo.EXPECT().GetByID(ctx, 12).Return(existing, nil)

	balanceErr := errors.New("balance update failed")
	accRepo.EXPECT().UpdateBalance(ctx, 5, -100.0).Return(balanceErr)

	err := svc.Delete(ctx, 12)
	assert.ErrorIs(t, err, balanceErr)
}

func TestTransactionService_GetAll(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	filter := domain.TransactionFilter{Month: "2026-01"}
	expected := []domain.Transaction{{ID: 1, Amount: 100.0}}
	txRepo.EXPECT().GetAll(ctx, filter).Return(expected, nil)

	got, err := svc.GetAll(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestTransactionService_Exists(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	txRepo.EXPECT().Exists(ctx, "2026-01-15", "Groceries", 42.5).Return(true, nil)

	exists, err := svc.Exists(ctx, "2026-01-15", "Groceries", 42.5)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestTransactionService_Update(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	tx := &domain.Transaction{ID: 1, Amount: 100.0}
	txRepo.EXPECT().Update(ctx, tx).Return(nil)

	err := svc.Update(ctx, tx)
	assert.NoError(t, err)
}

func TestTransactionService_Create_UpdateBalanceError(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	tx := &domain.Transaction{AccountID: 1, Amount: 100.0, Type: "income"}
	txRepo.EXPECT().Create(ctx, tx).Return(nil)

	balanceErr := errors.New("balance update failed")
	accRepo.EXPECT().UpdateBalance(ctx, 1, 100.0).Return(balanceErr)

	err := svc.Create(ctx, tx)
	assert.ErrorIs(t, err, balanceErr)
}

func TestTransactionService_Delete_DeleteRepoError(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	existing := &domain.Transaction{ID: 20, AccountID: 6, Amount: 50.0, Type: "expense"}
	txRepo.EXPECT().GetByID(ctx, 20).Return(existing, nil)
	accRepo.EXPECT().UpdateBalance(ctx, 6, 50.0).Return(nil)

	deleteErr := errors.New("delete failed")
	txRepo.EXPECT().Delete(ctx, 20).Return(deleteErr)

	err := svc.Delete(ctx, 20)
	assert.ErrorIs(t, err, deleteErr)
}

// Verify accRepo is never called when txRepo.Create fails.
func TestTransactionService_Create_NoBalanceUpdateOnError(t *testing.T) {
	txRepo := mockdomain.NewMockTransactionRepository(t)
	accRepo := mockdomain.NewMockAccountRepository(t)
	svc := service.NewTransaction(txRepo, accRepo)
	ctx := context.Background()

	tx := &domain.Transaction{AccountID: 1, Amount: 100.0, Type: "income"}
	txRepo.EXPECT().Create(ctx, tx).Return(errors.New("insert failed"))
	// accRepo.UpdateBalance must NOT be called — testify/mock asserts this automatically
	// because no expectation is set and mock.AssertExpectations runs via t.Cleanup.

	err := svc.Create(ctx, tx)
	assert.Error(t, err)
}
