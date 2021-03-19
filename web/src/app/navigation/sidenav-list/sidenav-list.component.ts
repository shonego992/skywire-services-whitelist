import { Component, OnInit, EventEmitter, Output, OnDestroy} from '@angular/core';
import { Subscription } from 'rxjs/Subscription';

import { AuthService } from '../../services/auth.service';
import {AdminClaims} from '../../models/admin.claims';

@Component({
  selector: 'app-sidenav-list',
  templateUrl: './sidenav-list.component.html',
  styleUrls: ['./sidenav-list.component.scss']
})
export class SidenavListComponent implements OnInit, OnDestroy {
  @Output() linkClick = new EventEmitter();
  public isAuth = false;
  public adminClaims: AdminClaims ;
  private authSubscription: Subscription;

  constructor(private authService: AuthService) { }

  ngOnInit() {
    this.isAuth = this.authService.isAuth();
    this.adminClaims = this.authService.getAdminClaims() || new AdminClaims();
    this.authSubscription = this.authService.authChange.subscribe(authStatus => {
      this.isAuth = authStatus.isAuth;
      this.adminClaims = authStatus.claims;
    });
  }

  clicked() {
    this.linkClick.emit();
  }

  ngOnDestroy() {
    this.authSubscription.unsubscribe();
  }

}
