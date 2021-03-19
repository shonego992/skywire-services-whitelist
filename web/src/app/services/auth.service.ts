import {Injectable, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Subject} from 'rxjs/Subject';
import 'rxjs/add/operator/map';
import {environment} from '../../environments/environment';
import {HttpService} from './http.service';
import {ApiRoutes} from '../shared/routes';
import {SharedService} from './shared.service';
import * as moment from 'moment';
import {User} from '../models/user.model';
import decode from 'jwt-decode';
import {AdminClaims} from '../models/admin.claims';

const TOKEN_KEY = 'token_key';
const TOKEN_EXPIRE = 'token_expire';

@Injectable()
export class AuthService implements OnInit {

  private userVerified = false;
  private user: User;
  private adminClaims: AdminClaims;
  private returnUrl;

  authChange = new Subject<any>();
  userInfo = new Subject<User>();

  constructor (private router: Router, private route: ActivatedRoute, private httpService: HttpService, private sharedService: SharedService) {}

  ngOnInit() {
  }

  public refreshUserData(): void {
      this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.AUTH.Info).subscribe((user: User) => {
      this.setUser(user);
    },
    (err: any) => {
      this.sharedService.showError('Can\'t get user details', err.split(': ')[1]);
    }
    );
  }

  public getAdminClaims(): AdminClaims {
    return this.adminClaims;
  }

  public getUser(): User {
    return this.user;
  }

  public isAdmin(): boolean {
    this.checkClaims();
    return this.canReveiew() || this.canAddVip();
  }

  public notConfirmed(): boolean {
    this.checkClaims();
    return this.adminClaims ? this.adminClaims.missing_confirmation : false;
  }

  public canReveiew(): boolean {
    this.checkClaims();
    return this.adminClaims ? this.adminClaims.review_whitelist : false;
  }

  public canAddVip(): boolean {
    this.checkClaims();
    return this.adminClaims ? this.adminClaims.flag_vip : false;
  }

  private checkClaims(): void {
    if (!this.adminClaims) {
      this.saveTokenClaims(localStorage.getItem(TOKEN_KEY));
    }
  }

  // TODO: need to set user as an observable and follow it where it is needed
  public setUser(user: User) {
    this.user = user;
    this.userVerified = true;
    this.saveTokenClaims(localStorage.getItem(TOKEN_KEY));
    this.userInfo.next(this.user);
  }

  public isUserVerified() {
    return this.userVerified;
  }

  public setUserVefified(value: boolean) {
   this.userVerified = value;
  }

  public getToken (): string {
    return localStorage.getItem(TOKEN_KEY);
  }

  // Save token data into local storage
  public saveToken (data: any) {
    const expiresAt = moment(data.expire);
    localStorage.removeItem(TOKEN_KEY);
    localStorage.setItem(TOKEN_KEY, data.token);

    localStorage.removeItem(TOKEN_EXPIRE);
    localStorage.setItem(TOKEN_EXPIRE, JSON.stringify(expiresAt.valueOf()));
    this.saveTokenClaims(data.token);
  }

  private saveTokenClaims(token: any): void {
    const tokenPayload: AdminClaims = decode(token);
    this.adminClaims = tokenPayload;
    this.authChange.next({isAuth: true, claims: this.adminClaims});
  }

  private getTokenExpirationDate (): any {
    const expiration =  localStorage.getItem(TOKEN_EXPIRE);
    const expiresAt = JSON.parse(expiration);
    return moment(expiresAt);
  }

  private isTokenExpired (): boolean {
    const currentTime = moment();
    const isExpired = currentTime.isAfter(this.getTokenExpirationDate());
    return isExpired;
  }

  // TODO: check this method and is auth.. for now the same, but include rolse inside?
  public isLoggedin () {
    return localStorage.getItem(TOKEN_KEY) && !this.isTokenExpired();
  }

  public logout (redirect: boolean = true): void {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(TOKEN_EXPIRE);
    localStorage.clear();

    this.adminClaims = new AdminClaims;
    this.authChange.next({ isAuth: false, claims: this.adminClaims });
    if (redirect) {
      this.router.navigate(['/']);
    }
  }

  private refreshToken () {
    return this.httpService.getFromUrl(environment.userService + ApiRoutes.AUTH.Refresh + localStorage.getItem(TOKEN_KEY))
      .map((response: any) => {
        this.saveToken(response);
      });
  }

  getAccount () {
    return this.httpService.getFromUrl(environment.userService + '/account');
  }

  getUserFromUrl () {
    return this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.AUTH.Info);
    // let payload: any = decode(localStorage.getItem(TOKEN_KEY));
    // console.log(payload);
  }

  public isAuth (): boolean {
    const isAuth: boolean = localStorage.getItem(TOKEN_KEY) && !this.isTokenExpired();
    return isAuth;
  }


  public authSuccessfully (): boolean {
    this.returnUrl = this.route.snapshot.queryParams['returnUrl'] || null;
    this.authChange.next({isAuth: true, claims: this.adminClaims});
    if (this.notConfirmed()) {
      this.logout(false);
      this.sharedService.showError('User not verified', 'Can\'t use Whitelisting System without verified user account');
      return false;
    } else if (this.isAdmin()) {
      this.router.navigate([this.returnUrl ? this.returnUrl : '/whitelist-app']);
    } else {
      this.router.navigate([this.returnUrl ? this.returnUrl : '/user-miners']);
    }
    return true;
  }
}
