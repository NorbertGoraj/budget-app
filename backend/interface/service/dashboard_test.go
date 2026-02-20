package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"budget-app/domain"
	"budget-app/interface/service"
	mockdomain "budget-app/mocks/domain"
)

func newDashboardSvc(t *testing.T) (
	*service.DashboardService,
	*mockdomain.MockAccountRepository,
	*mockdomain.MockTransactionRepository,
	*mockdomain.MockBudgetRepository,
	*mockdomain.MockPurchaseRepository,
	*mockdomain.MockInvestmentRepository,
) {
	t.Helper()
	accRepo := mockdomain.NewMockAccountRepository(t)
	txRepo := mockdomain.NewMockTransactionRepository(t)
	budgetRepo := mockdomain.NewMockBudgetRepository(t)
	purchaseRepo := mockdomain.NewMockPurchaseRepository(t)
	invRepo := mockdomain.NewMockInvestmentRepository(t)
	svc := service.NewDashboard(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo)
	return svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo
}

// setupStage1 sets up the common stage-1 expectations. All context args use
// mock.Anything because errgroup wraps the caller's context.
func setupStage1(
	accRepo *mockdomain.MockAccountRepository,
	txRepo *mockdomain.MockTransactionRepository,
	budgetRepo *mockdomain.MockBudgetRepository,
	purchaseRepo *mockdomain.MockPurchaseRepository,
	invRepo *mockdomain.MockInvestmentRepository,
	accounts []domain.Account,
	income, expenses float64,
	budgets []domain.Budget,
	monthlyInvTotal float64,
	purchases []domain.PlannedPurchase,
	investments []domain.Investment,
) {
	accRepo.EXPECT().GetAll(mock.Anything).Return(accounts, nil)
	txRepo.EXPECT().MonthlySums(mock.Anything, mock.AnythingOfType("string")).Return(income, expenses, nil)
	budgetRepo.EXPECT().GetAll(mock.Anything).Return(budgets, nil)
	invRepo.EXPECT().MonthlyTotal(mock.Anything).Return(monthlyInvTotal, nil)
	purchaseRepo.EXPECT().GetAll(mock.Anything).Return(purchases, nil)
	invRepo.EXPECT().GetAll(mock.Anything).Return(investments, nil)
}

func TestDashboardService_TotalBalance(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	accounts := []domain.Account{
		{ID: 1, Balance: 1000.0},
		{ID: 2, Balance: 500.0},
	}
	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		accounts, 0, 0, nil, 0, nil, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1500.0, resp.TotalBalance)
}

func TestDashboardService_MonthlySurplus(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		nil, 3000.0, 1200.0, nil, 300.0, nil, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	// surplus = income - expenses - investments = 3000 - 1200 - 300 = 1500
	assert.Equal(t, 1500.0, resp.AvailableForInvestment)
	assert.Equal(t, 3000.0, resp.MonthlyIncome)
	assert.Equal(t, 1200.0, resp.MonthlyExpenses)
	assert.Equal(t, 300.0, resp.MonthlyInvestmentTotal)
}

func TestDashboardService_BudgetStatus(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	budgets := []domain.Budget{
		{ID: 1, Category: "Food", MonthlyLimit: 500.0},
		{ID: 2, Category: "Transport", MonthlyLimit: 200.0},
	}
	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		nil, 0, 0, budgets, 0, nil, nil)

	txRepo.EXPECT().SpentByCategory(mock.Anything, "Food", mock.AnythingOfType("string")).Return(350.0, nil)
	txRepo.EXPECT().SpentByCategory(mock.Anything, "Transport", mock.AnythingOfType("string")).Return(80.0, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	require.Len(t, resp.BudgetStatus, 2)

	byCategory := make(map[string]domain.BudgetStatus)
	for _, bs := range resp.BudgetStatus {
		byCategory[bs.Category] = bs
	}

	food := byCategory["Food"]
	assert.Equal(t, 500.0, food.Limit)
	assert.Equal(t, 350.0, food.Spent)
	assert.Equal(t, 150.0, food.Remaining)

	transport := byCategory["Transport"]
	assert.Equal(t, 200.0, transport.Limit)
	assert.Equal(t, 80.0, transport.Spent)
	assert.Equal(t, 120.0, transport.Remaining)
}

func TestDashboardService_PurchaseAffordability_Affordable(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	accounts := []domain.Account{{Balance: 5000.0}}
	purchases := []domain.PlannedPurchase{
		{ID: 1, Name: "Laptop", EstimatedCost: 1500.0, Status: "pending", Priority: "high"},
	}
	// income=3000, expenses=1000, invTotal=0 → surplus=2000 ≥ 0; balance=5000 ≥ 1500
	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		accounts, 3000.0, 1000.0, nil, 0, purchases, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	require.Len(t, resp.PlannedPurchases, 1)

	pa := resp.PlannedPurchases[0]
	assert.True(t, pa.Affordable)
	assert.Equal(t, "Laptop", pa.Name)
	// reason should mention balance remaining after purchase: 5000 - 1500 = 3500
	assert.Contains(t, pa.Reason, "3500.00")
}

