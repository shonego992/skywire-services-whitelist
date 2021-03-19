import {Component, EventEmitter, OnInit, HostListener} from '@angular/core';

import {HttpService} from '../services/http.service';
import {environment} from '../../environments/environment';
import {SharedService} from '../services/shared.service';
import {ApiRoutes} from '../shared/routes';
import {ApplicationWhitelistReq} from '../models/requests/application-whitelist-req';
import {NodeKey} from '../models/requests/node-keys';
import {humanizeBytes, UploaderOptions, UploadFile, UploadInput, UploadOutput} from 'ngx-uploader';
import {AuthService} from '../services/auth.service';
import {User} from '../models/user.model';
import {Subscription} from 'rxjs/Rx';
import {WhitelistApplication} from '../models/application-model';
import {ChangeHistory} from '../models/change-history-model';
import {WhitelistImageModel} from '../models/whitelist-image-model';
import {Router} from '@angular/router';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-whitelist-form',
  templateUrl: './whitelist-form.component.html',
  styleUrls: ['./whitelist-form.component.scss']
})
export class WhitelistFormComponent implements OnInit {
  private maxFiles: number = 5;
  private maxSize: number = 3 * 1024; //KB

  public nodes: NodeKey[] = [{key: ''}];
  public addNodeKeysOnClick = 1;
  public location: string;
  public description: string;
  public user: User;
  public applicationStatus: string;
  public applicationInProgress: boolean = null;
  private authSubscription: Subscription;
  public images: WhitelistImageModel[] = [];
  public wrongKeys: number[] = [];
  public userComment: string;
  public buttonDisabled:boolean = false;
  public loading = false;
  urls = new Array<string>();
  public dateCreated: string;
  public lastUpdate: string;

  // For file upload
  options: UploaderOptions;
  files: UploadFile[];
  uploadInput: EventEmitter<UploadInput>;
  humanizeBytes: Function;
  dragOver: boolean;

  private readonly statusMap = {
    'PENDING': ['orange', 'stop'],
    'APPROVED': ['green', 'done'],
    'DECLINED': ['red', 'close'],
    'DISABLED': ['red', 'close'],
    'CANCELED': ['red', 'close'],
  };

  get statusClass() {
    return this.statusMap[this.applicationStatus][0];
  }

  get statusIcon() {
    return this.statusMap[this.applicationStatus][1];
  }

  constructor(private httpService: HttpService, private sharedService: SharedService, private authService: AuthService,
                private router: Router, private http: HttpClient) {
    this.options = { concurrency: 5, maxUploads: 5 };
    this.files = []; // local uploading files array
    this.uploadInput = new EventEmitter<UploadInput>(); // input events, we use this to emit data to ngx-uploader
    this.humanizeBytes = humanizeBytes;
  }

  // Used for file upload
  onUploadOutput(output: UploadOutput): void {
    if (output.type === 'allAddedToQueue') {
      // when all files added in queue
    } else if (output.type === 'addedToQueue'  && typeof output.file !== 'undefined') { // add file to array when added
      this.checkFile(output.file);
    } else if (output.type === 'uploading' && typeof output.file !== 'undefined') {
      // update current data in files array for uploading file
      const index = this.files.findIndex(file => typeof output.file !== 'undefined' && file.id === output.file.id);
      this.files[index] = output.file;
    } else if (output.type === 'removed') {
      // remove file from array when removed
      this.files = this.files.filter((file: UploadFile) => file !== output.file);
    } else if (output.type === 'dragOver') {
      this.dragOver = true;
    } else if (output.type === 'dragOut') {
      this.dragOver = false;
    } else if (output.type === 'drop') {
      this.dragOver = false;
    } else if (output.type === 'rejected' && typeof output.file !== 'undefined') {
      console.log(output.file.name + ' rejected');
    } else if (output.type === 'done') {

    }
  }

  private checkFile(file: any): void {
    if (file.type.includes('image')) {
      if (file.size / 1024 > this.maxSize) {
        this.sharedService.showError('Error adding file', 'file too big: ' + file.name);
        return;
      }
      let totalFiles = 0;
      if (this.files) {
        totalFiles += this.files.length;
      }
      if (this.images) {
        totalFiles += this.images.length;
      }
      if (totalFiles >= 5) {
        this.sharedService.showError('Error adding file', 'cannot add file to queue, too many files: ' + file.name);
        return;
      }
      this.files.push(file);
    } else {
      this.sharedService.showError('Error adding file', 'file not an image: ' + file.name);
    }
  }

