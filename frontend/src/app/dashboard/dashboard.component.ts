import { Component, OnInit, inject } from '@angular/core';
import { CommonModule, CurrencyPipe, DecimalPipe } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatDividerModule } from '@angular/material/divider';
import { RouterLink } from '@angular/router';
import { ApiService } from '../services/api.service';
import { Dashboard } from '../services/models';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    CurrencyPipe,
    DecimalPipe,
    RouterLink,
    MatCardModule,
    MatProgressBarModule,
    MatIconModule,
    MatListModule,
    MatDividerModule,
  ],
  template: `
    <div class="dashboard-container" *ngIf="dashboard">
      <!-- Total Balance -->
      <mat-card class="balance-card">
        <mat-card-header>
          <mat-icon mat-card-avatar>account_balance_wallet</mat-icon>
          <mat-card-title>Total Balance</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <h1 class="balance-amount">{{ dashboard.total_balance | currency:'PLN':'symbol':'1.2-2' }}</h1>
        </mat-card-content>
      </mat-card>

      <!-- Account Balances -->
      <mat-card>
        <mat-card-header>
          <mat-icon mat-card-avatar>account_balance</mat-icon>
          <mat-card-title>Accounts</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <mat-list>
            @for (account of dashboard.accounts; track account.id) {
              <mat-list-item>
                <mat-icon matListItemIcon>
                  {{ account.type === 'bank' ? 'account_balance' : 'payments' }}
                </mat-icon>
                <span matListItemTitle>{{ account.name }}</span>
                <span matListItemLine>{{ account.type | titlecase }}</span>
                <span matListItemMeta>{{ account.balance | currency:'PLN':'symbol':'1.2-2' }}</span>
              </mat-list-item>
            }
          </mat-list>
        </mat-card-content>
      </mat-card>

      <!-- Monthly Income / Expense Summary -->
      <mat-card>
        <mat-card-header>
          <mat-icon mat-card-avatar>bar_chart</mat-icon>
          <mat-card-title>Monthly Summary</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="summary-row">
            <div class="summary-item income">
              <mat-icon>trending_up</mat-icon>
              <div>
                <span class="label">Income</span>
                <span class="value">{{ dashboard.monthly_income | currency:'PLN':'symbol':'1.2-2' }}</span>
              </div>
            </div>
            <div class="summary-item expense">
              <mat-icon>trending_down</mat-icon>
              <div>
                <span class="label">Expenses</span>
                <span class="value">{{ dashboard.monthly_expenses | currency:'PLN':'symbol':'1.2-2' }}</span>
              </div>
            </div>
            <div class="summary-item net">
              <mat-icon>swap_vert</mat-icon>
              <div>
                <span class="label">Net</span>
                <span class="value"
                  [class.positive]="dashboard.monthly_income - dashboard.monthly_expenses >= 0"
                  [class.negative]="dashboard.monthly_income - dashboard.monthly_expenses < 0">
                  {{ dashboard.monthly_income - dashboard.monthly_expenses | currency:'PLN':'symbol':'1.2-2' }}
                </span>
              </div>
            </div>
          </div>
        </mat-card-content>
      </mat-card>

      <!-- Budget Status -->
      <mat-card>
        <mat-card-header>
          <mat-icon mat-card-avatar>pie_chart</mat-icon>
          <mat-card-title>Budget Status</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          @for (budget of dashboard.budget_status; track budget.category) {
            <div class="budget-item">
              <div class="budget-header">
                <span class="budget-category">{{ budget.category }}</span>
                <span class="budget-values">
                  {{ budget.spent | currency:'PLN':'symbol':'1.2-2' }}
                  / {{ budget.limit | currency:'PLN':'symbol':'1.2-2' }}
                </span>
              </div>
              <mat-progress-bar
                [mode]="'determinate'"
                [value]="getBudgetPercent(budget)"
                [color]="getBudgetPercent(budget) > 90 ? 'warn' : 'primary'">
              </mat-progress-bar>
              <div class="budget-remaining"
                [class.over-budget]="budget.remaining < 0">
                {{ budget.remaining >= 0 ? 'Remaining: ' : 'Over by: ' }}
                {{ (budget.remaining >= 0 ? budget.remaining : -budget.remaining) | currency:'PLN':'symbol':'1.2-2' }}
              </div>
            </div>
          }
          @if (dashboard.budget_status.length === 0) {
            <p class="empty-message">No budgets configured yet.</p>
          }
        </mat-card-content>
      </mat-card>

      <!-- Investment Summary -->
      <mat-card>
        <mat-card-header>
          <mat-icon mat-card-avatar>show_chart</mat-icon>
          <mat-card-title>Investments</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="summary-row">
            <div class="summary-item">
              <mat-icon>folder_open</mat-icon>
              <div>
                <span class="label">Active</span>
                <span class="value">{{ getActiveInvestments() }}</span>
              </div>
            </div>
            <div class="summary-item">
              <mat-icon>calendar_month</mat-icon>
              <div>
                <span class="label">Monthly Total</span>
                <span class="value">{{ dashboard.monthly_investment_total | currency:'PLN':'symbol':'1.2-2' }}</span>
              </div>
            </div>
            <div class="summary-item">
              <mat-icon>savings</mat-icon>
              <div>
                <span class="label">Available</span>
                <span class="value">{{ dashboard.available_for_investments | currency:'PLN':'symbol':'1.2-2' }}</span>
              </div>
            </div>
          </div>
          <mat-divider></mat-divider>
          <mat-list>
            @for (inv of dashboard.investments; track inv.name) {
              <mat-list-item>
                <mat-icon matListItemIcon>
                  {{ inv.status === 'active' ? 'check_circle' : 'pause_circle' }}
                </mat-icon>
                <span matListItemTitle>{{ inv.name }}</span>
                <span matListItemLine>{{ inv.category }} &middot; {{ inv.frequency }}</span>
                <span matListItemMeta>{{ inv.amount | currency:'PLN':'symbol':'1.2-2' }}</span>
              </mat-list-item>
            }
          </mat-list>
        </mat-card-content>
      </mat-card>

      <!-- Planned Purchases -->
      <mat-card>
        <mat-card-header>
          <mat-icon mat-card-avatar>shopping_cart</mat-icon>
          <mat-card-title>Planned Purchases</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          @for (purchase of dashboard.planned_purchases; track purchase.id) {
            <div class="purchase-item" [class.affordable]="purchase.affordable" [class.not-affordable]="!purchase.affordable">
              <div class="purchase-header">
                <mat-icon [style.color]="purchase.affordable ? '#4caf50' : '#f44336'">
                  {{ purchase.affordable ? 'check_circle' : 'cancel' }}
                </mat-icon>
                <span class="purchase-name">{{ purchase.name }}</span>
                <span class="purchase-cost">{{ purchase.cost | currency:'PLN':'symbol':'1.2-2' }}</span>
              </div>
              <div class="purchase-details">
                <span>Target: {{ purchase.target_month }}</span>
                <span>Priority: {{ purchase.priority | titlecase }}</span>
              </div>
              @if (!purchase.affordable && purchase.suggested_month) {
                <div class="purchase-suggestion">
                  <mat-icon>info</mat-icon>
                  Suggested month: {{ purchase.suggested_month }}
                  @if (purchase.reason) {
                    &mdash; {{ purchase.reason }}
                  }
                </div>
              }
            </div>
          }
          @if (dashboard.planned_purchases.length === 0) {
            <p class="empty-message">No planned purchases.</p>
          }
        </mat-card-content>
      </mat-card>

      <!-- Debt Summary -->
      <mat-card>
        <mat-card-header>
          <mat-icon mat-card-avatar>credit_card</mat-icon>
          <mat-card-title>Debts</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="summary-row">
            <div class="summary-item">
              <mat-icon>account_balance</mat-icon>
              <div>
                <span class="label">Total Debt</span>
                <span class="value debt-total">
                  {{ dashboard.debt_summary.total_debt | currency:'PLN':'symbol':'1.2-2' }}
                </span>
              </div>
            </div>
            <div class="summary-item">
              <mat-icon>calendar_month</mat-icon>
              <div>
                <span class="label">Monthly Minimums</span>
                <span class="value">
                  {{ dashboard.debt_summary.monthly_minimum_payments | currency:'PLN':'symbol':'1.2-2' }}
                </span>
              </div>
            </div>
          </div>
          <mat-divider></mat-divider>
          <mat-list>
            @for (debt of dashboard.debt_summary.active_debts; track debt.id) {
              <mat-list-item>
                <mat-icon matListItemIcon>credit_card</mat-icon>
                <span matListItemTitle>{{ debt.name }}</span>
                <span matListItemLine>{{ debt.type }} &middot; {{ debt.interest_rate | number:'1.1-2' }}% APR &middot; due day {{ debt.due_day }}</span>
                <span matListItemMeta>{{ debt.current_balance | currency:'PLN':'symbol':'1.2-2' }}</span>
              </mat-list-item>
            }
          </mat-list>
          @if (dashboard.debt_summary.active_debts.length === 0) {
            <p class="empty-message">No active debts. <a routerLink="/debts">Add one</a></p>
          }
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .dashboard-container {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
      gap: 16px;
      padding: 16px;
    }
    .balance-card {
      grid-column: 1 / -1;
    }
    .balance-amount {
      font-size: 2.5rem;
      font-weight: 300;
      margin: 16px 0 0;
    }
    .summary-row {
      display: flex;
      gap: 24px;
      flex-wrap: wrap;
      margin: 16px 0;
    }
    .summary-item {
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .summary-item .label {
      display: block;
      font-size: 0.8rem;
      color: rgba(0,0,0,0.6);
    }
    .summary-item .value {
      display: block;
      font-size: 1.2rem;
      font-weight: 500;
    }
    .income .value { color: #4caf50; }
    .expense .value { color: #f44336; }
    .positive { color: #4caf50; }
    .negative { color: #f44336; }

    .budget-item {
      margin: 16px 0;
    }
    .budget-header {
      display: flex;
      justify-content: space-between;
      margin-bottom: 4px;
    }
    .budget-category {
      font-weight: 500;
    }
    .budget-values {
      font-size: 0.85rem;
      color: rgba(0,0,0,0.6);
    }
    .budget-remaining {
      font-size: 0.8rem;
      margin-top: 4px;
      color: rgba(0,0,0,0.6);
    }
    .over-budget {
      color: #f44336;
      font-weight: 500;
    }

    .purchase-item {
      padding: 12px;
      margin: 8px 0;
      border-radius: 8px;
      border-left: 4px solid;
    }
    .purchase-item.affordable {
      border-left-color: #4caf50;
      background: rgba(76, 175, 80, 0.05);
    }
    .purchase-item.not-affordable {
      border-left-color: #f44336;
      background: rgba(244, 67, 54, 0.05);
    }
    .purchase-header {
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .purchase-name {
      flex: 1;
      font-weight: 500;
    }
    .purchase-cost {
      font-weight: 500;
    }
    .purchase-details {
      display: flex;
      gap: 16px;
      margin-top: 4px;
      padding-left: 32px;
      font-size: 0.85rem;
      color: rgba(0,0,0,0.6);
    }
    .purchase-suggestion {
      display: flex;
      align-items: center;
      gap: 4px;
      margin-top: 4px;
      padding-left: 32px;
      font-size: 0.85rem;
      color: #ff9800;
    }
    .purchase-suggestion mat-icon {
      font-size: 18px;
      width: 18px;
      height: 18px;
      color: #ff9800;
    }
    .empty-message {
      color: rgba(0,0,0,0.5);
      font-style: italic;
      text-align: center;
      padding: 16px;
    }
    .debt-total { color: #c62828; }
  `],
})
export class DashboardComponent implements OnInit {
  private api = inject(ApiService);

  dashboard: Dashboard | null = null;

  ngOnInit(): void {
    this.loadDashboard();
  }

  loadDashboard(): void {
    this.api.getDashboard().subscribe((data) => {
      this.dashboard = data;
    });
  }

  getBudgetPercent(budget: { spent: number; limit: number }): number {
    if (budget.limit <= 0) return 0;
    return Math.min((budget.spent / budget.limit) * 100, 100);
  }

  getActiveInvestments(): number {
    return this.dashboard?.investments.filter((i) => i.status === 'active').length ?? 0;
  }
}
