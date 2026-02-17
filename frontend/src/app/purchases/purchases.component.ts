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
import { MatChipsModule } from '@angular/material/chips';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ApiService } from '../services/api.service';
import { PlannedPurchase, PurchaseAffordability } from '../services/models';

@Component({
  selector: 'app-purchases',
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
    MatChipsModule,
    MatTooltipModule,
    MatSnackBarModule,
  ],
  template: `
    <div class="purchases-container">
      <h2>Planned Purchases</h2>

      <!-- Add Purchase Form -->
      <mat-card class="form-card">
        <mat-card-header>
          <mat-card-title>Add Planned Purchase</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="form" (ngSubmit)="addPurchase()" class="purchase-form">
            <mat-form-field>
              <mat-label>Name</mat-label>
              <input matInput formControlName="name" />
            </mat-form-field>

            <mat-form-field>
              <mat-label>Estimated Cost</mat-label>
              <input matInput type="number" formControlName="estimated_cost" />
              <span matTextPrefix>$&nbsp;</span>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Category</mat-label>
              <input matInput formControlName="category" />
            </mat-form-field>

            <mat-form-field>
              <mat-label>Priority</mat-label>
              <mat-select formControlName="priority">
                <mat-option value="high">High</mat-option>
                <mat-option value="medium">Medium</mat-option>
                <mat-option value="low">Low</mat-option>
              </mat-select>
            </mat-form-field>

            <mat-form-field>
              <mat-label>Target Month (YYYY-MM)</mat-label>
              <input matInput formControlName="target_month" placeholder="2026-03" />
            </mat-form-field>

            <mat-form-field>
              <mat-label>Notes</mat-label>
              <textarea matInput formControlName="notes" rows="2"></textarea>
            </mat-form-field>

            <button mat-raised-button color="primary" type="submit" [disabled]="form.invalid">
              <mat-icon>add</mat-icon> Add Purchase
            </button>
          </form>
        </mat-card-content>
      </mat-card>

      <!-- Purchases Table -->
      <mat-card class="table-card">
        <mat-card-content>
          <table mat-table [dataSource]="purchases()" class="full-width">
            <ng-container matColumnDef="name">
              <th mat-header-cell *matHeaderCellDef>Name</th>
              <td mat-cell *matCellDef="let p">{{ p.name }}</td>
            </ng-container>

            <ng-container matColumnDef="estimated_cost">
              <th mat-header-cell *matHeaderCellDef>Estimated Cost</th>
              <td mat-cell *matCellDef="let p">{{ p.estimated_cost | currency }}</td>
            </ng-container>

            <ng-container matColumnDef="category">
              <th mat-header-cell *matHeaderCellDef>Category</th>
              <td mat-cell *matCellDef="let p">{{ p.category }}</td>
            </ng-container>

            <ng-container matColumnDef="priority">
              <th mat-header-cell *matHeaderCellDef>Priority</th>
              <td mat-cell *matCellDef="let p">
                <span class="priority-badge" [class]="'priority-' + p.priority">
                  {{ p.priority }}
                </span>
              </td>
            </ng-container>

            <ng-container matColumnDef="target_month">
              <th mat-header-cell *matHeaderCellDef>Target Month</th>
              <td mat-cell *matCellDef="let p">{{ p.target_month }}</td>
            </ng-container>

            <ng-container matColumnDef="status">
              <th mat-header-cell *matHeaderCellDef>Status</th>
              <td mat-cell *matCellDef="let p">
                <span class="status-chip" [class]="'status-' + p.status">
                  {{ p.status }}
                </span>
              </td>
            </ng-container>

            <ng-container matColumnDef="affordability">
              <th mat-header-cell *matHeaderCellDef>Affordability</th>
              <td mat-cell *matCellDef="let p">
                @if (getAffordability(p.id); as aff) {
                  @if (aff.affordable) {
                    <span class="afford-chip affordable">Affordable</span>
                  } @else {
                    <span class="afford-chip not-affordable"
                          [matTooltip]="aff.reason || ''">
                      Not Affordable
                    </span>
                    @if (aff.suggested_month) {
                      <span class="suggested-month">
                        Suggested: {{ aff.suggested_month }}
                      </span>
                    }
                  }
                }
              </td>
            </ng-container>

            <ng-container matColumnDef="actions">
              <th mat-header-cell *matHeaderCellDef>Actions</th>
              <td mat-cell *matCellDef="let p">
                @if (p.status === 'planned') {
                  <button mat-icon-button color="primary"
                          matTooltip="Mark as Purchased"
                          (click)="updateStatus(p, 'purchased')">
                    <mat-icon>shopping_cart</mat-icon>
                  </button>
                  <button mat-icon-button color="accent"
                          matTooltip="Defer"
                          (click)="updateStatus(p, 'deferred')">
                    <mat-icon>schedule</mat-icon>
                  </button>
                  <button mat-icon-button color="warn"
                          matTooltip="Cancel"
                          (click)="updateStatus(p, 'cancelled')">
                    <mat-icon>cancel</mat-icon>
                  </button>
                }
                @if (p.status === 'deferred') {
                  <button mat-icon-button color="primary"
                          matTooltip="Replan"
                          (click)="updateStatus(p, 'planned')">
                    <mat-icon>replay</mat-icon>
                  </button>
                }
                <button mat-icon-button color="warn"
                        matTooltip="Delete"
                        (click)="deletePurchase(p.id)">
                  <mat-icon>delete</mat-icon>
                </button>
              </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;"
                [class.not-affordable-row]="!isAffordable(row.id)"></tr>
          </table>

          @if (purchases().length === 0) {
            <p class="empty-message">No planned purchases yet. Add one above.</p>
          }
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .purchases-container {
      padding: 24px;
      max-width: 1200px;
      margin: 0 auto;
    }
    .form-card {
      margin-bottom: 24px;
    }
    .purchase-form {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      align-items: flex-start;
      padding-top: 16px;
    }
    .purchase-form mat-form-field {
      flex: 1 1 200px;
    }
    .purchase-form button {
      margin-top: 8px;
    }
    .table-card {
      overflow-x: auto;
    }
    .full-width {
      width: 100%;
    }
    .priority-badge {
      padding: 4px 10px;
      border-radius: 12px;
      font-size: 12px;
      font-weight: 500;
      text-transform: uppercase;
    }
    .priority-high { background: #ffcdd2; color: #b71c1c; }
    .priority-medium { background: #fff9c4; color: #f57f17; }
    .priority-low { background: #c8e6c9; color: #1b5e20; }
    .status-chip {
      padding: 4px 10px;
      border-radius: 12px;
      font-size: 12px;
      font-weight: 500;
      text-transform: uppercase;
    }
    .status-planned { background: #bbdefb; color: #0d47a1; }
    .status-purchased { background: #c8e6c9; color: #1b5e20; }
    .status-cancelled { background: #ffcdd2; color: #b71c1c; }
    .status-deferred { background: #fff9c4; color: #f57f17; }
    .afford-chip {
      padding: 4px 10px;
      border-radius: 12px;
      font-size: 12px;
      font-weight: 500;
    }
    .affordable { background: #c8e6c9; color: #1b5e20; }
    .not-affordable { background: #ffcdd2; color: #b71c1c; }
    .suggested-month {
      display: block;
      font-size: 11px;
      color: #666;
      margin-top: 4px;
    }
    .not-affordable-row {
      background: #fff3f0;
    }
    .empty-message {
      text-align: center;
      padding: 24px;
      color: #666;
    }
  `],
})
export class PurchasesComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private snackBar = inject(MatSnackBar);

  purchases = signal<PlannedPurchase[]>([]);
  affordabilityMap = signal<Map<number, PurchaseAffordability>>(new Map());

  displayedColumns = [
    'name', 'estimated_cost', 'category', 'priority',
    'target_month', 'status', 'affordability', 'actions',
  ];

  form = this.fb.group({
    name: ['', Validators.required],
    estimated_cost: [null as number | null, [Validators.required, Validators.min(0.01)]],
    category: ['', Validators.required],
    priority: ['medium' as 'high' | 'medium' | 'low', Validators.required],
    target_month: ['', [Validators.required, Validators.pattern(/^\d{4}-\d{2}$/)]],
    notes: [''],
  });

  ngOnInit(): void {
    this.loadPurchases();
    this.loadAffordability();
  }

  loadPurchases(): void {
    this.api.getPurchases().subscribe({
      next: (data) => this.purchases.set(data),
      error: () => this.snackBar.open('Failed to load purchases', 'Close', { duration: 3000 }),
    });
  }

  loadAffordability(): void {
    this.api.getDashboard().subscribe({
      next: (dashboard) => {
        const map = new Map<number, PurchaseAffordability>();
        for (const pa of dashboard.planned_purchases) {
          map.set(pa.id, pa);
        }
        this.affordabilityMap.set(map);
      },
    });
  }

  getAffordability(id: number): PurchaseAffordability | undefined {
    return this.affordabilityMap().get(id);
  }

  isAffordable(id: number): boolean {
    const aff = this.affordabilityMap().get(id);
    return !aff || aff.affordable;
  }

  addPurchase(): void {
    if (this.form.invalid) return;
    this.api.createPurchase(this.form.value as Partial<PlannedPurchase>).subscribe({
      next: () => {
        this.snackBar.open('Purchase added', 'Close', { duration: 2000 });
        this.form.reset({ priority: 'medium' });
        this.loadPurchases();
        this.loadAffordability();
      },
      error: () => this.snackBar.open('Failed to add purchase', 'Close', { duration: 3000 }),
    });
  }

  updateStatus(purchase: PlannedPurchase, status: PlannedPurchase['status']): void {
    this.api.updatePurchase(purchase.id, { status }).subscribe({
      next: () => {
        this.snackBar.open(`Status updated to ${status}`, 'Close', { duration: 2000 });
        this.loadPurchases();
        this.loadAffordability();
      },
      error: () => this.snackBar.open('Failed to update status', 'Close', { duration: 3000 }),
    });
  }

  deletePurchase(id: number): void {
    this.api.deletePurchase(id).subscribe({
      next: () => {
        this.snackBar.open('Purchase deleted', 'Close', { duration: 2000 });
        this.loadPurchases();
        this.loadAffordability();
      },
      error: () => this.snackBar.open('Failed to delete purchase', 'Close', { duration: 3000 }),
    });
  }
}
