package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"budget-app/db"
	"budget-app/models"

	"github.com/gin-gonic/gin"
)

type BudgetStatus struct {
	Category  string  `json:"category"`
	Limit     float64 `json:"limit"`
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
}

type PurchaseAffordability struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Cost           float64 `json:"cost"`
	TargetMonth    string  `json:"target_month"`
	Priority       string  `json:"priority"`
	Affordable     bool    `json:"affordable"`
	SuggestedMonth string  `json:"suggested_month,omitempty"`
	Reason         string  `json:"reason,omitempty"`
}

type InvestmentSummary struct {
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Frequency string  `json:"frequency"`
	Status    string  `json:"status"`
	Category  string  `json:"category"`
}

type DashboardResponse struct {
	TotalBalance           float64                 `json:"total_balance"`
	Accounts               []models.Account        `json:"accounts"`
	MonthlyIncome          float64                 `json:"monthly_income"`
	MonthlyExpenses        float64                 `json:"monthly_expenses"`
	MonthlyInvestmentTotal float64                 `json:"monthly_investment_total"`
	AvailableForInvestment float64                 `json:"available_for_investments"`
	BudgetStatus           []BudgetStatus          `json:"budget_status"`
	PlannedPurchases       []PurchaseAffordability  `json:"planned_purchases"`
	Investments            []InvestmentSummary     `json:"investments"`
}

func GetDashboard(c *gin.Context) {
	now := time.Now()
	currentMonth := now.Format("2006-01")

	accounts, err := models.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalBalance float64
	for _, a := range accounts {
		totalBalance += a.Balance
	}

	var monthlyIncome, monthlyExpenses float64
	err = db.Pool.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(CASE WHEN type='income' THEN amount ELSE 0 END), 0), COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0) FROM transactions WHERE TO_CHAR(date, 'YYYY-MM') = $1",
		currentMonth).Scan(&monthlyIncome, &monthlyExpenses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	budgets, err := models.GetAllBudgets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var budgetStatuses []BudgetStatus
	for _, b := range budgets {
		var spent float64
		db.Pool.QueryRow(context.Background(),
			"SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE category=$1 AND type='expense' AND TO_CHAR(date, 'YYYY-MM')=$2",
			b.Category, currentMonth).Scan(&spent)
		budgetStatuses = append(budgetStatuses, BudgetStatus{
			Category:  b.Category,
			Limit:     b.MonthlyLimit,
			Spent:     spent,
			Remaining: b.MonthlyLimit - spent,
		})
	}
	if budgetStatuses == nil {
		budgetStatuses = []BudgetStatus{}
	}

	monthlyInvTotal, err := models.GetMonthlyInvestmentTotal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	availableForInvestment := monthlyIncome - monthlyExpenses - monthlyInvTotal

	purchases, err := models.GetAllPurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var purchaseAffordability []PurchaseAffordability
	for _, p := range purchases {
		if p.Status == "purchased" || p.Status == "cancelled" {
			continue
		}
		monthlySurplus := monthlyIncome - monthlyExpenses - monthlyInvTotal
		affordable := totalBalance >= p.EstimatedCost && monthlySurplus >= 0

		pa := PurchaseAffordability{
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
				suggested := now.AddDate(0, monthsNeeded, 0)
				pa.SuggestedMonth = suggested.Format("2006-01")
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
	if purchaseAffordability == nil {
		purchaseAffordability = []PurchaseAffordability{}
	}

	investments, err := models.GetAllInvestments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var investmentSummaries []InvestmentSummary
	for _, inv := range investments {
		if inv.Status != "active" {
			continue
		}
		investmentSummaries = append(investmentSummaries, InvestmentSummary{
			Name:      inv.Name,
			Amount:    inv.Amount,
			Frequency: inv.Frequency,
			Status:    inv.Status,
			Category:  inv.Category,
		})
	}
	if investmentSummaries == nil {
		investmentSummaries = []InvestmentSummary{}
	}

	if accounts == nil {
		accounts = []models.Account{}
	}

	c.JSON(http.StatusOK, DashboardResponse{
		TotalBalance:           totalBalance,
		Accounts:               accounts,
		MonthlyIncome:          monthlyIncome,
		MonthlyExpenses:        monthlyExpenses,
		MonthlyInvestmentTotal: monthlyInvTotal,
		AvailableForInvestment: availableForInvestment,
		BudgetStatus:           budgetStatuses,
		PlannedPurchases:       purchaseAffordability,
		Investments:            investmentSummaries,
	})
}
