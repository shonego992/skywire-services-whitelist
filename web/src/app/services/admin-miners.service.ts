import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Issue } from '../models/issue';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {WhitelistApplication} from '../models/application-model';
import {SharedService} from './shared.service';

@Injectable()
export class AdminMinersService {
  private readonly API_URL = environment.whitelistService + ApiRoutes.WHITELIST.AdminAllMiners;

  dataChange: BehaviorSubject<any[]> = new BehaviorSubject<any[]>([]);
  // Temporarily stores data from dialogs
  dialogData: any;

  constructor (private httpClient: HttpClient, private sharedService: SharedService) {}

  get data(): any[] {
    return this.dataChange.value;
  }

  /** CRUD METHODS */
  getAllMiners(): void {
    this.httpClient.get<any[]>(this.API_URL).subscribe(data => {

        this.dataChange.next(data);
      },
      (error: HttpErrorResponse) => {
        console.log (error.name + ' ' + error.message);
      });
  }
}
