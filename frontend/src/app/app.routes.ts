import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', loadComponent: () => import('./dashboard/dashboard.component').then(m => m.DashboardComponent) },
  { path: 'accounts', loadComponent: () => import('./accounts/accounts.component').then(m => m.AccountsComponent) },
  { path: 'transactions', loadComponent: () => import('./transactions/transactions.component').then(m => m.TransactionsComponent) },
  { path: 'import', loadComponent: () => import('./import/import.component').then(m => m.ImportComponent) },
  { path: 'budgets', loadComponent: () => import('./budgets/budgets.component').then(m => m.BudgetsComponent) },
  { path: 'purchases', loadComponent: () => import('./purchases/purchases.component').then(m => m.PurchasesComponent) },
  { path: 'investments', loadComponent: () => import('./investments/investments.component').then(m => m.InvestmentsComponent) },
  { path: 'debts', loadComponent: () => import('./debts/debts.component').then(m => m.DebtsComponent) },
  { path: '**', redirectTo: '' },
];
