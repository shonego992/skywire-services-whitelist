import { Component, OnInit } from '@angular/core';
import {HttpService} from '../../services/http.service';
import {environment} from '../../../environments/environment';
import {ApiRoutes} from '../../shared/routes';
import {UploadService} from '../../services/uploader.service';
import {HttpEventType, HttpResponse} from '@angular/common/http';
import {SharedService} from '../../services/shared.service';

@Component({
  selector: 'app-upload-user-list',
  templateUrl: './upload-user-list.component.html',
  styleUrls: ['./upload-user-list.component.scss']
})
export class UploadUserListComponent implements OnInit {

  constructor(private httpService: HttpService, private uploadService: UploadService,
              private sharedService: SharedService) { }

  public showImport: boolean = null;
  private importedData: any;

  ngOnInit() {
    this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.WHITELIST.Import).subscribe(
      (data: any) => {
        this.showImport = false;
        this.importedData = data;
        this.sharedService.setImportedUsers(this.importedData);

      },
      (err: any) => {
        this.showImport = true;
      }
    );

  }

  public selectFile(event) {
    this.uploadFile(event.target.files);
  }

  public uploadFile(files: FileList) {
    if (files.length === 0) {
      console.log('No file selected!');
      return;
    }
    let file: File = files[0];

    this.uploadService.uploadFile(environment.whitelistService + ApiRoutes.WHITELIST.UserList, file)
      .subscribe(
        event => {
          if (event.type === HttpEventType.UploadProgress) {
            const percentDone = Math.round(100 * event.loaded / event.total);
            console.log(`File is ${percentDone}% loaded.`);
          } else if (event instanceof HttpResponse) {
            console.log('File is completely loaded!');
            this.sharedService.showSuccess('Users file imported successfully');
            window.location.reload();//TODO consider workaround for this
          }
        },
        (err) => {
          console.log('Upload Error:', err);
          this.sharedService.showError('Error on import', 'Users file was not imported');
        }, () => {
          console.log('Upload done');
        }
      );
  }

}
