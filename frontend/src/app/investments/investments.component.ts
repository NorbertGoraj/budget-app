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
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ApiService } from '../services/api.service';
import { Account, Investment } from '../services/models';

@Component({
  selector: 'app-investments',
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
    MatButtonToggleModule,
    MatTooltipModule,
    MatSnackBarModule,
  ],
  template: `
    <div class="investments-container">
      <h2>Investments</h2>

      <!-- Summary Card -->
      <div class="summary-row">
        <mat-card class="summary-card">
          <mat-card-header>
            <mat-card-title>Monthly Investment Commitment</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <span class="summary-value">{{ monthlyTotal() | currency }}</span>
          </mat-card-content>
        </mat-card>
        <mat-card class="summary-card">
          <mat-card-header>
            <mat-card-title>Available for Investments</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <span class="summary-value" [class.negative]="availableForInvestments() < 0">
              {{ availableForInvestments() | currency }}
            </span>
          </mat-card-content>
        </mat-card>
      </div>

      <!-- Add Investment Form -->
      <mat-card class="form-card">
        <mat-card-header>
          <mat-card-title>Add Investment</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="form" (ngSubmit)="addInvestment()" class="investment-form">
            <mat-form-field>
              <mat-label>Name</mat-label>
              <input matInput formControlName="name" />
            </mat-form-field>

            <mat-form-field>
              <mat-label>Type</mat-label>
              <mat-select formControlName="type">
                <mat-option value="recurring">Recurring</mat-option>
                <mat-option value="one_time">One-Time</mat-option>
              </mat-select>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Amount</mat-label>
              <input matInput type="number" formControlName="amount" />
              <span matTextPrefix>$&nbsp;</span>
            </mat-form-field>

            @if (form.get('type')?.value === 'recurring') {
              <mat-form-field>
                <mat-label>Frequency</mat-label>
                <mat-select formControlName="frequency">
                  <mat-option value="weekly">Weekly</mat-option>
                  <mat-option value="monthly">Monthly</mat-option>
                  <mat-option value="quarterly">Quarterly</mat-option>
                  <mat-option value="yearly">Yearly</mat-option>
                </mat-select>
              </mat-form-field>
            }

            <mat-form-field>
              <mat-label>Account</mat-label>
              <mat-select formControlName="account_id">
                <mat-option [value]="null">None</mat-option>
                @for (account of accounts(); track account.id) {
                  <mat-option [value]="account.id">{{ account.name }}</mat-option>
                }
              </mat-select>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Category</mat-label>
              <mat-select formControlName="category">
                <mat-option value="stocks">Stocks</mat-option>
                <mat-option value="ETF">ETF</mat-option>
                <mat-option value="crypto">Crypto</mat-option>
                <mat-option value="savings">Savings</mat-option>
                <mat-option value="other">Other</mat-option>
              </mat-select>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Notes</mat-label>
              <textarea matInput formControlName="notes" rows="2"></textarea>
            </mat-form-field>

            <button mat-raised-button color="primary" type="submit" [disabled]="form.invalid">
              <mat-icon>add</mat-icon> Add Investment
            </button>
          </form>
        </mat-card-content>
      </mat-card>

      <!-- Investments Table -->
      <mat-card class="table-card">
        <mat-card-content>
          <table mat-table [dataSource]="investments()" class="full-width">
            <ng-container matColumnDef="name">
              <th mat-header-cell *matHeaderCellDef>Name</th>
              <td mat-cell *matCellDef="let i">{{ i.name }}</td>
            </ng-container>

            <ng-container matColumnDef="type">
              <th mat-header-cell *matHeaderCellDef>Type</th>
              <td mat-cell *matCellDef="let i">
                <span class="type-badge" [class]="'type-' + i.type">
                  {{ i.type === 'one_time' ? 'One-Time' : 'Recurring' }}
                </span>
              </td>
            </ng-container>

            <ng-container matColumnDef="amount">
              <th mat-header-cell *matHeaderCellDef>Amount</th>
              <td mat-cell *matCellDef="let i">{{ i.amount | currency }}</td>
            </ng-container>

            <ng-container matColumnDef="frequency">
              <th mat-header-cell *matHeaderCellDef>Frequency</th>
              <td mat-cell *matCellDef="let i">
                {{ i.type === 'recurring' ? i.frequency : '-' }}
              </td>
            </ng-container>

            <ng-container matColumnDef="category">
              <th mat-header-cell *matHeaderCellDef>Category</th>
              <td mat-cell *matCellDef="let i">{{ i.category }}</td>
            </ng-container>

            <ng-container matColumnDef="status">
              <th mat-header-cell *matHeaderCellDef>Status</th>
              <td mat-cell *matCellDef="let i">
                <mat-button-toggle-group [value]="i.status"
                                         (change)="updateStatus(i, $event.value)"
                                         appearance="standard">
                  <mat-button-toggle value="planned"
                                     matTooltip="Planned">
                    <mat-icon>event_note</mat-icon>
                  </mat-button-toggle>
                  <mat-button-toggle value="active"
                                     matTooltip="Active">
                    <mat-icon>play_arrow</mat-icon>
                  </mat-button-toggle>
                  <mat-button-toggle value="paused"
                                     matTooltip="Paused">
                    <mat-icon>pause</mat-icon>
                  </mat-button-toggle>
                </mat-button-toggle-group>
              </td>
            </ng-container>

            <ng-container matColumnDef="actions">
              <th mat-header-cell *matHeaderCellDef>Actions</th>
              <td mat-cell *matCellDef="let i">
                <button mat-icon-button color="warn"
                        matTooltip="Delete"
                        (click)="deleteInvestment(i.id)">
                  <mat-icon>delete</mat-icon>
                </button>
              </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;"></tr>
          </table>

          @if (investments().length === 0) {
            <p class="empty-message">No investments yet. Add one above.</p>
          }
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .investments-container {
      padding: 24px;
      max-width: 1200px;
      margin: 0 auto;
    }
    .summary-row {
      display: flex;
      gap: 16px;
      margin-bottom: 24px;
    }
    .summary-card {
      flex: 1;
    }
    .summary-value {
      font-size: 28px;
      font-weight: 600;
      color: #1565c0;
    }
    .summary-value.negative {
      color: #c62828;
    }
    .form-card {
      margin-bottom: 24px;
    }
    .investment-form {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      align-items: flex-start;
      padding-top: 16px;
    }
    .investment-form mat-form-field {
      flex: 1 1 200px;
    }
    .investment-form button {
      margin-top: 8px;
    }
    .table-card {
      overflow-x: auto;
    }
    .full-width {
      width: 100%;
    }
    .type-badge {
      padding: 4px 10px;
      border-radius: 12px;
      font-size: 12px;
      font-weight: 500;
      text-transform: uppercase;
    }
    .type-recurring { background: #e3f2fd; color: #0d47a1; }
    .type-one_time { background: #f3e5f5; color: #6a1b9a; }
    .empty-message {
      text-align: center;
      padding: 24px;
      color: #666;
    }
    mat-button-toggle-group {
      height: 36px;
    }
  `],
})
export class InvestmentsComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private snackBar = inject(MatSnackBar);

  investments = signal<Investment[]>([]);
  accounts = signal<Account[]>([]);
  monthlyTotal = signal(0);
  availableForInvestments = signal(0);

  displayedColumns = ['name', 'type', 'amount', 'frequency', 'category', 'status', 'actions'];

  form = this.fb.group({
    name: ['', Validators.required],
    type: ['recurring' as 'recurring' | 'one_time', Validators.required],
    amount: [null as number | null, [Validators.required, Validators.min(0.01)]],
    frequency: ['monthly'],
    account_id: [null as number | null],
    category: ['stocks', Validators.required],
    notes: [''],
  });

  ngOnInit(): void {
    this.loadInvestments();
    this.loadAccounts();
    this.loadDashboard();
  }

  loadInvestments(): void {
    this.api.getInvestments().subscribe({
      next: (data) => this.investments.set(data),
      error: () => this.snackBar.open('Failed to load investments', 'Close', { duration: 3000 }),
    });
  }

  loadAccounts(): void {
    this.api.getAccounts().subscribe({
      next: (data) => this.accounts.set(data),
    });
  }

  loadDashboard(): void {
    this.api.getDashboard().subscribe({
      next: (dashboard) => {
        this.monthlyTotal.set(dashboard.monthly_investment_total);
        this.availableForInvestments.set(dashboard.available_for_investments);
      },
    });
  }

  addInvestment(): void {
    if (this.form.invalid) return;
    const value = { ...this.form.value } as Partial<Investment>;
    if (value.type === 'one_time') {
      value.frequency = '';
    }
    this.api.createInvestment(value).subscribe({
      next: () => {
        this.snackBar.open('Investment added', 'Close', { duration: 2000 });
        this.form.reset({ type: 'recurring', frequency: 'monthly', category: 'stocks' });
        this.loadInvestments();
        this.loadDashboard();
      },
      error: () => this.snackBar.open('Failed to add investment', 'Close', { duration: 3000 }),
    });
  }

  updateStatus(investment: Investment, status: string): void {
    this.api.updateInvestment(investment.id, { status: status as Investment['status'] }).subscribe({
      next: () => {
        this.snackBar.open(`Status updated to ${status}`, 'Close', { duration: 2000 });
        this.loadInvestments();
        this.loadDashboard();
      },
      error: () => this.snackBar.open('Failed to update status', 'Close', { duration: 3000 }),
    });
  }

  deleteInvestment(id: number): void {
    this.api.deleteInvestment(id).subscribe({
      next: () => {
        this.snackBar.open('Investment deleted', 'Close', { duration: 2000 });
        this.loadInvestments();
        this.loadDashboard();
      },
      error: () => this.snackBar.open('Failed to delete investment', 'Close', { duration: 3000 }),
    });
  }
}
