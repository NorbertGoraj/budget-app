import { Component, OnInit, inject } from '@angular/core';
import { CommonModule, CurrencyPipe, DatePipe } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDividerModule } from '@angular/material/divider';
import { ApiService } from '../services/api.service';
import { Account, Transaction } from '../services/models';

@Component({
  selector: 'app-transactions',
  standalone: true,
  imports: [
    CommonModule,
    CurrencyPipe,
    DatePipe,
    ReactiveFormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatSnackBarModule,
    MatDividerModule,
  ],
  template: `
    <div class="transactions-container">
      <!-- Add Transaction Form -->
      <mat-card class="form-card">
        <mat-card-header>
          <mat-card-title>Add Transaction</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="form" (ngSubmit)="onSubmit()" class="transaction-form">
            <mat-form-field appearance="outline">
              <mat-label>Account</mat-label>
              <mat-select formControlName="account_id">
                @for (account of accounts; track account.id) {
                  <mat-option [value]="account.id">{{ account.name }}</mat-option>
                }
              </mat-select>
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Type</mat-label>
              <mat-select formControlName="type">
                <mat-option value="income">Income</mat-option>
                <mat-option value="expense">Expense</mat-option>
              </mat-select>
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Amount</mat-label>
              <input matInput type="number" formControlName="amount" placeholder="0.00" min="0.01" step="0.01">
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Description</mat-label>
              <input matInput formControlName="description" placeholder="Description">
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Category</mat-label>
              <input matInput formControlName="category" placeholder="Category">
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Date</mat-label>
              <input matInput type="date" formControlName="date">
            </mat-form-field>

            <div class="form-actions">
              <button mat-raised-button color="primary" type="submit" [disabled]="form.invalid">
                Add
              </button>
            </div>
          </form>
        </mat-card-content>
      </mat-card>

      <!-- Filters -->
      <mat-card>
        <mat-card-header>
          <mat-card-title>Filters</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="filterForm" class="filter-form">
            <mat-form-field appearance="outline">
              <mat-label>Month (YYYY-MM)</mat-label>
              <input matInput formControlName="month" placeholder="2026-02" maxlength="7">
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Account</mat-label>
              <mat-select formControlName="account_id">
                <mat-option [value]="''">All</mat-option>
                @for (account of accounts; track account.id) {
                  <mat-option [value]="account.id.toString()">{{ account.name }}</mat-option>
                }
              </mat-select>
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Category</mat-label>
              <input matInput formControlName="category" placeholder="Filter by category">
            </mat-form-field>

            <div class="form-actions">
              <button mat-raised-button color="accent" type="button" (click)="applyFilters()">
                <mat-icon>filter_list</mat-icon> Apply
              </button>
              <button mat-button type="button" (click)="clearFilters()">Clear</button>
            </div>
          </form>
        </mat-card-content>
      </mat-card>

      <!-- Transactions Table -->
      <mat-card>
        <mat-card-header>
          <mat-card-title>Transactions</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <table mat-table [dataSource]="transactions" class="full-width">
            <ng-container matColumnDef="date">
              <th mat-header-cell *matHeaderCellDef>Date</th>
              <td mat-cell *matCellDef="let t">{{ t.date | date:'yyyy-MM-dd' }}</td>
            </ng-container>

            <ng-container matColumnDef="description">
              <th mat-header-cell *matHeaderCellDef>Description</th>
              <td mat-cell *matCellDef="let t">{{ t.description }}</td>
            </ng-container>

            <ng-container matColumnDef="category">
              <th mat-header-cell *matHeaderCellDef>Category</th>
              <td mat-cell *matCellDef="let t">{{ t.category }}</td>
            </ng-container>

            <ng-container matColumnDef="amount">
              <th mat-header-cell *matHeaderCellDef>Amount</th>
              <td mat-cell *matCellDef="let t"
                [class.income-text]="t.type === 'income'"
                [class.expense-text]="t.type === 'expense'">
                {{ t.type === 'expense' ? '-' : '+' }}{{ t.amount | currency:'PLN':'symbol':'1.2-2' }}
              </td>
            </ng-container>

            <ng-container matColumnDef="type">
              <th mat-header-cell *matHeaderCellDef>Type</th>
              <td mat-cell *matCellDef="let t"
                [class.income-text]="t.type === 'income'"
                [class.expense-text]="t.type === 'expense'">
                {{ t.type | titlecase }}
              </td>
            </ng-container>

            <ng-container matColumnDef="account">
              <th mat-header-cell *matHeaderCellDef>Account</th>
              <td mat-cell *matCellDef="let t">{{ getAccountName(t.account_id) }}</td>
            </ng-container>

            <ng-container matColumnDef="actions">
              <th mat-header-cell *matHeaderCellDef>Actions</th>
              <td mat-cell *matCellDef="let t">
                <button mat-icon-button color="warn" (click)="deleteTransaction(t)">
                  <mat-icon>delete</mat-icon>
                </button>
              </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;"></tr>
          </table>

          @if (transactions.length === 0) {
            <p class="empty-message">No transactions found.</p>
          }
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .transactions-container {
      padding: 16px;
      display: flex;
      flex-direction: column;
      gap: 16px;
      max-width: 1100px;
      margin: 0 auto;
    }
    .transaction-form, .filter-form {
      display: flex;
      gap: 12px;
      align-items: flex-start;
      flex-wrap: wrap;
    }
    .transaction-form mat-form-field,
    .filter-form mat-form-field {
      flex: 1;
      min-width: 150px;
    }
    .form-actions {
      display: flex;
      gap: 8px;
      align-items: center;
      padding-top: 8px;
    }
    .full-width {
      width: 100%;
    }
    .income-text {
      color: #4caf50;
      font-weight: 500;
    }
    .expense-text {
      color: #f44336;
      font-weight: 500;
    }
    .empty-message {
      text-align: center;
      color: rgba(0,0,0,0.5);
      font-style: italic;
      padding: 24px;
    }
  `],
})
export class TransactionsComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private snackBar = inject(MatSnackBar);

  transactions: Transaction[] = [];
  accounts: Account[] = [];
  displayedColumns = ['date', 'description', 'category', 'amount', 'type', 'account', 'actions'];

  form: FormGroup = this.fb.group({
    account_id: [null, Validators.required],
    type: ['expense', Validators.required],
    amount: [null, [Validators.required, Validators.min(0.01)]],
    description: ['', Validators.required],
    category: ['', Validators.required],
    date: [this.todayString(), Validators.required],
  });

  filterForm: FormGroup = this.fb.group({
    month: [''],
    account_id: [''],
    category: [''],
  });

  ngOnInit(): void {
    this.api.getAccounts().subscribe((data) => {
      this.accounts = data;
    });
    this.loadTransactions();
  }

  loadTransactions(): void {
    const f = this.filterForm.value;
    this.api
      .getTransactions(f.month || undefined, f.account_id || undefined, f.category || undefined)
      .subscribe((data) => {
        this.transactions = data;
      });
  }

  onSubmit(): void {
    if (this.form.invalid) return;
    this.api.createTransaction(this.form.value).subscribe(() => {
      this.snackBar.open('Transaction added', 'Close', { duration: 2000 });
      this.form.reset({ account_id: null, type: 'expense', amount: null, description: '', category: '', date: this.todayString() });
      this.loadTransactions();
    });
  }

  applyFilters(): void {
    this.loadTransactions();
  }

  clearFilters(): void {
    this.filterForm.reset({ month: '', account_id: '', category: '' });
    this.loadTransactions();
  }

  deleteTransaction(t: Transaction): void {
    if (!confirm(`Delete transaction "${t.description}"?`)) return;
    this.api.deleteTransaction(t.id).subscribe(() => {
      this.snackBar.open('Transaction deleted', 'Close', { duration: 2000 });
      this.loadTransactions();
    });
  }

  getAccountName(accountId: number): string {
    const account = this.accounts.find((a) => a.id === accountId);
    return account ? account.name : `#${accountId}`;
  }

  private todayString(): string {
    return new Date().toISOString().substring(0, 10);
  }
}
