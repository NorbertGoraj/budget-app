import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ApiService } from '../services/api.service';
import { Account } from '../services/models';

interface CsvRow {
  [key: string]: string;
}

@Component({
  selector: 'app-import',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatSelectModule,
    MatCardModule,
    MatTableModule,
    MatProgressBarModule,
    MatSnackBarModule,
  ],
  template: `
    <div class="import-container">
      <h2>CSV Import</h2>

      <!-- Account Selector -->
      <mat-card class="config-card">
        <mat-card-header>
          <mat-card-title>Import Configuration</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="form" class="config-form">
            <mat-form-field>
              <mat-label>Import into Account</mat-label>
              <mat-select formControlName="account_id">
                @for (account of accounts(); track account.id) {
                  <mat-option [value]="account.id">{{ account.name }} ({{ account.type }})</mat-option>
                }
              </mat-select>
            </mat-form-field>
          </form>
        </mat-card-content>
      </mat-card>

      <!-- File Upload Area -->
      <mat-card class="upload-card">
        <mat-card-content>
          <div class="drop-zone"
               [class.drag-over]="isDragOver()"
               (dragover)="onDragOver($event)"
               (dragleave)="onDragLeave($event)"
               (drop)="onDrop($event)"
               (click)="fileInput.click()">
            <input #fileInput type="file" accept=".csv" hidden
                   (change)="onFileSelected($event)" />
            <mat-icon class="upload-icon">cloud_upload</mat-icon>
            @if (selectedFile()) {
              <p class="file-name">{{ selectedFile()!.name }}</p>
              <p class="file-size">{{ (selectedFile()!.size / 1024).toFixed(1) }} KB</p>
            } @else {
              <p class="drop-text">Drag & drop a CSV file here, or click to browse</p>
            }
          </div>
        </mat-card-content>
      </mat-card>

      <!-- CSV Preview -->
      @if (previewRows().length > 0) {
        <mat-card class="preview-card">
          <mat-card-header>
            <mat-card-title>Preview ({{ previewRows().length }} rows)</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <div class="table-wrapper">
              <table mat-table [dataSource]="previewRows()" class="full-width">
                @for (col of previewColumns(); track col) {
                  <ng-container [matColumnDef]="col">
                    <th mat-header-cell *matHeaderCellDef>{{ col }}</th>
                    <td mat-cell *matCellDef="let row">{{ row[col] }}</td>
                  </ng-container>
                }
                <tr mat-header-row *matHeaderRowDef="previewColumns()"></tr>
                <tr mat-row *matRowDef="let row; columns: previewColumns();"></tr>
              </table>
            </div>

            <div class="import-actions">
              <button mat-raised-button color="primary"
                      [disabled]="form.invalid || !selectedFile() || importing()"
                      (click)="importCsv()">
                <mat-icon>file_upload</mat-icon> Import {{ previewRows().length }} Rows
              </button>
              <button mat-button (click)="clearFile()">
                <mat-icon>clear</mat-icon> Clear
              </button>
            </div>

            @if (importing()) {
              <mat-progress-bar mode="indeterminate"></mat-progress-bar>
            }
          </mat-card-content>
        </mat-card>
      }

      <!-- Import Results -->
      @if (importResult()) {
        <mat-card class="results-card">
          <mat-card-header>
            <mat-card-title>Import Results</mat-card-title>
          </mat-card-header>
          <mat-card-content>
            <div class="results">
              <div class="result-item success">
                <mat-icon>check_circle</mat-icon>
                <span>{{ importResult()!.imported }} imported</span>
              </div>
              <div class="result-item skipped">
                <mat-icon>skip_next</mat-icon>
                <span>{{ importResult()!.skipped }} skipped</span>
              </div>
            </div>
          </mat-card-content>
        </mat-card>
      }
    </div>
  `,
  styles: [`
    .import-container {
      padding: 24px;
      max-width: 1000px;
      margin: 0 auto;
    }
    .config-card, .upload-card, .preview-card, .results-card {
      margin-bottom: 24px;
    }
    .config-form {
      padding-top: 16px;
    }
    .config-form mat-form-field {
      width: 100%;
      max-width: 400px;
    }
    .drop-zone {
      border: 2px dashed #bdbdbd;
      border-radius: 12px;
      padding: 48px 24px;
      text-align: center;
      cursor: pointer;
      transition: all 0.2s ease;
      background: #fafafa;
    }
    .drop-zone:hover, .drop-zone.drag-over {
      border-color: #1565c0;
      background: #e3f2fd;
    }
    .upload-icon {
      font-size: 48px;
      width: 48px;
      height: 48px;
      color: #9e9e9e;
    }
    .drag-over .upload-icon {
      color: #1565c0;
    }
    .drop-text {
      color: #757575;
      margin-top: 8px;
    }
    .file-name {
      font-weight: 500;
      font-size: 16px;
      margin-top: 8px;
    }
    .file-size {
      color: #757575;
      font-size: 13px;
    }
    .table-wrapper {
      overflow-x: auto;
      max-height: 400px;
      overflow-y: auto;
    }
    .full-width {
      width: 100%;
    }
    .import-actions {
      display: flex;
      gap: 12px;
      align-items: center;
      margin-top: 16px;
    }
    .results {
      display: flex;
      gap: 32px;
      padding: 16px 0;
    }
    .result-item {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 18px;
      font-weight: 500;
    }
    .result-item.success {
      color: #2e7d32;
    }
    .result-item.success mat-icon {
      color: #2e7d32;
    }
    .result-item.skipped {
      color: #f57f17;
    }
    .result-item.skipped mat-icon {
      color: #f57f17;
    }
  `],
})
export class ImportComponent implements OnInit {
  private api = inject(ApiService);
  private fb = inject(FormBuilder);
  private snackBar = inject(MatSnackBar);

