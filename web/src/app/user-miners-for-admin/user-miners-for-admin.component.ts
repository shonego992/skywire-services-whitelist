import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Router} from '@angular/router';
import {HttpClient} from '@angular/common/http';
import {MatDialog, MatPaginator, MatSort} from '@angular/material';
import {DropdownOption, Issue} from '../models/issue';
import {Observable} from 'rxjs/Observable';
import {BehaviorSubject} from 'rxjs/BehaviorSubject';
import {DataSource} from '@angular/cdk/collections';
import 'rxjs/add/observable/merge';
import 'rxjs/add/observable/fromEvent';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import {SharedService} from '../services/shared.service';
import {UserMinersForAdmin} from '../services/miners-service-for-admin.service';
import {DeleteDialogComponent} from '../shared/dialogs/delete/delete.dialog.component';
import {ApiRoutes} from '../shared/routes';
import {environment} from '../../environments/environment';
import {HttpService} from '../services/http.service';

@Component({
  selector: 'app-user-miners-for-admin',
  templateUrl: './user-miners-for-admin.component.html',
  styleUrls: ['./user-miners-for-admin.component.scss']
})
export class UserMinersForAdminComponent  implements OnInit {
  displayedColumns = ['id', 'created_at', 'updated_at', 'type', 'label', 'gifted', 'is_active', 'actions'];
  exampleDatabase: UserMinersForAdmin | null;
  dataSource: ExampleDataSource | null;
  index: number;
  id: number;
  userId: string;
  filterActiveState: string = '';
  activeStates: DropdownOption[] = [
    { value: '', viewValue: 'All' },
    { value: 'ACTIVE', viewValue: 'Active' },
    { value: 'DISABLED', viewValue: 'Disabled' }
  ];

  @Input()
  public username: string;

  constructor(public httpClient: HttpClient,
              public dialog: MatDialog,
              public router: Router,
              public sharedService: SharedService,
              public httpService: HttpService) {}

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild('filter') filter: ElementRef;

  public mapMinerType(value: number): string {
    return this.sharedService.mapMinerType(value);
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
    if (this.username) {
      this.userId = this.username;
    } else {
      // TODO find from route?
    }
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
          },
          (err) => {
            this.sharedService.showError('Can\'t remove miner key ', err.split(': ')[1]);
            console.log(err);
          }
        );
      }
    });
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
    this.exampleDatabase = new UserMinersForAdmin(this.httpClient, this.sharedService);
    this.dataSource = new ExampleDataSource(this.exampleDatabase, this.paginator, this.sort, this.username);
    if (this.filter) {
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

  chooseActiveState(value) {
    this.dataSource.activeState = value;
    this.dataSource._paginator.firstPage();
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

  constructor(public _exampleDatabase: UserMinersForAdmin,
              public _paginator: MatPaginator,
              public _sort: MatSort,
              public _username: string) {
    super();
    // Reset to the first page when the user changes the filter.
    this._filterChange.subscribe(() => this._paginator.pageIndex = 0);
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

    this._exampleDatabase.getMiners(this._username);
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
        const searchStr = (miner.id.toString()).toLowerCase();
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
