import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {Router} from '@angular/router';
import {HttpClient} from '@angular/common/http';
import {MatDialog, MatPaginator, MatSort} from '@angular/material';
import {DropdownOption, Issue} from '../../models/issue';
import {Observable} from 'rxjs/Observable';
import {BehaviorSubject} from 'rxjs/BehaviorSubject';
import {DataSource} from '@angular/cdk/collections';
import 'rxjs/add/observable/merge';
import 'rxjs/add/observable/fromEvent';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import {DeleteDialogComponent} from '../../shared/dialogs/delete/delete.dialog.component';
import {SharedService} from '../../services/shared.service';
import {AdminMinersService} from '../../services/admin-miners.service';
import {environment} from '../../../environments/environment';
import {ApiRoutes} from '../../shared/routes';
import {HttpService} from '../../services/http.service';

@Component({
  selector: 'app-admin-miners-view',
  templateUrl: './admin-all-miners-view.component.html',
  styleUrls: ['./admin-all-miners-view.component.scss']
})
export class AdminAllMinersViewComponent  implements OnInit {
  displayedColumns = ['id', 'username', 'type', 'created_at', 'updated_at', 'label', 'gifted', 'is_active', 'actions'];
  exampleDatabase: AdminMinersService | null;
  dataSource: ExampleDataSource | null;
  id: number;
  startDate: Date = null;
  endDate: Date = null;
  filterActiveState = '';
  activeStates: DropdownOption[] = [
    { value: '', viewValue: 'All' },
    { value: 'ACTIVE', viewValue: 'Active' },
    { value: 'DISABLED', viewValue: 'Deleted' }
  ];

  constructor(public httpClient: HttpClient,
              public dialog: MatDialog,
              public router: Router,
              public sharedService: SharedService,
              public httpService: HttpService) {}

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild('filter') filter: ElementRef;


  public exportMinersData(): void {
    let request = {
      startDate: '',
      endDate: ''
    };
    if (this.startDate && this.endDate) {
      request = {
        startDate: (this.startDate.getTime() / 1000).toString(),
        endDate: (this.endDate.getTime() / 1000).toString()
      };
    }

      this.getStartDate();
      this.getEndDate();

    this.httpClient.post(environment.whitelistService + ApiRoutes.WHITELIST.ExportMiners, request, {responseType: 'text'}).subscribe(
      (data: any) => {
        const options = { type: 'text/csv;charset=utf-8;' };
        const filename = 'minerExport_' + this.formatDate(this.startDate) + '-' + this.formatDate(this.endDate) + '.csv';
        this.createAndDownloadBlobFile(data, options, filename);

      },
      (err: any) => {
        this.sharedService.showError('Can\'t load miners data', err.split(': ')[1]);
      }
    );
  }

  public exportMinersNoRestrictions() {
    let request = {
      startDate: '',
      endDate: ''
    };
    if (this.startDate && this.endDate) {
      request = {
        startDate: (this.startDate.getTime() / 1000).toString(),
        endDate: (this.endDate.getTime() / 1000).toString()
      };
    }
      this.getStartDate();
      this.getEndDate();

    this.httpClient.post(environment.whitelistService + ApiRoutes.WHITELIST.ExportMinersNoLimitations, request, {responseType: 'text'}).subscribe(
      (data: any) => {
        const options = { type: 'text/csv;charset=utf-8;' };
        const filename = 'minerExport_' + this.formatDate(this.startDate) + '-' + this.formatDate(this.endDate) + '_all' + '.csv';
        this.createAndDownloadBlobFile(data, options, filename);

      },
      (err: any) => {
        this.sharedService.showError('Can\'t load miners data', err.split(': ')[1]);
      }
    );
  }

  private getStartDate() {
    if (!this.startDate) {
      this.startDate = new Date();
      this.startDate.setDate(1);
      this.startDate.setMonth(this.startDate.getMonth() - 1);
     }
  }

  private getEndDate() {
    if (!this.endDate) {
      this.endDate = new Date();
      this.endDate.setDate(0);
     }
  }

  private formatDate(thisDate : Date): string{
   let stringDate = thisDate.getDate() + "_" + (thisDate.getMonth() + 1) + "_" + thisDate.getFullYear()
    return stringDate
  }

  private createAndDownloadBlobFile(body, options, filename) {
    const blob = new Blob([body], options);
    if (navigator.msSaveBlob) {
      // IE 10+
      navigator.msSaveBlob(blob, filename);
    } else {
      let link = document.createElement('a');
      // Browsers that support HTML5 download attribute
      if (link.download !== undefined) {
        const url = URL.createObjectURL(blob);
        link.setAttribute('href', url);
        link.setAttribute('download', filename);
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      }
    }
  }


  public getDateValue(value: string): string {
    const time = new Date(value);
    return time.toUTCString();
  }
  //
  // PENDING = 0
  // APPROVED = 1
  // DENIED = 2
  // CANCELED = 3
  public mapStatus(value: number): string {
    return this.sharedService.mapStatus(value);
  }

  ngOnInit() {

    this.loadData();
  }

  refresh() {
    this.loadData();
  }

  chooseActiveState(value) {
    this.dataSource.activeState = value;
    this.dataSource._paginator.firstPage();
  }