  uploadFiles() {
    this.buttonDisabled = true;
    this.loading = true;
    let headers = { 'Authorization': 'Bearer ' + this.authService.getToken() };
    let url: string = environment.whitelistService + ApiRoutes.WHITELIST.Application;
    if (this.applicationInProgress) {
      url = environment.whitelistService + ApiRoutes.WHITELIST.UpdateApplication;
    }
    let data = this.createApplicationRequestData(this.applicationInProgress, true);
    const formData: FormData = new FormData();

    for (let i = 0; i < this.files.length; i++) {
      const file: any = this.files[i];
      formData.append('file', file.nativeFile, file.name);
    }
    Object.keys(data).forEach(key => formData.append(key, data[key]));
    this.httpService.postToUrlFormData(url , formData, {headers: headers}).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Application submitted. You can review your application and make changes to it.');
        this.buttonDisabled = false;
        this.router.navigate(['/whitelist-form']);
        this.authService.refreshUserData();
        this.files = [];
        this.loading = false;
      },
      (err: any) => {
        this.sharedService.showError('Cannot submit application to server: ', err.split(': ')[1]);
        this.buttonDisabled = false;
        this.loading = false;
      }
  );
  }

  detectFiles(event) {
    this.urls = [];
    let files = event.target.files;
    if (files) {
      for (let file of files) {
        let reader = new FileReader();
        reader.onload = (e: any) => {
          this.urls.push(e.target.result);
        }
        reader.readAsDataURL(file);
      }
    }
  }

  ngOnInit (): void {
    if (this.authService.getUser()) {
      this.user = this.authService.getUser();
      this.getUserApplication();
    }
    this.authSubscription = this.authService.userInfo.subscribe(user => {
      this.user = user;
      this.getUserApplication();
    });
  }

  private getUserApplication(): void {
    if (!this.user.applications || this.user.applications.length === 0) {
      this.applicationInProgress = false;
      return;
    }

    this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.WHITELIST.Application).subscribe(
      (data: any) => {
        if (data.id) {
          const status = data.currentStatus;
          this.applicationInProgress = this.checkIfApplicationIsInProgress(status);
          if (!this.applicationInProgress) {
            return;
          }
        } else {
          this.applicationInProgress = false;
          return;
        }
        this.preloadDataInForm(data);
        this.dateCreated = this.getDateValue(data.createdAt);
        this.lastUpdate = this.getDateValue(data.changeHistory[data.changeHistory.length-1].createdAt);
      },
      (err: any) => {
          this.sharedService.showError('Can\'t get application from server: ', err.split(': ')[1]);
      }
    );
  }

  public getDateValue(value: string): string {
    const time = new Date(value);
    return time.toUTCString();
  }

  private preloadDataInForm(data: WhitelistApplication): void {
    const latestChange: ChangeHistory = data.changeHistory[data.changeHistory.length - 1];
    this.applicationStatus = this.mapApplicationStatusToString(data.currentStatus);
    this.nodes = latestChange.nodes;
    this.location = latestChange.location;
    this.description = latestChange.description;
    this.images = latestChange.images;
    this.userComment = latestChange.userComment;
  }

  // Map the status from the backend, and return if the user has application in progress or can create a new one
  // Statuses from backend:
  // PENDING = 0
  // APPROVED = 1
  // DENIED = 2
  // CANCELED = 3
  private checkIfApplicationIsInProgress(status: number): boolean {
    if (status === 0 || status === 2) {
      return true;
    }
    return false;
  }

  // map numeric value of status to string value for user
  public mapApplicationStatusToString(status: number): string {
    return this.sharedService.mapStatus(status);
  }

  public addNodeKey(): void {
    if (!this.nodes) {
      this.nodes = [];
    }
    for (let i = 0; i < this.addNodeKeysOnClick; i++) {
      this.nodes.push({key: ''});
      this.buttonDisabled = false;
    }
  }

  public getImagePath(image: string): string {
    return environment.imageBaseURL + image;
  }

  public deleteNodeKey(i: number) {
    if (i > -1 && i < this.nodes.length) {
      this.nodes.splice(i, 1);
    }
    if (i === 0) {
      this.buttonDisabled = true;
    }
  }

  public deleteOldImage(index: number): void {
    this.images.splice(index, 1);
    this.urls.splice(index, 1);
  }

  public createAndPostApplicationNoImages(): void {
    const body: ApplicationWhitelistReq = this.createApplicationRequestData(this.applicationInProgress, false);

    // if (this.myFiles.length > 0 && (!this.isValidFiles(this.myFiles))) {
      this.httpService.postToUrl(environment.whitelistService + ApiRoutes.WHITELIST.ApplicationNoImages, body).subscribe(
        res => {console.log(res); }
      );
    // }

  }

  // create data for posting as form update.
  // app in progress shows if application is in progress,
  // stringify if data should be serialized before sending
  private createApplicationRequestData(appInProgress, stringify: boolean): any {
    let oldImages: any = [];
    if (appInProgress) {
      oldImages = stringify ? JSON.stringify(this.images) : this.images;
    }
    const applicationRequest: ApplicationWhitelistReq = {
      location: this.location,
      nodes: stringify ? JSON.stringify(this.nodes) : this.nodes,
      description: this.description,
      oldImages: oldImages
    };
    return applicationRequest;
  }

  public checkIfDisabled(): boolean {
    let totalFiles = 0;
    if (this.files) {
      totalFiles += this.files.length;
    }
    if (this.images) {
      totalFiles += this.images.length;
    }
    const correctNodeKeys = this.checkIfNodeKeysAreCorrect();

    return totalFiles < 3 || !correctNodeKeys;
  }

  private checkIfNodeKeysAreCorrect(): boolean {
    if (this.nodes) {
      for (const node of this.nodes) {
        if (node.key === '') {
          return false;
        }
      }
    } else {
      return false;
    }
    return true;
  }

  onValidator(event) {
    this.sharedService.inputValidator(event);
  }
}


