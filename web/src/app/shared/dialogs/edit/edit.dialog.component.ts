import {Component, OnInit, ViewChild} from '@angular/core';
import {FormControl, Validators} from '@angular/forms';
import {ActivatedRoute, Router} from '@angular/router';

import {DataService} from '../../../services/data.service';
import {ApiRoutes} from '../../routes';
import {environment} from '../../../../environments/environment';
import {WhitelistApplication} from '../../../models/application-model';
import {HttpService} from '../../../services/http.service';
import {SharedService} from '../../../services/shared.service';
import {ChangeApplicationStatusReq} from '../../../models/requests/admin-change-application-status-req';
import {GALLERY_CONF, GALLERY_IMAGE, NgxImageGalleryComponent} from 'ngx-image-gallery';

@Component({
  selector: 'app-edit',
  templateUrl: './edit.component.html',
  styleUrls: ['./edit.component.scss']
})
export class AdminEditWhitelistComponent implements OnInit {

  public whitelistId: number;
  public application: WhitelistApplication = new WhitelistApplication();
  public currentStatus: string;
  public dateCreated: string;
  public viewingApplication: any = {};
  public viewingApplicationIndex: number;
  public newUserComment: string;
  public newAdminComment: string;
  public whitelistAction = '1';
  public multipleChanges = false;
  public images = [];
  public username = '';
  public currentNodeNumbers;

  @ViewChild(NgxImageGalleryComponent) currentImages: NgxImageGalleryComponent;

  // gallery configuration
  conf: GALLERY_CONF = {
    imageOffset: '0px',
    showDeleteControl: false,
    showImageTitle: false,
    inline: true,
    showExtUrlControl: true
  };

  constructor(public dataService: DataService,
              public router: Router, public activeRoute: ActivatedRoute, public httpService: HttpService,
              public sharedService: SharedService) { }

  formControl = new FormControl('', [
    Validators.required
    // Validators.email,
  ]);

  // load id from url, and get the full details about the application with that ID
  ngOnInit(): void {
    this.activeRoute.queryParams.subscribe(params => {
      this.whitelistId = params['id'];
      this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.WHITELIST.Whitelist + '?id=' + this.whitelistId).subscribe(
        (data: any) => {
          this.application = data;
          this.currentStatus = this.mapStatus(this.application.currentStatus);
          this.username = this.application.userId;
          this.dateCreated = this.getDateValue(this.application.createdAt);
          this.viewingApplication = this.application.changeHistory[this.application.changeHistory.length - 1];
          this.viewingApplicationIndex = this.application.changeHistory.length - 1;
          this.images = this.transformImages(this.application, this.application.changeHistory.length - 1);
          if (this.application.miner.nodes !== null) {
            this.currentNodeNumbers = this.application.miner.nodes.length;
          } else {
            this.currentNodeNumbers = 0;
          }
          if (this.application.changeHistory.length > 1) {
            this.multipleChanges = true;
          }
        },
        (err: any) => {
          this.sharedService.showError('Can\'t load whitelist application data from server: ', err.split(': ')[1]);
        }
      );
    });
  }

  public transformImages(application: WhitelistApplication, index: number): any[] {
    const result: any[] = [];
    if(application.changeHistory[index].images !== null) {
      for (const image of application.changeHistory[index].images) {
        const changed = {
          url: environment.imageBaseURL + image.path,
          extUrl: environment.imageBaseURL + image.path
        };
        result.push(changed);
    }
    }
        return result;  
  }

  getErrorMessage() {
    return this.formControl.hasError('required') ? 'Required field' :
      this.formControl.hasError('email') ? 'Not a valid email' :
        '';
  }

  public mapStatus(value: number): string {
    return this.sharedService.mapStatus(value);
  }

  public getDateValue(value: string): string {
    const time = new Date(value);
    return time.toUTCString();
  }

  public checkIfSubmitDisabled(): boolean {
    return !(this.newAdminComment && this.newAdminComment.length !== 0);
  }

  public nextChangeHistory(): void {
    this.viewingApplicationIndex++;
    this.viewingApplication = this.application.changeHistory[this.viewingApplicationIndex];
    this.images = this.transformImages(this.application, this.viewingApplicationIndex);
  }

  public previousChangeHistory(): void {
    this.viewingApplicationIndex--;
    this.viewingApplication = this.application.changeHistory[this.viewingApplicationIndex];
    this.images = this.transformImages(this.application, this.viewingApplicationIndex);
  }

  public setViewingIndex(i: number) {
    this.viewingApplicationIndex = i;
    this.viewingApplication = this.application.changeHistory[i];
    this.images = this.transformImages(this.application, this.viewingApplicationIndex);
  }


  closeEditWhitelist(): void {
    window.close();
  }

  // stopEdit(): void {
  //   this.dataService.updateIssue(this);
  // }

  updateStatus(): void {
    const body: ChangeApplicationStatusReq = new ChangeApplicationStatusReq();
    body.applicationId = Number(this.whitelistId);
    body.status = Number(this.whitelistAction);
    body.userComment = this.newUserComment;
    body.adminComment = this.newAdminComment;

    this.httpService.postToUrl(environment.whitelistService + ApiRoutes.WHITELIST.Whitelist, body).subscribe(res => {
      this.sharedService.showSuccess('Action done');
      this.router.navigate(['whitelist-app']);
    },
    (err: any) => {
      this.sharedService.showError('Error in managing whitelist', err.split(': ')[1]);
    });
  }


  // EVENTS for image gallery
  // callback on gallery closed
  galleryClosed(value) {
    document.getElementById(value).classList.add('inline');
  }

  // callback on gallery image clicked
  galleryImageClicked(index, value) {
   if (document.getElementById(value).classList.contains('inline')) {
     document.getElementById(value).classList.remove('inline');
     this.currentImages.open(index);
   }
  }
}
