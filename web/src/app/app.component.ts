import {Component, OnInit} from '@angular/core';
import {TranslateService} from '@ngx-translate/core';
import {AuthService} from './services/auth.service';
import {Router} from '@angular/router';
import {environment} from '../environments/environment';
import {User} from './models/user.model';
import {HttpService} from './services/http.service';
import {SharedService} from './services/shared.service';
import {ApiRoutes} from './shared/routes';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit{

  constructor(translate: TranslateService, private authService: AuthService, private router: Router,
              private httpService: HttpService, private sharedService: SharedService) {
    // this language will be used as a fallback when a translation isn't found in the current language
    translate.setDefaultLang('en');

    // the lang to use, if the lang isn't available, it will use the current loader to get them
    translate.use('en');
    if (window.location.href.includes('verify-profile')) {
      return;
    }

    // TODO: figure out a proper way to route for loged in user
    // if (this.authService.isLoggedin()) {
    //   this.router.navigate(['/account-info']);
    // }
  }

    ngOnInit(): void {
      if (this.authService.isAuth()) {
        this.getUserInfo();
      }
    }

  // get user info and store it in auth.service
  private getUserInfo(): void {
    this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.AUTH.Info).subscribe((user: User) => {
        this.authService.setUser(user);
      },
      (err: any) => {
        this.sharedService.showError('Can\'t get user details', err.split(': ')[1]);
      }
    );
  }

}
