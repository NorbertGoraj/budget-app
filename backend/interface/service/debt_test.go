package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"budget-app/domain"
	"budget-app/interface/service"
	mockdomain "budget-app/mocks/domain"
)

func newDebtSvc(t *testing.T) (
	*service.DebtService,
	*mockdomain.MockDebtRepository,
	*mockdomain.MockDebtPaymentRepository,
) {
	t.Helper()
	debtRepo := mockdomain.NewMockDebtRepository(t)
	paymentRepo := mockdomain.NewMockDebtPaymentRepository(t)
	svc := service.NewDebt(debtRepo, paymentRepo)
	return svc, debtRepo, paymentRepo
}

func TestDebtService_RecordPayment_ReducesBalance(t *testing.T) {
	svc, debtRepo, paymentRepo := newDebtSvc(t)

	debt := &domain.Debt{ID: 1, Name: "Credit Card", CurrentBalance: 1000.0, Status: "active"}
	payment := &domain.DebtPayment{Amount: 200.0, PaidAt: "2026-02-01"}

	debtRepo.EXPECT().GetByID(context.Background(), 1).Return(debt, nil)
	paymentRepo.EXPECT().Create(context.Background(), payment).Return(nil)
	debtRepo.EXPECT().Update(context.Background(), &domain.Debt{
		ID: 1, Name: "Credit Card", CurrentBalance: 800.0, Status: "active",
	}).Return(nil)

	err := svc.RecordPayment(context.Background(), 1, payment)
	require.NoError(t, err)
	assert.Equal(t, 800.0, debt.CurrentBalance)
	assert.Equal(t, "active", debt.Status)
	assert.Equal(t, 1, payment.DebtID)
}

func TestDebtService_RecordPayment_PaysOffDebt(t *testing.T) {
	svc, debtRepo, paymentRepo := newDebtSvc(t)

	debt := &domain.Debt{ID: 2, Name: "Car Loan", CurrentBalance: 500.0, Status: "active"}
	payment := &domain.DebtPayment{Amount: 600.0, PaidAt: "2026-02-01"}

	debtRepo.EXPECT().GetByID(context.Background(), 2).Return(debt, nil)
	paymentRepo.EXPECT().Create(context.Background(), payment).Return(nil)
	debtRepo.EXPECT().Update(context.Background(), &domain.Debt{
		ID: 2, Name: "Car Loan", CurrentBalance: 0.0, Status: "paid_off",
	}).Return(nil)

	err := svc.RecordPayment(context.Background(), 2, payment)
	require.NoError(t, err)
	assert.Equal(t, 0.0, debt.CurrentBalance)
	assert.Equal(t, "paid_off", debt.Status)
}

func TestDebtService_RecordPayment_ExactPayoffClampsToZero(t *testing.T) {
	svc, debtRepo, paymentRepo := newDebtSvc(t)

	debt := &domain.Debt{ID: 3, Name: "Student Loan", CurrentBalance: 300.0, Status: "active"}
	payment := &domain.DebtPayment{Amount: 300.0, PaidAt: "2026-02-01"}

	debtRepo.EXPECT().GetByID(context.Background(), 3).Return(debt, nil)
	paymentRepo.EXPECT().Create(context.Background(), payment).Return(nil)
	debtRepo.EXPECT().Update(context.Background(), &domain.Debt{
		ID: 3, Name: "Student Loan", CurrentBalance: 0.0, Status: "paid_off",
	}).Return(nil)

	err := svc.RecordPayment(context.Background(), 3, payment)
	require.NoError(t, err)
	assert.Equal(t, 0.0, debt.CurrentBalance)
	assert.Equal(t, "paid_off", debt.Status)
}

func TestDebtService_DeletePayment_ReversesBalance(t *testing.T) {
	svc, debtRepo, paymentRepo := newDebtSvc(t)

	payment := &domain.DebtPayment{ID: 10, DebtID: 1, Amount: 200.0}
	debt := &domain.Debt{ID: 1, Name: "Credit Card", CurrentBalance: 800.0, Status: "active"}

	paymentRepo.EXPECT().GetByID(context.Background(), 10).Return(payment, nil)
	debtRepo.EXPECT().GetByID(context.Background(), 1).Return(debt, nil)
	debtRepo.EXPECT().Update(context.Background(), &domain.Debt{
		ID: 1, Name: "Credit Card", CurrentBalance: 1000.0, Status: "active",
	}).Return(nil)
	paymentRepo.EXPECT().Delete(context.Background(), 10).Return(nil)

	err := svc.DeletePayment(context.Background(), 10)
	require.NoError(t, err)
	assert.Equal(t, 1000.0, debt.CurrentBalance)
	assert.Equal(t, "active", debt.Status)
}

func TestDebtService_DeletePayment_ReactivatesPaidOffDebt(t *testing.T) {
	svc, debtRepo, paymentRepo := newDebtSvc(t)

	payment := &domain.DebtPayment{ID: 11, DebtID: 2, Amount: 500.0}
	debt := &domain.Debt{ID: 2, Name: "Car Loan", CurrentBalance: 0.0, Status: "paid_off"}

	paymentRepo.EXPECT().GetByID(context.Background(), 11).Return(payment, nil)
	debtRepo.EXPECT().GetByID(context.Background(), 2).Return(debt, nil)
	debtRepo.EXPECT().Update(context.Background(), &domain.Debt{
		ID: 2, Name: "Car Loan", CurrentBalance: 500.0, Status: "active",
	}).Return(nil)
	paymentRepo.EXPECT().Delete(context.Background(), 11).Return(nil)

	err := svc.DeletePayment(context.Background(), 11)
	require.NoError(t, err)
	assert.Equal(t, 500.0, debt.CurrentBalance)
	assert.Equal(t, "active", debt.Status)
}