  public deleteMiner(id: string): void {
    const dialogRef = this.dialog.open(DeleteDialogComponent, {
      data: {username: id, message: 'MINER.DELETE_MINER'}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result === 1) {
        this.httpService.deleteFromUrl(environment.whitelistService + ApiRoutes.WHITELIST.Miner + '/' + id).subscribe(
          (res: any) => {
            const foundIndex = this.exampleDatabase.dataChange.value.findIndex(x => x.id === id);
            this.exampleDatabase.dataChange.value[foundIndex].disabled = true;
            // this.refreshTable();
          },
          (err) => {
            this.sharedService.showError('Can\'t remove miner key ', err.split(': ')[1]);
            console.log(err);
          }
        );
      }
    });
  }

  public activateMiner(id: string): void {
    const dialogRef = this.dialog.open(DeleteDialogComponent, {
      data: {username: id, message: 'MINER.ACTIVATE_MINER'}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result === 1) {
        this.httpService.getFromUrl(environment.whitelistService + ApiRoutes.WHITELIST.Miner + '/' + id + '/activate').subscribe(
          (res: any) => {
            const foundIndex = this.exampleDatabase.dataChange.value.findIndex(x => x.id === id);
            this.exampleDatabase.dataChange.value[foundIndex].disabled = false;
            // this.refreshTable();
          },
          (err) => {
            this.sharedService.showError('Can\'t activate miner ', err.split(': ')[1]);
            console.log(err);
          }
        );
      }
    });
  }

  public mapMinerType(value: number): string {
    return this.sharedService.mapMinerType(value);
  }

  // If you don't need a filter or a pagination this can be simplified, you just use code from else block
  private refreshTable() {
    // if there's a paginator active we're using it for refresh
    if (this.dataSource._paginator.hasNextPage()) {
      this.dataSource._paginator.nextPage();
      this.dataSource._paginator.previousPage();
      // in case we're on last page this if will tick
    } else if (this.dataSource._paginator.hasPreviousPage()) {
      this.dataSource._paginator.previousPage();
      this.dataSource._paginator.nextPage();
      // in all other cases including active filter we do it like this
    } else {
      this.dataSource.filter = '';
      this.dataSource.filter = this.filter.nativeElement.value;
    }
  }

  public loadData() {
    this.exampleDatabase = new AdminMinersService(this.httpClient, this.sharedService);
    this.dataSource = new ExampleDataSource(this.exampleDatabase, this.paginator, this.sort, this.sharedService);
    Observable.fromEvent(this.filter.nativeElement, 'keyup')
      .debounceTime(150)
      .distinctUntilChanged()
      .subscribe(() => {
        if (!this.dataSource) {
          return;
        }
        this.dataSource.filter = this.filter.nativeElement.value;
      });
  }
}

export class ExampleDataSource extends DataSource<any> {
  _filterChange = new BehaviorSubject('');
  _filterStatus = new BehaviorSubject('');

  get filter(): string {
    return this._filterChange.value;
  }

  set filter(filter: string) {
    this._filterChange.next(filter);
  }

  get activeState(): string {
    return this._filterStatus.value;
  }

  set activeState(activeState: string) {
    this._filterStatus.next(activeState);
  }

  filteredData: any[] = [];
  renderedData: any[] = [];

  constructor(public _exampleDatabase: AdminMinersService,
              public _paginator: MatPaginator,
              public _sort: MatSort,
              public sharedService: SharedService) {
    super();
    // Reset to the first page when the user changes the filter.
    this._filterChange.subscribe(() => this._paginator.pageIndex = 0);
  }

  public mapMinerType(value: number): string {
    return this.sharedService.mapMinerType(value);
  }

  /** Connect function called by the table to retrieve one stream containing the data to render. */
  connect(): Observable<any[]> {
    // Listen for any changes in the base data, sorting, filtering, or pagination
    const displayDataChanges = [
      this._exampleDatabase.dataChange,
      this._sort.sortChange,
      this._filterChange,
      this._filterStatus,
      this._paginator.page
    ];

    this._exampleDatabase.getAllMiners();
    return Observable.merge(...displayDataChanges).map(() => {
      // Filter data
      this.filteredData = this._exampleDatabase.data.slice().filter((miner: any) => {
        if (this.activeState && this.activeState.length > 0) {
          if (this.activeState === "ACTIVE" && miner.disabled) {
            return false;
          } else if (this.activeState === "DISABLED" && !miner.disabled) {
            return false;
          }
        }
        const minerType = miner.type == 1 ? 'DIY' : 'OFFICIAL';
        const searchStr = (miner.username.toString() + miner.id.toString() + minerType + miner.batchLabel).toLowerCase();
        return searchStr.indexOf(this.filter.toLowerCase()) !== -1;
      });

      // Sort filtered data
      const sortedData = this.sortData(this.filteredData.slice());

      // Grab the page's slice of the filtered sorted data.
      const startIndex = this._paginator.pageIndex * this._paginator.pageSize;
      this.renderedData = sortedData.splice(startIndex, this._paginator.pageSize);
      return this.renderedData;
    });
  }
  disconnect() {
  }

  /** Returns a sorted copy of the database data. */
  sortData(data: any[]): any[] {
    if (!this._sort.active || this._sort.direction === '') {
      return data;
    }

    return data.sort((a, b) => {
      let propertyA: number | string = '';
      let propertyB: number | string = '';

      switch (this._sort.active) {
        case 'id': [propertyA, propertyB] = [a.id, b.id]; break;
        case 'title': [propertyA, propertyB] = [a.title, b.title]; break;
        case 'state': [propertyA, propertyB] = [a.state, b.state]; break;
        case 'url': [propertyA, propertyB] = [a.url, b.url]; break;
        case 'created_at': [propertyA, propertyB] = [a.createdAt, b.createdAt]; break;
        case 'updated_at': [propertyA, propertyB] = [a.updatedAt, b.updatedAt]; break;
      }

      const valueA = isNaN(+propertyA) ? propertyA : +propertyA;
      const valueB = isNaN(+propertyB) ? propertyB : +propertyB;

      return (valueA < valueB ? -1 : 1) * (this._sort.direction === 'asc' ? 1 : -1);
    });
  }
}