func TestDashboardService_PurchaseAffordability_NotAffordable_NegativeSurplus(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	accounts := []domain.Account{{Balance: 10000.0}}
	purchases := []domain.PlannedPurchase{
		{ID: 2, Name: "Car", EstimatedCost: 500.0, Status: "pending"},
	}
	// surplus = 1000 - 2000 - 0 = -1000 (negative)
	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		accounts, 1000.0, 2000.0, nil, 0, purchases, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	require.Len(t, resp.PlannedPurchases, 1)

	pa := resp.PlannedPurchases[0]
	assert.False(t, pa.Affordable)
	assert.Contains(t, pa.Reason, "surplus is zero or negative")
}

func TestDashboardService_PurchaseAffordability_NotAffordable_InsufficientBalance(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	accounts := []domain.Account{{Balance: 100.0}}
	purchases := []domain.PlannedPurchase{
		{ID: 3, Name: "Vacation", EstimatedCost: 2000.0, Status: "pending"},
	}
	// surplus = 3000 - 1000 - 0 = 2000 (positive), but balance=100 < 2000
	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		accounts, 3000.0, 1000.0, nil, 0, purchases, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	require.Len(t, resp.PlannedPurchases, 1)

	pa := resp.PlannedPurchases[0]
	assert.False(t, pa.Affordable)
	// needed = 2000 - 100 = 1900; suggested month should be provided
	assert.NotEmpty(t, pa.SuggestedMonth)
	assert.Contains(t, pa.Reason, "1900.00")
}

func TestDashboardService_PurchaseAffordability_SkipPurchasedAndCancelled(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	purchases := []domain.PlannedPurchase{
		{ID: 1, Name: "Purchased Item", Status: "purchased"},
		{ID: 2, Name: "Cancelled Item", Status: "cancelled"},
		{ID: 3, Name: "Pending Item", EstimatedCost: 50.0, Status: "pending"},
	}
	accounts := []domain.Account{{Balance: 1000.0}}
	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		accounts, 2000.0, 500.0, nil, 0, purchases, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	// only the pending item should appear
	assert.Len(t, resp.PlannedPurchases, 1)
	assert.Equal(t, "Pending Item", resp.PlannedPurchases[0].Name)
}

func TestDashboardService_InvestmentSummaries_OnlyActive(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	investments := []domain.Investment{
		{ID: 1, Name: "ETF", Status: "active", Amount: 500.0, Frequency: "monthly", Category: "stocks"},
		{ID: 2, Name: "Old Bond", Status: "closed", Amount: 1000.0},
		{ID: 3, Name: "Crypto", Status: "active", Amount: 200.0, Frequency: "weekly", Category: "crypto"},
	}
	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		nil, 0, 0, nil, 0, nil, investments)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	require.Len(t, resp.Investments, 2)

	names := []string{resp.Investments[0].Name, resp.Investments[1].Name}
	assert.Contains(t, names, "ETF")
	assert.Contains(t, names, "Crypto")
	assert.NotContains(t, names, "Old Bond")
}

func TestDashboardService_EmptyAccounts_ReturnsEmptySlice(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	setupStage1(accRepo, txRepo, budgetRepo, purchaseRepo, invRepo,
		nil, 0, 0, nil, 0, nil, nil)

	resp, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, resp.Accounts)
	assert.Empty(t, resp.Accounts)
}

func TestDashboardService_ReturnsErrorOnStage1Failure(t *testing.T) {
	svc, accRepo, txRepo, budgetRepo, purchaseRepo, invRepo := newDashboardSvc(t)

	dbErr := errors.New("database unavailable")
	accRepo.EXPECT().GetAll(mock.Anything).Return(nil, dbErr)
	txRepo.EXPECT().MonthlySums(mock.Anything, mock.AnythingOfType("string")).Return(0.0, 0.0, nil).Maybe()
	budgetRepo.EXPECT().GetAll(mock.Anything).Return(nil, nil).Maybe()
	invRepo.EXPECT().MonthlyTotal(mock.Anything).Return(0.0, nil).Maybe()
	purchaseRepo.EXPECT().GetAll(mock.Anything).Return(nil, nil).Maybe()
	invRepo.EXPECT().GetAll(mock.Anything).Return(nil, nil).Maybe()

	_, err := svc.GetDashboard(context.Background())
	assert.Error(t, err)
}
