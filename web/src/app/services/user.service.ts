import {Injectable} from '@angular/core';

import {environment} from '../../environments/environment';
import {HttpService} from './http.service';
import {ApiRoutes} from '../shared/routes';
import {SharedService} from './shared.service';
import {AuthService} from './auth.service';
import {AuthData} from '../models/auth-data.model';
import {Router} from '@angular/router';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Subject} from 'rxjs/Rx';

@Injectable()
export class UserService {

  constructor(private httpService: HttpService,  private sharedService: SharedService, private authService: AuthService,
              private router: Router, private httpClient: HttpClient) {}
  public tokenNeeded: Subject<any> = new Subject<any>();

  public registerUser (authData: AuthData) {
    this.httpService.postToUrl<AuthData>(environment.userService + ApiRoutes.USER.Users, authData).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Registration successful. Please check your email to confirm profile');
        this.sharedService.sleep(1000);
        this.router.navigate(['']);
      },
      (err: any) => {
        this.sharedService.showError('Can\'t sign up new user', err.split(': ')[1]);

        // TODO  customize messages
        // switch (err) {
        //   case "user service: provided email is already taken by another user": {this.showMessage('Can\'t register new user', err.split(': ')[1]); break;}
        //   case "user service: provided email is already taken by another user": {this.showMessage('Can\'t register new user', err.split(': ')[1]); break;}
        // }
      }
    );
  }

  public login (username: string, password: string, token?: string) {
    if (!token) {
      token = '';
    }
    const data = {
      username: username,
      password: password
    };
    this.httpClient
      .post(environment.userService + '/auth/login', data, {
        headers: new HttpHeaders().set('2fa', token),
      })
      .subscribe((res: any) => {
          if (res && res.token && res.expire) {
            this.authService.saveToken(res);
            this.authService.authSuccessfully();
            this.authService.refreshUserData();
          } else {
            this.sharedService.showError('Unexpected error on sign in', res);
          }
        },
        (err: any) => {
          var errString = err.split(': ')[1];
          if (errString) {
            if (errString.indexOf('2FA') !== -1) {
              this.tokenNeeded.next(true);
              return;
            }
            this.sharedService.showError('Can\'t sign in', errString);
          }
        }
      );
  }

  public updateAddress (newAddress: string) {
    const data = {
      address: newAddress
    };
    this.httpService.patchToUrl<any>(environment.whitelistService + ApiRoutes.USER.Address, data).subscribe(
      (res: any) => {
        this.sharedService.showSuccess('Address successfully updated');
        this.authService.refreshUserData();
      },
      (err: any) => {
        this.sharedService.showError('Can\'t update user\'s address', err.split(': ')[1]);
        // TODO  customize messages
      }
    );
  }

}
