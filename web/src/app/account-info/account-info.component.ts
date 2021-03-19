import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {FormControl, FormGroup} from '@angular/forms';
import {HttpService} from '../services/http.service';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {SharedService} from '../services/shared.service';
import {AuthService} from '../services/auth.service';
import {User} from '../models/user.model';
import {UserService} from '../services/user.service';
import {MatDialog} from '@angular/material';
import {DeleteDialogComponent} from '../shared/dialogs/delete/delete.dialog.component';
import {Subscription} from 'rxjs/Rx';
import {SkycoinAddressModel} from '../models/skycoin-address.model';
import { EditAddressComponent } from '../shared/dialogs/edit-address/edit-address.component';

@Component({
  selector: 'app-account-info',
  templateUrl: './account-info.component.html',
  styleUrls: ['./account-info.component.scss']
})
export class AccountInfoComponent implements OnInit {
  public user: User;
  public keys: string[] = [];
  public lastId = 1;
  private authSubscription: Subscription;
  updateAddress: FormGroup;
  displaySkycoinAddressForm: boolean = false;
  addressHistory: SkycoinAddressModel[];

  constructor(private router: Router, private sharedService: SharedService, private httpService: HttpService,
                private userService: UserService, private dialog: MatDialog, private authService: AuthService) { }

  ngOnInit() {
    this.getAllUserKeys();
    if (this.authService.getUser()) {
      this.user = this.authService.getUser();
    }
    this.authSubscription = this.authService.userInfo.subscribe(user => {
      this.user = user;
    });
    this.updateAddress = new FormGroup({
      skycoinAddress: new FormControl()
    });
  }

  onSubmit() {
  // this.router.navigate(['update-address']);
  this.displaySkycoinAddressForm = true;
  }

  onUpdate() {
    const val = this.updateAddress.value;
    if (val.skycoinAddress) {
      this.userService.updateAddress(val.skycoinAddress);
      this.displaySkycoinAddressForm = false;
    }
  }

  copyToClipboard(key) {
      const event = (e: ClipboardEvent) => {
          e.clipboardData.setData('text/plain', key);
          e.preventDefault();
          // ...('copy', e), as event is outside scope
          document.removeEventListener('copy', event);
      };
      document.addEventListener('copy', event);
      document.execCommand('copy');
      this.sharedService.showSuccess('Link copied to clipboard.');
  }

  private getAllUserKeys(): void {
    this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.USER.Keys).subscribe(
      (data: any) => {
        this.keys = data;
      },
      (err: any) => {
        this.sharedService.showError('Can\'t get api keys from server: ', err.split(': ')[1]);
      }
    );
  }

  public generateApiKey() {
    this.httpService.postToUrl(environment.whitelistService + ApiRoutes.USER.Keys, {}).subscribe(
      (res: string) => {
        this.keys.push(res);
      },
      (err) => {
        this.sharedService.showError('Can\'t create new api key ', err.split(': ')[1]);
      });
  }

  public deleteApiKey(i: number) {
    this.lastId = i;
    const dialogRef: any = this.dialog.open(DeleteDialogComponent, {
      height: '200px',
      width: '600px',
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        const req: any = { body: {key: this.keys[this.lastId]}};
        this.httpService.deleteFromUrl(environment.whitelistService + ApiRoutes.USER.Keys, req).subscribe(
          (res: any) => {
            this.keys.splice(this.lastId, 1);
          },
          (err) => {
            this.sharedService.showError('Can\'t remove api key ', err.split(': ')[1]);
            console.log(err);
          }
        );
      }

    });
  }

  public getUser(): User {
    return this.user;
  }

  public showSkycoinAddressForm(): boolean {
    return this.displaySkycoinAddressForm;
  }

  public getLatestSkycoinAddress(): string|null {
    return this.user.addressHistory && this.user.addressHistory.length > 0
      ? this.user.addressHistory[0].skycoinAddress
      : null;
  }

  public editAddress() {
    this.dialog.open(EditAddressComponent).afterClosed().subscribe(result => {
      if (result) {
        this.userService.updateAddress(result);
      }
    });
  }
}
