import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Issue } from '../models/issue';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {WhitelistApplication} from '../models/application-model';
import {SharedService} from './shared.service';

@Injectable()
export class ImportedUsersService {
  private readonly API_URL = environment.whitelistService + ApiRoutes.WHITELIST.Import;

  dataChange: BehaviorSubject<any[]> = new BehaviorSubject<any[]>([]);
  // Temporarily stores data from dialogs
  dialogData: any;

  constructor (private httpClient: HttpClient, private sharedService: SharedService) {}

  get data(): any[] {
    return this.dataChange.value;
  }

  getDialogData() {
    return this.dialogData;
  }

  /** CRUD METHODS */
  getImportedUsers(): void {
    if (this.sharedService.getImportedUsers()) {
      const data = this.sharedService.getImportedUsers();
      this.dataChange.next(data);

    } else {
      this.httpClient.get<WhitelistApplication[]>(this.API_URL).subscribe(data => {
          for (const item of data) {
            item.currentStatus = this.sharedService.mapStatus(item.currentStatus);
          }

          this.dataChange.next(data);
        },
        (error: HttpErrorResponse) => {
          console.log (error.name + ' ' + error.message);
        });
    }

  }

  addIssue (issue: Issue): void {
    this.dialogData = issue;
  }

  updateIssue (issue: Issue): void {
    this.dialogData = issue;
  }

  deleteIssue (id: number): void {
    console.log(id);
  }
}
