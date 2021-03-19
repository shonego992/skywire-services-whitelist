import {Component, OnInit, Inject} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import { DOCUMENT } from '@angular/common';
import {UserService} from '../../services/user.service';
import { environment } from '../../../environments/environment';
import {Subscription} from 'rxjs/Rx';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  private readonly SIGN_UP_URL = environment.signUpURL + '?redirectURL=';
  private readonly RESET_PASS_URL = environment.resetPasswordURL + '?redirectURL=';
  loginForm: FormGroup;
  tokenNeeded = false;
  private tokenSub: Subscription;

  constructor(private userService: UserService,
    @Inject(DOCUMENT) private document: any) { }

  ngOnInit() {
    this.loginForm = new FormGroup({
      username: new FormControl(null, {validators: [this.noWhitespaceValidator]}),
      password: new FormControl('', {validators: [Validators.required, Validators.minLength(8)]}),
      token: new FormControl()
    });
  }

  public noWhitespaceValidator(control: FormControl) {
    const re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{1,}))$/;
    const valid =  re.test((String(control.value || '').trim()).toLowerCase());
    return valid ? null : { 'error': true };
  }

  onSubmit(): void {
    const val = this.loginForm.value;
    if (val.username && val.password) {
      this.userService.login(val.username.trim(), val.password, val.token);
      this.tokenSub = this.userService.tokenNeeded.subscribe((value: any) => {
        this.tokenNeeded = value;
      });
    }
  }

  public signUpURL(): void {
    this.document.location.href = this.SIGN_UP_URL + this.document.location;
  }

  public resetPasswordURL(): void {
    this.document.location.href = this.RESET_PASS_URL + this.document.location;
  }
}
