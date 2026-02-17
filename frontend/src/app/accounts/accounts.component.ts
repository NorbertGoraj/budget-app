import { Component, OnInit, inject } from '@angular/core';
import { CommonModule, CurrencyPipe } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ApiService } from '../services/api.service';
import { Account } from '../services/models';

@Component({
  selector: 'app-accounts',
  standalone: true,
  imports: [
    CommonModule,
    CurrencyPipe,
    ReactiveFormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatSnackBarModule,
  ],
  template: `
    <div class="accounts-container">
      <!-- Add / Edit Form -->
      <mat-card class="form-card">
        <mat-card-header>
          <mat-card-title>{{ editingId ? 'Edit Account' : 'Add Account' }}</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="form" (ngSubmit)="onSubmit()" class="account-form">
            <mat-form-field appearance="outline">
              <mat-label>Name</mat-label>
              <input matInput formControlName="name" placeholder="Account name">
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>Type</mat-label>
              <mat-select formControlName="type">
                <mat-option value="bank">Bank</mat-option>
                <mat-option value="cash">Cash</mat-option>
              </mat-select>
            </mat-form-field>

            <mat-form-field appearance="outline">
              <mat-label>{{ editingId ? 'Balance' : 'Initial Balance' }}</mat-label>
              <input matInput type="number" formControlName="balance" placeholder="0.00">
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

      <!-- Accounts Table -->
      <mat-card>
        <mat-card-header>
          <mat-card-title>Accounts</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <table mat-table [dataSource]="accounts" class="full-width">
            <ng-container matColumnDef="name">
              <th mat-header-cell *matHeaderCellDef>Name</th>
              <td mat-cell *matCellDef="let account">{{ account.name }}</td>
            </ng-container>

            <ng-container matColumnDef="type">
              <th mat-header-cell *matHeaderCellDef>Type</th>
              <td mat-cell *matCellDef="let account">{{ account.type | titlecase }}</td>
            </ng-container>

            <ng-container matColumnDef="balance">
              <th mat-header-cell *matHeaderCellDef>Balance</th>
              <td mat-cell *matCellDef="let account">
                {{ account.balance | currency:'PLN':'symbol':'1.2-2' }}
              </td>
            </ng-container>

            <ng-container matColumnDef="actions">
              <th mat-header-cell *matHeaderCellDef>Actions</th>
              <td mat-cell *matCellDef="let account">
                <button mat-icon-button color="primary" (click)="startEdit(account)" matTooltip="Edit">
                  <mat-icon>edit</mat-icon>
                </button>
                <button mat-icon-button color="warn" (click)="deleteAccount(account)" matTooltip="Delete">
                  <mat-icon>delete</mat-icon>
                </button>
              </td>
            </ng-container>

            <tr mat-header-row *matHeaderRowDef="displayedColumns"></tr>
            <tr mat-row *matRowDef="let row; columns: displayedColumns;"></tr>
          </table>

          @if (accounts.length === 0) {
            <p class="empty-message">No accounts yet. Add one above.</p>
          }
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .accounts-container {
      padding: 16px;
      display: flex;
      flex-direction: column;
      gap: 16px;
      max-width: 900px;
      margin: 0 auto;
    }
    .account-form {
      display: flex;
      gap: 12px;
      align-items: flex-start;
      flex-wrap: wrap;
    }
    .account-form mat-form-field {
      flex: 1;
      min-width: 180px;
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
    .empty-message {
      text-align: center;
      color: rgba(0,0,0,0.5);
      font-style: italic;
      padding: 24px;
    }
  `],
})
export class AccountsComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private snackBar = inject(MatSnackBar);

  accounts: Account[] = [];
  editingId: number | null = null;
  displayedColumns = ['name', 'type', 'balance', 'actions'];

  form: FormGroup = this.fb.group({
    name: ['', Validators.required],
    type: ['bank', Validators.required],
    balance: [0, [Validators.required, Validators.min(0)]],
  });

  ngOnInit(): void {
    this.loadAccounts();
  }

  loadAccounts(): void {
    this.api.getAccounts().subscribe((data) => {
      this.accounts = data;
    });
  }

  onSubmit(): void {
    if (this.form.invalid) return;
    const value = this.form.value;

    if (this.editingId) {
      this.api.updateAccount(this.editingId, value).subscribe(() => {
        this.snackBar.open('Account updated', 'Close', { duration: 2000 });
        this.cancelEdit();
        this.loadAccounts();
      });
    } else {
      this.api.createAccount(value).subscribe(() => {
        this.snackBar.open('Account created', 'Close', { duration: 2000 });
        this.form.reset({ name: '', type: 'bank', balance: 0 });
        this.loadAccounts();
      });
    }
  }

  startEdit(account: Account): void {
    this.editingId = account.id;
    this.form.patchValue({
      name: account.name,
      type: account.type,
      balance: account.balance,
    });
  }

  cancelEdit(): void {
    this.editingId = null;
    this.form.reset({ name: '', type: 'bank', balance: 0 });
  }

  deleteAccount(account: Account): void {
    if (!confirm(`Delete account "${account.name}"? This cannot be undone.`)) return;
    this.api.deleteAccount(account.id).subscribe(() => {
      this.snackBar.open('Account deleted', 'Close', { duration: 2000 });
      this.loadAccounts();
    });
  }
}
