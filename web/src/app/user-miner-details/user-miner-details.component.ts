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
import {AuthService} from '../services/auth.service';

@Component({
  selector: 'app-user-miner-details',
  templateUrl: './user-miner-details.component.html',
  styleUrls: ['./user-miner-details.component.scss']
})
export class UserMinerDetailsComponent implements OnInit {

  public minerId: number;
  public miner: any;
  public miners: any[];
  public transferMail: string;
  public nodes: NodeKey[] = [{key: ''}];
  public addNodeKeysOnClick = 1;
  public images = [];
  public buttonDisabled:boolean = false;
  public dateCreated: string;
  public dateUpdated: string;
  public dateDisabled: string;

  @ViewChild(NgxImageGalleryComponent) currentImages: NgxImageGalleryComponent;

  // gallery configuration
  conf: GALLERY_CONF = {
    imageOffset: '0px',
    showDeleteControl: false,
    showImageTitle: false,
    inline: true,
    showExtUrlControl: true
  };

  constructor (public activeRoute: ActivatedRoute, private sharedService: SharedService, private authService: AuthService,
               private httpService: HttpService, private router: Router, private dialog: MatDialog) { }

  ngOnInit () {
    this.activeRoute.queryParams.subscribe(params => {
      this.minerId = params['id'];
      this.miners = this.sharedService.getUserMiners();
      if (this.miners) {
        for (let m of this.miners) {
          if (m.id === this.minerId) {
            this.miner = m;
          }
        }
      } else {
        this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.WHITELIST.Miner + '?id=' + this.minerId).subscribe(
          (data: any) => {
            this.miner = data;
            this.images = this.transformImages();
            this.dateCreated = this.getDateValue(this.miner.createdAt);
            this.dateUpdated = this.getDateValue(this.miner.updatedAt);
            this.dateDisabled = this.getDateValue(this.miner.deletedAt);
          },
          (err: any) => {
            this.sharedService.showError('Can\'t load miner data from server: ', err.split(': ')[1]);
          }
        );
      }
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

  public transferMiner (): void {
    const dialogRef = this.dialog.open(DeleteDialogComponent, {
      data: {username: this.transferMail, message: 'MINER.SURE_TO_TRANSFER'}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result === 1) {
        const reqData = {
          minerId: this.minerId,
          transferTo: this.transferMail
        };
        this.httpService.postToUrl<any>(environment.whitelistService + ApiRoutes.WHITELIST.TransferMiner, reqData).subscribe(
          (data: any) => {
            this.sharedService.showSuccess('Miner transferred.');
            this.router.navigate(['/user-miners']);
          },
          (err: any) => {
            this.sharedService.showError('Can\'t transfer miner', err.split(': ')[1]);
          }
        );
      }
    });
  }

  public checkMail (): boolean {
    return !this.transferMail;
  }

  public checkIfDisabled () {
    return this.miner && this.miner.nodes.length === 0;
  }

  public saveMinerChanges () {
    if (this.miner.type === 0 || this.miner.nodes.length <= this.miner.approvedNodeCount) {
      this.makeTheCall();
    } else {
      const dialogRef = this.dialog.open(EditMinerComponent, {
        data: {username: this.miner.node, message: 'MINER.ADDED_NEW_NODE_KEYS'}
      });
      dialogRef.afterClosed().subscribe(result => {
        if (result === 1) {
          this.makeTheCall(window.location.reload);
        }
      });
    }
  }

  private makeTheCall(callback = null): void {
    this.buttonDisabled = true;
    const data = {
      Id: this.miner.id + '',
      Nodes: this.miner.nodes
    };
    this.httpService.postToUrl(environment.whitelistService + ApiRoutes.WHITELIST.Miner, data).subscribe((res: any) => {
      console.log(data);
      this.authService.refreshUserData();
      this.sharedService.showSuccess('Miner updated');
      this.buttonDisabled = false;
      if (callback) {
        callback();
      }
    },
      (err: any) => {
        this.sharedService.showError('Can\'t update miner', err.toString().split(': ')[1]);
        this.buttonDisabled = false;
      }
    );
  }

  public addNodeKey (): void {
    for (let i = 0; i < this.addNodeKeysOnClick; i++) {
      this.miner.nodes.push({key: '', uptime: []});
    }
  }

  public deleteNodeKey (i: number) {
    if (i > -1 && i < this.miner.nodes.length) {
      this.miner.nodes.splice(i, 1);
    }
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

  onValidator(event) {
    this.sharedService.inputValidator(event);
  }
}
