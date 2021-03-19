import {Injectable} from '@angular/core';
import {CanActivate, Router} from '@angular/router';

import {AuthService} from './auth.service';

@Injectable()
export class LoginGuard implements CanActivate {

  constructor(private router: Router, private authService: AuthService) { }

  canActivate() {
    if (this.authService.isAuth()) {
      if (this.authService.isAdmin()) {
        this.router.navigate(['/whitelist-app']);
        return false;
      }
      this.router.navigate(['/user-miners']);
      return false;
    }
    return true;
  }

}
