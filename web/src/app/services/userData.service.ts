import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import { User } from '../models/user.model';

@Injectable()
export class UserDataService {
  private readonly API_URL_USERS = environment.whitelistService + ApiRoutes.ADMIN.Users;

  dataChange: BehaviorSubject<User[]> = new BehaviorSubject<User[]>([]);
  // Temporarily stores data from dialogs
  dialogData: any;

  constructor (private httpClient: HttpClient) {}

  get data(): User[] {
    return this.dataChange.value;
  }

  getDialogData() {
    return this.dialogData;
  }

  /** CRUD METHODS */
  getAllUsers(): void {
    this.httpClient.get<User[]>(this.API_URL_USERS).subscribe(data => {
        this.dataChange.next(data);
      },
      (error: HttpErrorResponse) => {
      console.log (error.name + ' ' + error.message);
      });
  }

  addUser (issue: User): void {
    this.dialogData = issue;
  }

  updateUser (issue: User): void {
    this.dialogData = issue;
  }

  deleteUser (username: string, callback: any): void {
    this.httpClient.delete<User>(this.API_URL_USERS + '/' + username).subscribe(
      data => {
        callback()
      },
      (error: HttpErrorResponse) => {
        console.log(error.name + ' ' + error.message);
      });
  }

}
