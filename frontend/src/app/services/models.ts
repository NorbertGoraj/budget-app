export interface Account {
  id: number;
  name: string;
  type: 'bank' | 'cash';
  balance: number;
  created_at: string;
}

export interface Transaction {
  id: number;
  account_id: number;
  amount: number;
  description: string;
  category: string;
  type: 'income' | 'expense';
  date: string;
  imported: boolean;
  created_at: string;
}

export interface Budget {
  id: number;
  category: string;
  monthly_limit: number;
  created_at: string;
}

export interface PlannedPurchase {
  id: number;
  name: string;
  estimated_cost: number;
  category: string;
  priority: 'high' | 'medium' | 'low';
  target_month: string;
  notes: string;
  status: 'planned' | 'purchased' | 'cancelled' | 'deferred';
  created_at: string;
}

export interface Investment {
  id: number;
  name: string;
  type: 'recurring' | 'one_time';
  amount: number;
  frequency: string;
  account_id: number | null;
  category: string;
  notes: string;
  status: 'planned' | 'active' | 'paused';
  created_at: string;
}

export interface BudgetStatus {
  category: string;
  limit: number;
  spent: number;
  remaining: number;
}

export interface PurchaseAffordability {
  id: number;
  name: string;
  cost: number;
  target_month: string;
  priority: string;
  affordable: boolean;
  suggested_month?: string;
  reason?: string;
}

export interface InvestmentSummary {
  name: string;
  amount: number;
  frequency: string;
  status: string;
  category: string;
}

export interface Dashboard {
  total_balance: number;
  accounts: Account[];
  monthly_income: number;
  monthly_expenses: number;
  monthly_investment_total: number;
  available_for_investments: number;
  budget_status: BudgetStatus[];
  planned_purchases: PurchaseAffordability[];
  investments: InvestmentSummary[];
}
