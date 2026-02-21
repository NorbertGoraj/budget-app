import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule, CurrencyPipe } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatCardModule } from '@angular/material/card';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ApiService } from '../services/api.service';
import { Debt, DebtPayment } from '../services/models';

@Component({
  selector: 'app-debts',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatCardModule,
    MatTooltipModule,
    MatChipsModule,
    MatDividerModule,
    MatSnackBarModule,
  ],
  template: `
    <div class="debts-container">
      <h2>Debts</h2>

      <!-- Summary Cards -->
      <div class="summary-row">
        <mat-card class="summary-card">
          <mat-card-header>
            <mat-card-title>Total Debt</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <span class="summary-value debt-value">{{ totalBalance() | currency }}</span>
          </mat-card-content>
        </mat-card>
        <mat-card class="summary-card">
          <mat-card-header>
            <mat-card-title>Monthly Minimums</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <span class="summary-value">{{ totalMonthlyMinimum() | currency }}</span>
          </mat-card-content>
        </mat-card>
        <mat-card class="summary-card">
          <mat-card-header>
            <mat-card-title>Active Debts</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <span class="summary-value">{{ activeCount() }}</span>
          </mat-card-content>
        </mat-card>
      </div>

      <!-- Add Debt Form -->
      <mat-card class="form-card">
        <mat-card-header>
          <mat-card-title>Add Debt</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="debtForm" (ngSubmit)="addDebt()" class="debt-form">
            <mat-form-field>
              <mat-label>Name</mat-label>
              <input matInput formControlName="name" placeholder="e.g. Visa Credit Card" />
            </mat-form-field>

            <mat-form-field>
              <mat-label>Type</mat-label>
              <mat-select formControlName="type">
                <mat-option value="credit_card">Credit Card</mat-option>
                <mat-option value="loan">Loan</mat-option>
                <mat-option value="mortgage">Mortgage</mat-option>
                <mat-option value="student_loan">Student Loan</mat-option>
                <mat-option value="car_loan">Car Loan</mat-option>
                <mat-option value="other">Other</mat-option>
              </mat-select>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Original Amount</mat-label>
              <input matInput type="number" formControlName="original_amount" />
              <span matTextPrefix>$&nbsp;</span>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Current Balance</mat-label>
              <input matInput type="number" formControlName="current_balance" />
              <span matTextPrefix>$&nbsp;</span>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Interest Rate (APR %)</mat-label>
              <input matInput type="number" formControlName="interest_rate" />
              <span matTextSuffix>%</span>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Minimum Payment</mat-label>
              <input matInput type="number" formControlName="minimum_payment" />
              <span matTextPrefix>$&nbsp;</span>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Due Day (1–28)</mat-label>
              <input matInput type="number" formControlName="due_day" />
            </mat-form-field>

            <mat-form-field>
              <mat-label>Notes</mat-label>
              <textarea matInput formControlName="notes" rows="2"></textarea>
            </mat-form-field>

            <button mat-raised-button color="primary" type="submit" [disabled]="debtForm.invalid">
              <mat-icon>add</mat-icon> Add Debt
            </button>
          </form>
        </mat-card-content>
      </mat-card>

      <!-- Debts Table -->
      <mat-card class="table-card">
        <mat-card-content>
          <table mat-table [dataSource]="debts()" class="full-width">

            <ng-container matColumnDef="name">
              <th mat-header-cell *matHeaderCellDef>Name</th>
              <td mat-cell *matCellDef="let d">
                <span class="debt-name">{{ d.name }}</span>
                <span class="type-badge" [class]="'type-' + d.type">{{ typeLabel(d.type) }}</span>
              </td>
            </ng-container>

            <ng-container matColumnDef="current_balance">
              <th mat-header-cell *matHeaderCellDef>Balance</th>
              <td mat-cell *matCellDef="let d">
                <span [class.paid-off-value]="d.status === 'paid_off'">
                  {{ d.current_balance | currency }}
                </span>
              </td>
            </ng-container>

            <ng-container matColumnDef="interest_rate">
              <th mat-header-cell *matHeaderCellDef>APR</th>
              <td mat-cell *matCellDef="let d">{{ d.interest_rate | number:'1.2-2' }}%</td>
            </ng-container>

            <ng-container matColumnDef="minimum_payment">
              <th mat-header-cell *matHeaderCellDef>Min. Payment</th>
              <td mat-cell *matCellDef="let d">{{ d.minimum_payment | currency }}</td>
            </ng-container>

            <ng-container matColumnDef="due_day">
              <th mat-header-cell *matHeaderCellDef>Due Day</th>
              <td mat-cell *matCellDef="let d">{{ d.due_day }}</td>
            </ng-container>

            <ng-container matColumnDef="status">
              <th mat-header-cell *matHeaderCellDef>Status</th>
              <td mat-cell *matCellDef="let d">
                <span class="status-badge" [class]="'status-' + d.status">
                  {{ d.status === 'paid_off' ? 'Paid Off' : 'Active' }}
                </span>
              </td>
            </ng-container>

            <ng-container matColumnDef="actions">
              <th mat-header-cell *matHeaderCellDef>Actions</th>
              <td mat-cell *matCellDef="let d">
                <button mat-icon-button color="primary"
                        matTooltip="Payments"
                        (click)="togglePayments(d)">
                  <mat-icon>payments</mat-icon>
                </button>
                <button mat-icon-button color="warn"
                        matTooltip="Delete"
                        (click)="deleteDebt(d.id)">
                  <mat-icon>delete</mat-icon>
                </button>
              </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;"
                [class.selected-row]="selectedDebt()?.id === row.id"></tr>
          </table>

          @if (debts().length === 0) {
            <p class="empty-message">No debts yet. Add one above.</p>
          }
        </mat-card-content>
      </mat-card>

      <!-- Payments Panel -->
      @if (selectedDebt()) {
        <mat-card class="payments-card">
          <mat-card-header>
            <mat-icon mat-card-avatar>receipt_long</mat-icon>
            <mat-card-title>Payments — {{ selectedDebt()!.name }}</mat-card-title>
            <mat-card-subtitle>
              Balance: {{ selectedDebt()!.current_balance | currency }}
            </mat-card-subtitle>
          </mat-card-header>
          <mat-card-content>

            <!-- Record Payment Form -->
            <form [formGroup]="paymentForm" (ngSubmit)="recordPayment()" class="payment-form">
              <mat-form-field>
                <mat-label>Amount</mat-label>
                <input matInput type="number" formControlName="amount" />
                <span matTextPrefix>$&nbsp;</span>
              </mat-form-field>

              <mat-form-field>
                <mat-label>Date</mat-label>
                <input matInput type="date" formControlName="paid_at" />
              </mat-form-field>

              <mat-form-field>
                <mat-label>Notes</mat-label>
                <input matInput formControlName="notes" />
              </mat-form-field>

              <button mat-raised-button color="accent" type="submit" [disabled]="paymentForm.invalid">
                <mat-icon>add</mat-icon> Record Payment
              </button>
            </form>

            <mat-divider></mat-divider>

            <!-- Payment History -->
            @if (payments().length > 0) {
              <table mat-table [dataSource]="payments()" class="full-width payments-table">

                <ng-container matColumnDef="paid_at">
                  <th mat-header-cell *matHeaderCellDef>Date</th>
                  <td mat-cell *matCellDef="let p">{{ p.paid_at }}</td>
                </ng-container>

                <ng-container matColumnDef="amount">
                  <th mat-header-cell *matHeaderCellDef>Amount</th>
                  <td mat-cell *matCellDef="let p">{{ p.amount | currency }}</td>
                </ng-container>

                <ng-container matColumnDef="notes">
                  <th mat-header-cell *matHeaderCellDef>Notes</th>
                  <td mat-cell *matCellDef="let p">{{ p.notes }}</td>
                </ng-container>

                <ng-container matColumnDef="actions">
                  <th mat-header-cell *matHeaderCellDef></th>
                  <td mat-cell *matCellDef="let p">
                    <button mat-icon-button color="warn"
                            matTooltip="Delete payment"
                            (click)="deletePayment(p.id)">
                      <mat-icon>delete</mat-icon>
                    </button>
                  </td>
                </ng-container>

                <tr mat-header-row *matHeaderRowDef="paymentColumns"></tr>
                <tr mat-row *matRowDef="let row; columns: paymentColumns;"></tr>
              </table>
            } @else {
              <p class="empty-message">No payments recorded yet.</p>
            }
          </mat-card-content>
        </mat-card>
      }
    </div>
  `,
  styles: [`
    .debts-container {
      padding: 24px;
      max-width: 1200px;
      margin: 0 auto;
    }
    .summary-row {
      display: flex;
      gap: 16px;
      margin-bottom: 24px;
      flex-wrap: wrap;
    }
    .summary-card {
      flex: 1;
      min-width: 160px;
    }
    .summary-value {
      font-size: 28px;
      font-weight: 600;
      color: #1565c0;
    }
    .summary-value.debt-value {
      color: #c62828;
    }
    .form-card, .table-card, .payments-card {
      margin-bottom: 24px;
    }
    .debt-form {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      align-items: flex-start;
      padding-top: 16px;
    }
    .debt-form mat-form-field {
      flex: 1 1 200px;
    }
    .debt-form button {
      margin-top: 8px;
    }
    .full-width {
      width: 100%;
    }
    .debt-name {
      font-weight: 500;
      margin-right: 8px;
    }
    .type-badge {
      padding: 2px 8px;
      border-radius: 10px;
      font-size: 11px;
      font-weight: 500;
      text-transform: uppercase;
      vertical-align: middle;
    }
    .type-credit_card  { background: #fce4ec; color: #880e4f; }
    .type-loan         { background: #e3f2fd; color: #0d47a1; }
    .type-mortgage     { background: #e8f5e9; color: #1b5e20; }
    .type-student_loan { background: #fff3e0; color: #e65100; }
    .type-car_loan     { background: #f3e5f5; color: #4a148c; }
    .type-other        { background: #f5f5f5; color: #424242; }
    .status-badge {
      padding: 4px 10px;
      border-radius: 12px;
      font-size: 12px;
      font-weight: 500;
    }
    .status-active   { background: #ffebee; color: #c62828; }
    .status-paid_off { background: #e8f5e9; color: #2e7d32; }
    .paid-off-value {
      color: #2e7d32;
    }
    .selected-row {
      background: rgba(21, 101, 192, 0.06);
    }
    .table-card {
      overflow-x: auto;
    }
    .payments-card {
      border-left: 4px solid #1565c0;
    }
    .payment-form {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      align-items: flex-start;
      padding: 16px 0;
    }
    .payment-form mat-form-field {
      flex: 1 1 160px;
    }
    .payment-form button {
      margin-top: 8px;
    }
    .payments-table {
      margin-top: 16px;
    }
    .empty-message {
      text-align: center;
      padding: 24px;
      color: #666;
    }
  `],
})
export class DebtsComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private snackBar = inject(MatSnackBar);

  debts = signal<Debt[]>([]);
  payments = signal<DebtPayment[]>([]);
  selectedDebt = signal<Debt | null>(null);

  displayedColumns = ['name', 'current_balance', 'interest_rate', 'minimum_payment', 'due_day', 'status', 'actions'];
  paymentColumns = ['paid_at', 'amount', 'notes', 'actions'];

  debtForm = this.fb.group({
    name:            ['', Validators.required],
    type:            ['credit_card' as Debt['type'], Validators.required],
    original_amount: [null as number | null, [Validators.required, Validators.min(0.01)]],
    current_balance: [null as number | null, [Validators.required, Validators.min(0)]],
    interest_rate:   [null as number | null, [Validators.required, Validators.min(0)]],
    minimum_payment: [null as number | null, [Validators.required, Validators.min(0)]],
    due_day:         [1, [Validators.required, Validators.min(1), Validators.max(28)]],
    notes:           [''],
  });

  paymentForm = this.fb.group({
    amount:  [null as number | null, [Validators.required, Validators.min(0.01)]],
    paid_at: [new Date().toISOString().slice(0, 10), Validators.required],
    notes:   [''],
  });

  ngOnInit(): void {
    this.loadDebts();
  }

  loadDebts(): void {
    this.api.getDebts().subscribe({
      next: (data) => {
        this.debts.set(data);
        // Refresh selected debt reference so balance shown in panel stays current
        const sel = this.selectedDebt();
        if (sel) {
          const refreshed = data.find(d => d.id === sel.id) ?? null;
          this.selectedDebt.set(refreshed);
        }
      },
      error: () => this.snackBar.open('Failed to load debts', 'Close', { duration: 3000 }),
    });
  }

  addDebt(): void {
    if (this.debtForm.invalid) return;
    this.api.createDebt(this.debtForm.value as Partial<Debt>).subscribe({
      next: () => {
        this.snackBar.open('Debt added', 'Close', { duration: 2000 });
        this.debtForm.reset({ type: 'credit_card', due_day: 1 });
        this.loadDebts();
      },
      error: () => this.snackBar.open('Failed to add debt', 'Close', { duration: 3000 }),
    });
  }

  deleteDebt(id: number): void {
    this.api.deleteDebt(id).subscribe({
      next: () => {
        this.snackBar.open('Debt deleted', 'Close', { duration: 2000 });
        if (this.selectedDebt()?.id === id) {
          this.selectedDebt.set(null);
          this.payments.set([]);
        }
        this.loadDebts();
      },
      error: () => this.snackBar.open('Failed to delete debt', 'Close', { duration: 3000 }),
    });
  }

  togglePayments(debt: Debt): void {
    if (this.selectedDebt()?.id === debt.id) {
      this.selectedDebt.set(null);
      this.payments.set([]);
    } else {
      this.selectedDebt.set(debt);
      this.loadPayments(debt.id);
    }
  }

  loadPayments(debtId: number): void {
    this.api.getDebtPayments(debtId).subscribe({
      next: (data) => this.payments.set(data),
      error: () => this.snackBar.open('Failed to load payments', 'Close', { duration: 3000 }),
    });
  }

  recordPayment(): void {
    const debt = this.selectedDebt();
    if (!debt || this.paymentForm.invalid) return;
    this.api.recordDebtPayment(debt.id, this.paymentForm.value as Partial<DebtPayment>).subscribe({
      next: () => {
        this.snackBar.open('Payment recorded', 'Close', { duration: 2000 });
        this.paymentForm.reset({ paid_at: new Date().toISOString().slice(0, 10) });
        this.loadPayments(debt.id);
        this.loadDebts();
      },
      error: () => this.snackBar.open('Failed to record payment', 'Close', { duration: 3000 }),
    });
  }

  deletePayment(paymentId: number): void {
    const debt = this.selectedDebt();
    if (!debt) return;
    this.api.deleteDebtPayment(paymentId).subscribe({
      next: () => {
        this.snackBar.open('Payment deleted', 'Close', { duration: 2000 });
        this.loadPayments(debt.id);
        this.loadDebts();
      },
      error: () => this.snackBar.open('Failed to delete payment', 'Close', { duration: 3000 }),
    });
  }

  // Computed summary helpers
  totalBalance(): number {
    return this.debts()
      .filter(d => d.status === 'active')
      .reduce((sum, d) => sum + d.current_balance, 0);
  }

  totalMonthlyMinimum(): number {
    return this.debts()
      .filter(d => d.status === 'active')
      .reduce((sum, d) => sum + d.minimum_payment, 0);
  }

  activeCount(): number {
    return this.debts().filter(d => d.status === 'active').length;
  }

  typeLabel(type: string): string {
    const labels: Record<string, string> = {
      credit_card:  'Credit Card',
      loan:         'Loan',
      mortgage:     'Mortgage',
      student_loan: 'Student Loan',
      car_loan:     'Car Loan',
      other:        'Other',
    };
    return labels[type] ?? type;
  }
}