  accounts = signal<Account[]>([]);
  selectedFile = signal<File | null>(null);
  previewRows = signal<CsvRow[]>([]);
  previewColumns = signal<string[]>([]);
  isDragOver = signal(false);
  importing = signal(false);
  importResult = signal<{ imported: number; skipped: number } | null>(null);

  form = this.fb.group({
    account_id: [null as number | null, Validators.required],
  });

  ngOnInit(): void {
    this.loadAccounts();
  }

  loadAccounts(): void {
    this.api.getAccounts().subscribe({
      next: (data) => this.accounts.set(data),
      error: () => this.snackBar.open('Failed to load accounts', 'Close', { duration: 3000 }),
    });
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragOver.set(true);
  }

  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragOver.set(false);
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragOver.set(false);
    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      const file = files[0];
      if (file.name.endsWith('.csv')) {
        this.handleFile(file);
      } else {
        this.snackBar.open('Please select a CSV file', 'Close', { duration: 3000 });
      }
    }
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.handleFile(input.files[0]);
      input.value = '';
    }
  }

  private handleFile(file: File): void {
    this.selectedFile.set(file);
    this.importResult.set(null);
    this.parseCsv(file);
  }

  private parseCsv(file: File): void {
    const reader = new FileReader();
    reader.onload = (e) => {
      const text = e.target?.result as string;
      if (!text) return;

      const lines = text.split('\n').filter((line) => line.trim().length > 0);
      if (lines.length < 2) {
        this.snackBar.open('CSV file must have a header row and at least one data row', 'Close', { duration: 3000 });
        return;
      }

      const headers = this.parseCsvLine(lines[0]);
      this.previewColumns.set(headers);

      const rows: CsvRow[] = [];
      for (let i = 1; i < lines.length; i++) {
        const values = this.parseCsvLine(lines[i]);
        const row: CsvRow = {};
        headers.forEach((header, idx) => {
          row[header] = values[idx] || '';
        });
        rows.push(row);
      }
      this.previewRows.set(rows);
    };
    reader.readAsText(file);
  }

  private parseCsvLine(line: string): string[] {
    const result: string[] = [];
    let current = '';
    let inQuotes = false;
    for (let i = 0; i < line.length; i++) {
      const char = line[i];
      if (inQuotes) {
        if (char === '"') {
          if (i + 1 < line.length && line[i + 1] === '"') {
            current += '"';
            i++;
          } else {
            inQuotes = false;
          }
        } else {
          current += char;
        }
      } else {
        if (char === '"') {
          inQuotes = true;
        } else if (char === ',') {
          result.push(current.trim());
          current = '';
        } else {
          current += char;
        }
      }
    }
    result.push(current.trim());
    return result;
  }

  clearFile(): void {
    this.selectedFile.set(null);
    this.previewRows.set([]);
    this.previewColumns.set([]);
    this.importResult.set(null);
  }

  importCsv(): void {
    const accountId = this.form.value.account_id;
    const file = this.selectedFile();
    if (!accountId || !file) return;

    this.importing.set(true);
    this.api.importCSV(accountId, file).subscribe({
      next: (result) => {
        this.importing.set(false);
        this.importResult.set({
          imported: result.imported ?? result.count ?? 0,
          skipped: result.skipped ?? 0,
        });
        this.snackBar.open('Import complete', 'Close', { duration: 2000 });
      },
      error: () => {
        this.importing.set(false);
        this.snackBar.open('Import failed', 'Close', { duration: 3000 });
      },
    });
  }
}
