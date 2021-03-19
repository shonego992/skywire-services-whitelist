import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '../../../../../node_modules/@angular/router';
import { HttpService } from '../../../services/http.service';
import { SharedService } from '../../../services/shared.service';
import { environment } from '../../../../environments/environment';
import { ApiRoutes } from '../../routes';
import { User } from '../../../models/user.model';

@Component({
  selector: 'app-edit-admin',
  templateUrl: './edit-admin.component.html',
  styleUrls: ['./edit-admin.component.scss']
})
export class EditAdminComponent implements OnInit {

  public adminEmail: number;
  public admin: User;
  public isDataAvailable: boolean = false;

  rights: Right[] = [
    { value: 'create-user', viewValue: 'Create Admin'},
    { value: 'disable-user', viewValue: 'Disable User'},
    { value: 'flag-vip', viewValue: 'Flag User As VIP' },
    { value: 'review-whitelist', viewValue: 'Review Whitelist Application' }
  ];

  constructor(public activeRoute: ActivatedRoute,
              public httpService: HttpService,
              public sharedService: SharedService) { }

  ngOnInit() {
    this.activeRoute.queryParams.subscribe(params => {
      this.adminEmail = params['email'];
      this.httpService.getFromUrl(environment.userService + ApiRoutes.ADMIN.Users + '/' + this.adminEmail).subscribe(
        (data: User) => {
          this.admin = data;
          this.isDataAvailable = true;
        },
        (err: any) => {
          this.sharedService.showError('Can\'t load admin data from server: ', err.split(': ')[1]);
        }
      );
    });
  }

  compareWithFunc(a, b) {
    return a === b;
  }

  saveEdit() {
    this.httpService.postToUrl(environment.userService + ApiRoutes.ADMIN.Admins + '/' + this.adminEmail + '/rights', this.admin).subscribe(res => {
      console.log(res);
    }
    );
  }

  onCancel() {
    window.close();
  }

}
export interface Right {
  value: string;
  viewValue: string;
}
