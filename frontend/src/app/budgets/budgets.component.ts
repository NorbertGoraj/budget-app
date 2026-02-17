import { Component, OnInit, inject } from '@angular/core';
import { CommonModule, CurrencyPipe } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ApiService } from '../services/api.service';
import { Budget, BudgetStatus } from '../services/models';

@Component({
  selector: 'app-budgets',
  standalone: true,
  imports: [
    CommonModule,
    CurrencyPipe,
    ReactiveFormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatProgressBarModule,
    MatSnackBarModule,
  ],
  template: `
    <div class="budgets-container">
      <!-- Add / Edit Form -->
      <mat-card class="form-card">
        <mat-card-header>
          <mat-card-title>{{ editingId ? 'Edit Budget' : 'Add Budget' }}</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="form" (ngSubmit)="onSubmit()" class="budget-form">
            <mat-form-field appearance="outline">
              <mat-label>Category</mat-label>
              <input matInput formControlName="category" placeholder="e.g. Groceries">
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Monthly Limit</mat-label>
              <input matInput type="number" formControlName="monthly_limit" placeholder="0.00" min="0.01" step="0.01">
            </mat-form-field>

            <div class="form-actions">
              <button mat-raised-button color="primary" type="submit" [disabled]="form.invalid">
                {{ editingId ? 'Update' : 'Add' }}
              </button>
              @if (editingId) {
                <button mat-button type="button" (click)="cancelEdit()">Cancel</button>
              }
            </div>
          </form>
        </mat-card-content>
      </mat-card>

      <!-- Budgets Table -->
      <mat-card>
        <mat-card-header>
          <mat-card-title>Budgets</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <table mat-table [dataSource]="budgets" class="full-width">
            <ng-container matColumnDef="category">
              <th mat-header-cell *matHeaderCellDef>Category</th>
              <td mat-cell *matCellDef="let b">{{ b.category }}</td>
            </ng-container>

            <ng-container matColumnDef="monthly_limit">
              <th mat-header-cell *matHeaderCellDef>Monthly Limit</th>
              <td mat-cell *matCellDef="let b">
                {{ b.monthly_limit | currency:'PLN':'symbol':'1.2-2' }}
              </td>
            </ng-container>

            <ng-container matColumnDef="spending">
              <th mat-header-cell *matHeaderCellDef>Current Month Spending</th>
              <td mat-cell *matCellDef="let b" class="spending-cell">
                <div class="spending-info">
                  <span class="spending-text"
                    [class.over-budget]="getSpent(b.category) > b.monthly_limit">
                    {{ getSpent(b.category) | currency:'PLN':'symbol':'1.2-2' }}
                    / {{ b.monthly_limit | currency:'PLN':'symbol':'1.2-2' }}
                  </span>
                  <mat-progress-bar
                    [mode]="'determinate'"
                    [value]="getSpentPercent(b.category, b.monthly_limit)"
                    [color]="getSpentPercent(b.category, b.monthly_limit) > 90 ? 'warn' : 'primary'">
                  </mat-progress-bar>
                </div>
              </td>
            </ng-container>

            <ng-container matColumnDef="actions">
              <th mat-header-cell *matHeaderCellDef>Actions</th>
              <td mat-cell *matCellDef="let b">
                <button mat-icon-button color="primary" (click)="startEdit(b)">
                  <mat-icon>edit</mat-icon>
                </button>
                <button mat-icon-button color="warn" (click)="deleteBudget(b)">
                  <mat-icon>delete</mat-icon>
                </button>
              </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;"></tr>
          </table>

          @if (budgets.length === 0) {
            <p class="empty-message">No budgets configured. Add one above to start tracking spending.</p>
          }
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .budgets-container {
      padding: 16px;
      display: flex;
      flex-direction: column;
      gap: 16px;
      max-width: 900px;
      margin: 0 auto;
    }
    .budget-form {
      display: flex;
      gap: 12px;
      align-items: flex-start;
      flex-wrap: wrap;
    }
    .budget-form mat-form-field {
      flex: 1;
      min-width: 200px;
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
    .spending-cell {
      min-width: 220px;
    }
    .spending-info {
      display: flex;
      flex-direction: column;
      gap: 4px;
      padding: 4px 0;
    }
    .spending-text {
      font-size: 0.85rem;
      color: rgba(0,0,0,0.7);
    }
    .over-budget {
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
export class BudgetsComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private snackBar = inject(MatSnackBar);

  budgets: Budget[] = [];
  budgetStatus: BudgetStatus[] = [];
  editingId: number | null = null;
  displayedColumns = ['category', 'monthly_limit', 'spending', 'actions'];

  form: FormGroup = this.fb.group({
    category: ['', Validators.required],
    monthly_limit: [null, [Validators.required, Validators.min(0.01)]],
  });

  ngOnInit(): void {
    this.loadData();
  }

  loadData(): void {
    this.api.getBudgets().subscribe((data) => {
      this.budgets = data;
    });
    this.api.getDashboard().subscribe((data) => {
      this.budgetStatus = data.budget_status;
    });
  }

  onSubmit(): void {
    if (this.form.invalid) return;
    const value = this.form.value;

    if (this.editingId) {
      this.api.updateBudget(this.editingId, value).subscribe(() => {
        this.snackBar.open('Budget updated', 'Close', { duration: 2000 });
        this.cancelEdit();
        this.loadData();
      });
    } else {
      this.api.createBudget(value).subscribe(() => {
        this.snackBar.open('Budget created', 'Close', { duration: 2000 });
        this.form.reset({ category: '', monthly_limit: null });
        this.loadData();
      });
    }
  }

  startEdit(budget: Budget): void {
    this.editingId = budget.id;
    this.form.patchValue({
      category: budget.category,
      monthly_limit: budget.monthly_limit,
    });
  }

  cancelEdit(): void {
    this.editingId = null;
    this.form.reset({ category: '', monthly_limit: null });
  }

  deleteBudget(budget: Budget): void {
    if (!confirm(`Delete budget for "${budget.category}"?`)) return;
    this.api.deleteBudget(budget.id).subscribe(() => {
      this.snackBar.open('Budget deleted', 'Close', { duration: 2000 });
      this.loadData();
    });
  }

  getSpent(category: string): number {
    const status = this.budgetStatus.find((s) => s.category === category);
    return status ? status.spent : 0;
  }

  getSpentPercent(category: string, limit: number): number {
    if (limit <= 0) return 0;
    const spent = this.getSpent(category);
    return Math.min((spent / limit) * 100, 100);
  }
}
