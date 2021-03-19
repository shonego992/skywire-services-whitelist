import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Issue } from '../models/issue';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {WhitelistApplication} from '../models/application-model';
import {SharedService} from './shared.service';

@Injectable()
export class MinersService {
  private readonly API_URL = environment.whitelistService + ApiRoutes.WHITELIST.Miners;

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
  getMiners(): void {
    this.httpClient.get<any[]>(this.API_URL).subscribe(data => {
        // for (let item of data) {
        //   item.currentStatus = this.sharedService.mapStatus(item.currentStatus);
        // }
        this.sharedService.setUserMiners(data);
        this.dataChange.next(data);
      },
      (error: HttpErrorResponse) => {
        console.log (error.name + ' ' + error.message);
      });
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
