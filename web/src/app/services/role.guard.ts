import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, Router} from '@angular/router';

import {AuthService} from './auth.service';

@Injectable()
export class RoleGuard implements CanActivate {

  constructor(private authService: AuthService, private router: Router) { }

  // ADMIN ROLES:
  // flag_vip, create_user, disable_user, review_whitelist, is_admin
  canActivate(route: ActivatedRouteSnapshot) {
    const expectedRole = route.data.expectedRole;
    if (this.authService.isAuth()) {
      switch (expectedRole) {
        case 'is_admin': {
         return this.authService.isAdmin();
        }
        case 'flag_vip': {
          return this.authService.canAddVip();
        }
        case 'review_whitelist': {
          return this.authService.canReveiew();
        }
        case 'can_manipulate_users': {
          return this.authService.canReveiew() || this.authService.canAddVip();
        }
      }
      return true;
    }

    this.router.navigate(['']);
    return false;
  }

}
