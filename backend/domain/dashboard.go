package domain

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

type DebtItem struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	CurrentBalance float64 `json:"current_balance"`
	InterestRate   float64 `json:"interest_rate"`
	MinimumPayment float64 `json:"minimum_payment"`
	DueDay         int     `json:"due_day"`
}

type DebtSummary struct {
	TotalDebt              float64    `json:"total_debt"`
	MonthlyMinimumPayments float64    `json:"monthly_minimum_payments"`
	ActiveDebts            []DebtItem `json:"active_debts"`
}

type DashboardResponse struct {
	TotalBalance           float64                 `json:"total_balance"`
	Accounts               []Account               `json:"accounts"`
	MonthlyIncome          float64                 `json:"monthly_income"`
	MonthlyExpenses        float64                 `json:"monthly_expenses"`
	MonthlyInvestmentTotal float64                 `json:"monthly_investment_total"`
	AvailableForInvestment float64                 `json:"available_for_investments"`
	BudgetStatus           []BudgetStatus          `json:"budget_status"`
	PlannedPurchases       []PurchaseAffordability `json:"planned_purchases"`
	Investments            []InvestmentSummary     `json:"investments"`
	DebtSummary            DebtSummary             `json:"debt_summary"`
}
