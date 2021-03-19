import {Injectable} from '@angular/core';
import {ToastrService} from 'ngx-toastr';

const TOAST_TIMEOUT = 3000;
const pattern = /^[a-zA-Z0-9]*$/;

@Injectable()
export class SharedService {

  public constructor(private toastr: ToastrService) {
  }

  private miners: any[];
  private importedUsers: any[];

  public setImportedUsers(users: any[]): void {
    this.importedUsers = users;
  }

  public getImportedUsers(): any [] {
    return this.importedUsers;
  }

  public setUserMiners(miners: any[]): void {
    this.miners = miners;
  }

  public getUserMiners(): any[] {
    return this.miners;
  }

  public sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  public showSuccess (message: string) {
    console.info('success: ', message);
    this.toastr.success(message, '', {
      timeOut: TOAST_TIMEOUT
    });
  }

  public showError (title: string, message: string) {
    console.error('error: ', title, message);
    this.toastr.error(message, title, {
      timeOut: TOAST_TIMEOUT
    });
  }

  // Map the status from the backend, and return if the user has application in progress or can create a new one
  // Statuses from backend:
  // PENDING = 0
  // APPROVED = 1
  // DECLINED = 2
  // DISABLED = 3
  // CANCELED = 4
  public checkIfApplicationIsInProgress(status: number): boolean {
    if (status === 0 || status === 2) {
      return true;
    }
    return false;
  }

  public mapStatus(value: number): string {
    switch (value) {
      case 0: {
        return 'PENDING';
      }
      case 1: {
        return 'APPROVED';
      }
      case 2: {
        return 'DECLINED';
      }
      case 3: {
        return 'DISABLED';
      }
      case 4: {
        return 'CANCELED';
      }
      default: {
        return 'UNKNOWN';
      }
    }
  }

  // Map miner type 0 as official, 1 as DIY
  public mapMinerType(value: number): string {
    switch (value) {
      case 0: {
        return 'OFFICIAL';
      }
      case 1: {
        return 'DIY';
      }
      default: {
        return 'UNKNOWN';
      }
    }
  }

  public inputValidator(event: any) {
    if (!pattern.test(event.target.value)) {
      event.target.value = event.target.value.replace(/[^a-zA-Z0-9]/g, '');
    }
  }

}
