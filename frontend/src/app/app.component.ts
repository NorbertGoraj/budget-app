import { Component, inject } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive, Router, NavigationEnd } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatDividerModule } from '@angular/material/divider';
import { MatMenuModule } from '@angular/material/menu';
import { MatBadgeModule } from '@angular/material/badge';
import { MatTooltipModule } from '@angular/material/tooltip';
import { filter, map, startWith } from 'rxjs/operators';
import { AsyncPipe } from '@angular/common';

const ROUTE_TITLES: Record<string, string> = {
  '/': 'Dashboard',
  '/accounts': 'Accounts',
  '/transactions': 'Transactions',
  '/import': 'Import CSV',
  '/budgets': 'Budgets',
  '/purchases': 'Planned Purchases',
  '/investments': 'Investments',
  '/debts': 'Debts',
};

@Component({
  selector: 'app-root',
  imports: [
    RouterOutlet, RouterLink, RouterLinkActive, AsyncPipe,
    MatToolbarModule, MatButtonModule, MatIconModule,
    MatSidenavModule, MatListModule, MatDividerModule,
    MatMenuModule, MatBadgeModule, MatTooltipModule,
  ],
  template: `
    <mat-sidenav-container class="app-container">
      <!-- ===== SIDENAV ===== -->
      <mat-sidenav mode="side" opened class="sidenav">

        <!-- Header / Logo -->
        <div class="sidenav-header">
          <mat-icon class="header-icon">account_balance_wallet</mat-icon>
          <span class="header-title">Budget App</span>
        </div>

        <!-- SECTION: Overview -->
        <div class="nav-section-label">Overview</div>
        <mat-nav-list>
          <a mat-list-item routerLink="/" routerLinkActive="nav-active" [routerLinkActiveOptions]="{exact: true}">
            <mat-icon matListItemIcon>dashboard</mat-icon>
            <span matListItemTitle>Dashboard</span>
          </a>
          <a mat-list-item routerLink="/accounts" routerLinkActive="nav-active">
            <mat-icon matListItemIcon>account_balance</mat-icon>
            <span matListItemTitle>Accounts</span>
          </a>
          <a mat-list-item routerLink="/transactions" routerLinkActive="nav-active">
            <mat-icon matListItemIcon>receipt_long</mat-icon>
            <span matListItemTitle>Transactions</span>
          </a>
        </mat-nav-list>

        <mat-divider class="nav-divider"></mat-divider>

        <!-- SECTION: Finance -->
        <div class="nav-section-label">Finance</div>
        <mat-nav-list>
          <a mat-list-item routerLink="/budgets" routerLinkActive="nav-active">
            <mat-icon matListItemIcon>pie_chart</mat-icon>
            <span matListItemTitle>Budgets</span>
          </a>
          <a mat-list-item routerLink="/debts" routerLinkActive="nav-active">
            <mat-icon matListItemIcon>credit_card</mat-icon>
            <span matListItemTitle>Debts</span>
          </a>
          <a mat-list-item routerLink="/import" routerLinkActive="nav-active">
            <mat-icon matListItemIcon>upload_file</mat-icon>
            <span matListItemTitle>Import CSV</span>
          </a>
        </mat-nav-list>

        <mat-divider class="nav-divider"></mat-divider>

        <!-- SECTION: Planning -->
        <div class="nav-section-label">Planning</div>
        <mat-nav-list>
          <a mat-list-item routerLink="/purchases" routerLinkActive="nav-active">
            <mat-icon matListItemIcon>shopping_cart</mat-icon>
            <span matListItemTitle>Purchases</span>
          </a>
          <a mat-list-item routerLink="/investments" routerLinkActive="nav-active">
            <mat-icon matListItemIcon>trending_up</mat-icon>
            <span matListItemTitle>Investments</span>
          </a>
        </mat-nav-list>

        <!-- Bottom user area -->
        <div class="sidenav-footer">
          <div class="footer-avatar">B</div>
          <span class="footer-username">My Budget</span>
        </div>

      </mat-sidenav>

      <!-- ===== MAIN CONTENT ===== -->
      <mat-sidenav-content class="content">

        <!-- Toolbar -->
        <mat-toolbar class="app-toolbar">
          <span class="page-title">{{ pageTitle$ | async }}</span>
          <span class="toolbar-spacer"></span>

          <!-- Notifications -->
          <button mat-icon-button matTooltip="Notifications" class="toolbar-icon-btn">
            <mat-icon>notifications</mat-icon>
            <span class="notif-badge">4</span>
          </button>

          <!-- User menu -->
          <button mat-button [matMenuTriggerFor]="userMenu" class="user-menu-btn">
            <div class="user-avatar">A</div>
            <span class="user-label">Hi, Admin</span>
          </button>
          <mat-menu #userMenu="matMenu" xPosition="before">
            <button mat-menu-item disabled>
              <mat-icon>person</mat-icon>
              <span>Profile</span>
            </button>
            <mat-divider></mat-divider>
            <button mat-menu-item disabled>
              <mat-icon>logout</mat-icon>
              <span>Sign out</span>
            </button>
          </mat-menu>
        </mat-toolbar>

        <!-- Page content -->
        <div class="page-content">
          <router-outlet />
        </div>

      </mat-sidenav-content>
    </mat-sidenav-container>
  `,
  styles: [`
    /* Layout */
    .app-container { height: 100vh; }

    /* ---- Sidenav ---- */
    .sidenav {
      width: 240px;
      background-color: #ffffff;
      display: flex;
      flex-direction: column;
      border-right: 1px solid #e8eaf0;
      box-shadow: 2px 0 8px rgba(0,0,0,0.06);
    }

    .sidenav-header {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 0 20px;
      height: 64px;
      background: linear-gradient(135deg, #3949ab 0%, #5c6bc0 100%);
      color: #fff;
    }
    .header-icon {
      font-size: 26px;
      width: 26px;
      height: 26px;
    }
    .header-title {
      font-size: 17px;
      font-weight: 600;
      letter-spacing: 0.3px;
    }

    .nav-section-label {
      padding: 16px 20px 4px;
      font-size: 10px;
      font-weight: 700;
      letter-spacing: 1.4px;
      text-transform: uppercase;
      color: #9e9e9e;
    }

    .nav-divider {
      border-top-color: #f0f0f0 !important;
      margin: 6px 0 !important;
    }

    /* Nav list items */
    .sidenav ::ng-deep .mat-mdc-list-item {
      color: #546e7a !important;
      border-radius: 8px !important;
      margin: 2px 10px !important;
      height: 44px !important;
    }
    .sidenav ::ng-deep .mat-mdc-list-item .mat-icon {
      color: #90a4ae !important;
    }
    .sidenav ::ng-deep .mat-mdc-list-item:hover {
      background: rgba(63, 81, 181, 0.07) !important;
    }
    .sidenav ::ng-deep .mat-mdc-list-item:hover .mat-icon {
      color: #3f51b5 !important;
    }
    .sidenav ::ng-deep .nav-active {
      background: #3f51b5 !important;
      border-radius: 8px !important;
    }
    .sidenav ::ng-deep .nav-active .mat-icon {
      color: #fff !important;
    }
    .sidenav ::ng-deep .nav-active .mdc-list-item__primary-text {
      color: #fff !important;
      font-weight: 600 !important;
    }

    /* Footer */
    .sidenav-footer {
      margin-top: auto;
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 14px 20px;
      border-top: 1px solid #f0f0f0;
    }
    .footer-avatar {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      background: linear-gradient(135deg, #3949ab, #5c6bc0);
      color: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      font-weight: 600;
    }
    .footer-username {
      font-size: 13px;
      color: #546e7a;
      font-weight: 500;
    }

    /* ---- Toolbar ---- */
    .app-toolbar {
      position: sticky;
      top: 0;
      z-index: 100;
      height: 64px;
      background: linear-gradient(135deg, #3949ab 0%, #5c6bc0 100%);
      color: #fff;
      box-shadow: 0 2px 10px rgba(57, 73, 171, 0.4);
    }
    .page-title {
      font-size: 17px;
      font-weight: 500;
      color: #fff;
      letter-spacing: 0.2px;
    }
    .toolbar-spacer { flex: 1; }

    .toolbar-icon-btn {
      color: rgba(255,255,255,0.9) !important;
      position: relative;
    }
    .notif-badge {
      position: absolute;
      top: 6px;
      right: 6px;
      width: 16px;
      height: 16px;
      border-radius: 50%;
      background: #ff9800;
      color: #fff;
      font-size: 10px;
      font-weight: 700;
      display: flex;
      align-items: center;
      justify-content: center;
      line-height: 1;
    }

    .user-menu-btn {
      color: #fff !important;
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 0 8px !important;
    }
    .user-avatar {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      background: #e91e63;
      color: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      font-weight: 700;
    }
    .user-label {
      font-size: 14px;
      font-weight: 500;
    }

    /* ---- Content ---- */
    .content {
      display: flex;
      flex-direction: column;
      height: 100%;
      background-color: #f5f7fa;
    }
    .page-content {
      padding: 24px;
      flex: 1;
      overflow: auto;
    }
  `],
})
export class AppComponent {
  private router = inject(Router);

  pageTitle$ = this.router.events.pipe(
    filter(e => e instanceof NavigationEnd),
    map((e) => ROUTE_TITLES[(e as NavigationEnd).urlAfterRedirects] ?? 'Budget App'),
    startWith(ROUTE_TITLES[this.router.url] ?? 'Budget App'),
  );
}
