import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';

@Component({
  selector: 'app-root',
  imports: [
    RouterOutlet, RouterLink, RouterLinkActive,
    MatToolbarModule, MatButtonModule, MatIconModule, MatSidenavModule, MatListModule,
  ],
  template: `
    <mat-sidenav-container class="app-container">
      <mat-sidenav mode="side" opened class="sidenav">
        <div class="logo">Budget App</div>
        <mat-nav-list>
          <a mat-list-item routerLink="/" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">
            <mat-icon matListItemIcon>dashboard</mat-icon>
            <span matListItemTitle>Dashboard</span>
          </a>
          <a mat-list-item routerLink="/accounts" routerLinkActive="active">
            <mat-icon matListItemIcon>account_balance</mat-icon>
            <span matListItemTitle>Accounts</span>
          </a>
          <a mat-list-item routerLink="/transactions" routerLinkActive="active">
            <mat-icon matListItemIcon>receipt_long</mat-icon>
            <span matListItemTitle>Transactions</span>
          </a>
          <a mat-list-item routerLink="/import" routerLinkActive="active">
            <mat-icon matListItemIcon>upload_file</mat-icon>
            <span matListItemTitle>Import CSV</span>
          </a>
          <a mat-list-item routerLink="/budgets" routerLinkActive="active">
            <mat-icon matListItemIcon>pie_chart</mat-icon>
            <span matListItemTitle>Budgets</span>
          </a>
          <a mat-list-item routerLink="/purchases" routerLinkActive="active">
            <mat-icon matListItemIcon>shopping_cart</mat-icon>
            <span matListItemTitle>Purchases</span>
          </a>
          <a mat-list-item routerLink="/investments" routerLinkActive="active">
            <mat-icon matListItemIcon>trending_up</mat-icon>
            <span matListItemTitle>Investments</span>
          </a>
          <a mat-list-item routerLink="/debts" routerLinkActive="active">
            <mat-icon matListItemIcon>credit_card</mat-icon>
            <span matListItemTitle>Debts</span>
          </a>
        </mat-nav-list>
      </mat-sidenav>
      <mat-sidenav-content class="content">
        <mat-toolbar color="primary">
          <span>Budget App</span>
        </mat-toolbar>
        <div class="page-content">
          <router-outlet />
        </div>
      </mat-sidenav-content>
    </mat-sidenav-container>
  `,
  styles: [`
    .app-container { height: 100vh; }
    .sidenav { width: 220px; }
    .logo { padding: 16px; font-size: 20px; font-weight: 600; text-align: center; border-bottom: 1px solid #e0e0e0; }
    .content { display: flex; flex-direction: column; height: 100%; }
    .page-content { padding: 24px; flex: 1; overflow: auto; }
    .active { background: rgba(0,0,0,0.08); }
  `],
})
export class AppComponent {}
