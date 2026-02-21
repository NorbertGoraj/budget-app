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

export interface Debt {
  id: number;
  name: string;
  type: 'credit_card' | 'loan' | 'mortgage' | 'student_loan' | 'car_loan' | 'other';
  original_amount: number;
  current_balance: number;
  interest_rate: number;
  minimum_payment: number;
  due_day: number;
  status: 'active' | 'paid_off';
  notes: string;
  created_at: string;
}

export interface DebtPayment {
  id: number;
  debt_id: number;
  amount: number;
  paid_at: string;
  notes: string;
  created_at: string;
}

export interface DebtItem {
  id: number;
  name: string;
  type: string;
  current_balance: number;
  interest_rate: number;
  minimum_payment: number;
  due_day: number;
}

export interface DebtSummary {
  total_debt: number;
  monthly_minimum_payments: number;
  active_debts: DebtItem[];
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
  debt_summary: DebtSummary;
}
