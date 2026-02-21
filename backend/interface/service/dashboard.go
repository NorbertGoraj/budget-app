package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"budget-app/domain"
)

type DashboardService struct {
	accountRepo     domain.AccountRepository
	transactionRepo domain.TransactionRepository
	budgetRepo      domain.BudgetRepository
	purchaseRepo    domain.PurchaseRepository
	investmentRepo  domain.InvestmentRepository
	debtRepo        domain.DebtRepository
}

func NewDashboard(
	accountRepo domain.AccountRepository,
	transactionRepo domain.TransactionRepository,
	budgetRepo domain.BudgetRepository,
	purchaseRepo domain.PurchaseRepository,
	investmentRepo domain.InvestmentRepository,
	debtRepo domain.DebtRepository,
) *DashboardService {
	return &DashboardService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
		purchaseRepo:    purchaseRepo,
		investmentRepo:  investmentRepo,
		debtRepo:        debtRepo,
	}
}

func (s *DashboardService) GetDashboard(ctx context.Context) (*domain.DashboardResponse, error) {
	now := time.Now()
	currentMonth := now.Format("2006-01")

	// ── Stage 1 ────────────────────────────────────────────────────────────────
	// All data sources are independent — fetch them in parallel.
	// errgroup inherits the request context: if the client disconnects, all
	// in-flight DB queries are cancelled automatically.
	var (
		accounts        []domain.Account
		monthlyIncome   float64
		monthlyExpenses float64
		budgets         []domain.Budget
		purchases       []domain.PlannedPurchase
		investments     []domain.Investment
		debts           []domain.Debt
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		accounts, err = s.accountRepo.GetAll(gCtx)
		return err
	})
	g.Go(func() error {
		var err error
		monthlyIncome, monthlyExpenses, err = s.transactionRepo.MonthlySums(gCtx, currentMonth)
		return err
	})
	g.Go(func() error {
		var err error
		budgets, err = s.budgetRepo.GetAll(gCtx)
		return err
	})
	g.Go(func() error {
		var err error
		purchases, err = s.purchaseRepo.GetAll(gCtx)
		return err
	})
	g.Go(func() error {
		var err error
		investments, err = s.investmentRepo.GetAll(gCtx)
		return err
	})
	g.Go(func() error {
		var err error
		debts, err = s.debtRepo.GetAll(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// ── Stage 2 ────────────────────────────────────────────────────────────────
	// Each budget needs one SpentByCategory query — run them in parallel.
	// Pre-allocate the slice so each goroutine writes to its own index (no mutex needed).
	budgetStatuses := make([]domain.BudgetStatus, len(budgets))

	g2, g2Ctx := errgroup.WithContext(ctx)
	for i, b := range budgets {
		i, b := i, b
		g2.Go(func() error {
			spent, err := s.transactionRepo.SpentByCategory(g2Ctx, b.Category, currentMonth)
			if err != nil {
				return err
			}
			budgetStatuses[i] = domain.BudgetStatus{
				Category:  b.Category,
				Limit:     b.MonthlyLimit,
				Spent:     spent,
				Remaining: b.MonthlyLimit - spent,
			}
			return nil
		})
	}

	if err := g2.Wait(); err != nil {
		return nil, err
	}

	// ── Stage 3 ────────────────────────────────────────────────────────────────
	// Pure computation — no IO, no goroutines needed.
	if accounts == nil {
		accounts = []domain.Account{}
	}

	var totalBalance float64
	for _, a := range accounts {
		totalBalance += a.Balance
	}

	var monthlyInvTotal float64
	for _, inv := range investments {
		if inv.Status != "active" {
			continue
		}
		switch inv.Frequency {
		case "weekly":
			monthlyInvTotal += inv.Amount * 4.33
		case "monthly":
			monthlyInvTotal += inv.Amount
		case "quarterly":
			monthlyInvTotal += inv.Amount / 3
		case "yearly":
			monthlyInvTotal += inv.Amount / 12
		}
	}

	monthlySurplus := monthlyIncome - monthlyExpenses - monthlyInvTotal

	purchaseAffordability := make([]domain.PurchaseAffordability, 0, len(purchases))
	for _, p := range purchases {
		if p.Status == "purchased" || p.Status == "cancelled" {
			continue
		}

		affordable := totalBalance >= p.EstimatedCost && monthlySurplus >= 0
		pa := domain.PurchaseAffordability{
			ID:          p.ID,
			Name:        p.Name,
			Cost:        p.EstimatedCost,
			TargetMonth: p.TargetMonth,
			Priority:    p.Priority,
			Affordable:  affordable,
		}

		if !affordable {
			needed := p.EstimatedCost - totalBalance
			if needed > 0 && monthlySurplus > 0 {
				monthsNeeded := int(needed/monthlySurplus) + 1
				pa.SuggestedMonth = now.AddDate(0, monthsNeeded, 0).Format("2006-01")
				pa.Reason = fmt.Sprintf("Need %.2f more. With monthly surplus of %.2f, affordable by %s.",
					needed, monthlySurplus, pa.SuggestedMonth)
			} else if monthlySurplus <= 0 {
				pa.Reason = "Monthly surplus is zero or negative. Reduce expenses or increase income first."
			}
		} else {
			pa.Reason = fmt.Sprintf("Balance after: %.2f", totalBalance-p.EstimatedCost)
		}

		purchaseAffordability = append(purchaseAffordability, pa)
	}

	investmentSummaries := make([]domain.InvestmentSummary, 0, len(investments))
	var debtTotalBal, debtMonthlyMin float64
	activeDebtItems := make([]domain.DebtItem, 0, len(debts))

	for _, inv := range investments {
		if inv.Status != "active" {
			continue
		}
		investmentSummaries = append(investmentSummaries, domain.InvestmentSummary{
			Name:      inv.Name,
			Amount:    inv.Amount,
			Frequency: inv.Frequency,
			Status:    inv.Status,
			Category:  inv.Category,
		})
	}

	for _, d := range debts {
		if d.Status != "active" {
			continue
		}
		debtTotalBal += d.CurrentBalance
		debtMonthlyMin += d.MinimumPayment
		activeDebtItems = append(activeDebtItems, domain.DebtItem{
			ID:             d.ID,
			Name:           d.Name,
			Type:           d.Type,
			CurrentBalance: d.CurrentBalance,
			InterestRate:   d.InterestRate,
			MinimumPayment: d.MinimumPayment,
			DueDay:         d.DueDay,
		})
	}

	return &domain.DashboardResponse{
		TotalBalance:           totalBalance,
		Accounts:               accounts,
		MonthlyIncome:          monthlyIncome,
		MonthlyExpenses:        monthlyExpenses,
		MonthlyInvestmentTotal: monthlyInvTotal,
		AvailableForInvestment: monthlySurplus,
		BudgetStatus:           budgetStatuses,
		PlannedPurchases:       purchaseAffordability,
		Investments:            investmentSummaries,
		DebtSummary: domain.DebtSummary{
			TotalDebt:              debtTotalBal,
			MonthlyMinimumPayments: debtMonthlyMin,
			ActiveDebts:            activeDebtItems,
		},
	}, nil
}
