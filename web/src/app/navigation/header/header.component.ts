import { Component, OnInit, EventEmitter, Output, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';

import { AuthService } from '../../services/auth.service';
import {AdminClaims} from '../../models/admin.claims';
import {User} from '../../models/user.model';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Output() sidenavToggle = new EventEmitter<void>();
  public isAuth = false;
  public adminClaims: AdminClaims ;
  private authSubscription: Subscription;
  private user: User;

  constructor(private authService: AuthService) { }

  ngOnInit() {
    this.isAuth = this.authService.isAuth();
    this.adminClaims = this.authService.getAdminClaims() || new AdminClaims();
    this.authSubscription = this.authService.authChange.subscribe(authStatus => {
      this.isAuth = authStatus.isAuth;
      this.adminClaims = authStatus.claims;
    });

    this.authSubscription = this.authService.userInfo.subscribe(user => {
      this.user = user;
    });
  }

  onToggleSidenav() {
    this.sidenavToggle.emit();
  }

  onLogout() {
    this.authService.logout();
  }

  ngOnDestroy() {
    this.authSubscription.unsubscribe();
  }

}
