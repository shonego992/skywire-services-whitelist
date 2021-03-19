import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '../../../../../node_modules/@angular/router';
import { HttpService } from 'src/app/services/http.service';
import { environment } from '../../../../environments/environment';
import { ApiRoutes } from '../../routes';
import { User } from '../../../models/user.model';
import { SharedService } from 'src/app/services/shared.service';

@Component({
  selector: 'app-edit-user',
  templateUrl: './edit-user.component.html',
  styleUrls: ['./edit-user.component.scss']
})
export class EditUserComponent implements OnInit {

  public userId: string;
  public user: User;
  public createdAt: string;

  constructor(public activeRoute: ActivatedRoute, public httpService: HttpService, public sharedService: SharedService) { }

  ngOnInit() {
    const _self = this
    this.activeRoute.queryParams.subscribe(params => {
      _self.userId = params['userId'];
    });
    this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.ADMIN.Users + '/' + this.userId).subscribe(
      (data: User) => {
        this.user = data;
        this.createdAt = this.getDateValue(this.user.createdAt);
      }
    );
  }
  public getDateValue(value: string): string {
    const time = new Date(value);
    return time.toUTCString();
  }

  public getLatestSkycoinAddress(): string|null {
    return this.user.addressHistory && this.user.addressHistory.length > 0
      ? this.user.addressHistory[0].skycoinAddress
      : null;
  }
}
