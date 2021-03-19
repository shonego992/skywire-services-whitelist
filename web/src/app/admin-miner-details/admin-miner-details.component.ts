import {Component, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {SharedService} from '../services/shared.service';
import {HttpService} from '../services/http.service';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {DeleteDialogComponent} from '../shared/dialogs/delete/delete.dialog.component';
import {MatDialog} from '@angular/material';
import {NodeKey} from '../models/requests/node-keys';
import {EditMinerComponent} from '../shared/dialogs/edit-miner/edit-miner.component';
import {GALLERY_CONF, NgxImageGalleryComponent} from 'ngx-image-gallery';

@Component({
  selector: 'app-admin-miner-details',
  templateUrl: './admin-miner-details.component.html',
  styleUrls: ['./admin-miner-details.component.scss']
})
export class AdminMinerDetailsComponent implements OnInit {

  public minerId: number;
  public miner: any;
  public miners: any[];
  public transferMail: string;
  public uptimes: string[] = [];
  public images = [];
  @ViewChild(NgxImageGalleryComponent) currentImages: NgxImageGalleryComponent;
  public dateCreated: string;
  public dateDisabled: string;
  public dateUpdated: string;
  public uptime;


  // gallery configuration
  conf: GALLERY_CONF = {
    imageOffset: '0px',
    showDeleteControl: false,
    showImageTitle: false,
    inline: true,
    showExtUrlControl: true
  };

  constructor (public activeRoute: ActivatedRoute, private sharedService: SharedService, private httpService: HttpService) { }

  ngOnInit () {
    this.activeRoute.queryParams.subscribe(params => {
      this.minerId = params['id'];
        this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.WHITELIST.MinerForAdmin + '?id=' + this.minerId).subscribe(
          (data: any) => {
            this.miner = data;
            // this.getUptimes(this.miner);
            this.images = this.transformImages();
            this.dateCreated = this.getDateValue(this.miner.createdAt);
            this.dateUpdated = this.getDateValue(this.miner.updatedAt);
            this.dateDisabled = this.getDateValue(this.miner.deletedAt);
          },
          (err: any) => {
            this.sharedService.showError('Can\'t load miner data from server: ', err.split(': ')[1]);
          }
        );
    });
  }

  public getDateValue(value: string): string {
    const time = new Date(value);
    return time.toUTCString();
  }

  public transformImages(): any[] {
    const result: any[] = [];
    for (const image of this.miner.images) {
      const changed = {
        url: environment.imageBaseURL + image.path,
        extUrl: environment.imageBaseURL + image.path
      };
      result.push(changed);
    }
    return result;
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

  public checkMail (): boolean {
    return !this.transferMail;
  }

  public checkIfDisabled () {
    return false;
  }

  public close(): void {
    window.close();
  }

  function(totalSeconds) {
    var hours   = Math.floor(totalSeconds / 3600);
    var minutes = Math.floor((totalSeconds - (hours * 3600)) / 60);
    var seconds = totalSeconds - (hours * 3600) - (minutes * 60);

    // round seconds
    seconds = Math.round(seconds * 100) / 100

    var result = (hours < 10 ? "0" + hours : hours);
        result += "-" + (minutes < 10 ? "0" + minutes : minutes);
        result += "-" + (seconds  < 10 ? "0" + seconds : seconds);
    return result;
  }
}
