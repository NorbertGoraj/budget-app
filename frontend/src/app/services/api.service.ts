import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Account, Transaction, Budget, PlannedPurchase, Investment, Dashboard, Debt, DebtPayment } from './models';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private base = '/api';

  constructor(private http: HttpClient) {}

  // Accounts
  getAccounts(): Observable<Account[]> {
    return this.http.get<Account[]>(`${this.base}/accounts`);
  }
  createAccount(a: Partial<Account>): Observable<Account> {
    return this.http.post<Account>(`${this.base}/accounts`, a);
  }
  updateAccount(id: number, a: Partial<Account>): Observable<Account> {
    return this.http.put<Account>(`${this.base}/accounts/${id}`, a);
  }
  deleteAccount(id: number): Observable<any> {
    return this.http.delete(`${this.base}/accounts/${id}`);
  }

  // Transactions
  getTransactions(month?: string, accountId?: string, category?: string): Observable<Transaction[]> {
    let params = new HttpParams();
    if (month) params = params.set('month', month);
    if (accountId) params = params.set('account_id', accountId);
    if (category) params = params.set('category', category);
    return this.http.get<Transaction[]>(`${this.base}/transactions`, { params });
  }
  createTransaction(t: Partial<Transaction>): Observable<Transaction> {
    return this.http.post<Transaction>(`${this.base}/transactions`, t);
  }
  updateTransaction(id: number, t: Partial<Transaction>): Observable<Transaction> {
    return this.http.put<Transaction>(`${this.base}/transactions/${id}`, t);
  }
  deleteTransaction(id: number): Observable<any> {
    return this.http.delete(`${this.base}/transactions/${id}`);
  }

  // Import
  importCSV(accountId: number, file: File): Observable<any> {
    const formData = new FormData();
    formData.append('account_id', accountId.toString());
    formData.append('file', file);
    return this.http.post(`${this.base}/import/csv`, formData);
  }

  // Budgets
  getBudgets(): Observable<Budget[]> {
    return this.http.get<Budget[]>(`${this.base}/budgets`);
  }
  createBudget(b: Partial<Budget>): Observable<Budget> {
    return this.http.post<Budget>(`${this.base}/budgets`, b);
  }
  updateBudget(id: number, b: Partial<Budget>): Observable<Budget> {
    return this.http.put<Budget>(`${this.base}/budgets/${id}`, b);
  }
  deleteBudget(id: number): Observable<any> {
    return this.http.delete(`${this.base}/budgets/${id}`);
  }

  // Purchases
  getPurchases(): Observable<PlannedPurchase[]> {
    return this.http.get<PlannedPurchase[]>(`${this.base}/purchases`);
  }
  createPurchase(p: Partial<PlannedPurchase>): Observable<PlannedPurchase> {
    return this.http.post<PlannedPurchase>(`${this.base}/purchases`, p);
  }
  updatePurchase(id: number, p: Partial<PlannedPurchase>): Observable<PlannedPurchase> {
    return this.http.put<PlannedPurchase>(`${this.base}/purchases/${id}`, p);
  }
  deletePurchase(id: number): Observable<any> {
    return this.http.delete(`${this.base}/purchases/${id}`);
  }

  // Investments
  getInvestments(): Observable<Investment[]> {
    return this.http.get<Investment[]>(`${this.base}/investments`);
  }
  createInvestment(i: Partial<Investment>): Observable<Investment> {
    return this.http.post<Investment>(`${this.base}/investments`, i);
  }
  updateInvestment(id: number, i: Partial<Investment>): Observable<Investment> {
    return this.http.put<Investment>(`${this.base}/investments/${id}`, i);
  }
  deleteInvestment(id: number): Observable<any> {
    return this.http.delete(`${this.base}/investments/${id}`);
  }

  // Debts
  getDebts(): Observable<Debt[]> {
    return this.http.get<Debt[]>(`${this.base}/debts`);
  }
  createDebt(d: Partial<Debt>): Observable<Debt> {
    return this.http.post<Debt>(`${this.base}/debts`, d);
  }
  updateDebt(id: number, d: Partial<Debt>): Observable<Debt> {
    return this.http.put<Debt>(`${this.base}/debts/${id}`, d);
  }
  deleteDebt(id: number): Observable<any> {
    return this.http.delete(`${this.base}/debts/${id}`);
  }
  getDebtPayments(debtId: number): Observable<DebtPayment[]> {
    return this.http.get<DebtPayment[]>(`${this.base}/debts/${debtId}/payments`);
  }
  recordDebtPayment(debtId: number, p: Partial<DebtPayment>): Observable<DebtPayment> {
    return this.http.post<DebtPayment>(`${this.base}/debts/${debtId}/payments`, p);
  }
  deleteDebtPayment(paymentId: number): Observable<any> {
    return this.http.delete(`${this.base}/debts/payments/${paymentId}`);
  }

  // Dashboard
  getDashboard(): Observable<Dashboard> {
    return this.http.get<Dashboard>(`${this.base}/dashboard`);
  }
}
