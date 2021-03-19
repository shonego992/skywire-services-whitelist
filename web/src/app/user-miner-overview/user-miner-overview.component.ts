import {Component, ElementRef, Input, OnInit, ViewChild, TemplateRef} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
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
import {WhitelistApplication} from '../models/application-model';
import {SharedService} from '../services/shared.service';
import {MinersService} from '../services/miners.service';

@Component({
  selector: 'app-user-miner-overview',
  templateUrl: './user-miner-overview.component.html',
  styleUrls: ['./user-miner-overview.component.scss']
})
export class UserMinerOverviewComponent implements OnInit {
  displayedColumns = ['id', 'type', 'created_at', 'updated_at', 'label', 'is_active', 'actions'];
  exampleDatabase: MinersService | null;
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
  noMiners = false;

  @Input()
  public username: string;

  constructor(public httpClient: HttpClient,
              public dialog: MatDialog,
              public dataService: MinersService,
              public router: Router,
              public sharedService: SharedService,
              public activeRoute: ActivatedRoute) {}

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild('filter') filter: ElementRef;

  chooseActiveState(value) {
    this.dataSource.activeState = value;
    this.dataSource._paginator.firstPage();
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

  public mapType(value: number): string {
    return this.sharedService.mapMinerType(value);
  }

  ngOnInit() {
    this.loadData();

    if (this.username) {
      this.userId = this.username;
    } else {
      // TODO find from route?
    }
    if (this.userId) {
      if (!this.username) {
      this.filter.nativeElement.value = this.userId.toString();
      }
      this.dataSource.filter = this.userId.toString();
    } else {
      // this.filter.nativeElement.value = 'PENDING';
      // this.dataSource.filter = 'PENDING';
    }
  }

  public createApp() {
    this.router.navigateByUrl('/whitelist-form');
  }

  refresh() {
    this.loadData();
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
    this.exampleDatabase = new MinersService(this.httpClient, this.sharedService);
    this.dataSource = new ExampleDataSource(this.exampleDatabase, this.paginator, this.sort, this);
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
}


export class ExampleDataSource extends DataSource<WhitelistApplication> {
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

  constructor(public _exampleDatabase: MinersService,
              public _paginator: MatPaginator,
              public _sort: MatSort,
              public _component: UserMinerOverviewComponent) {
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

    this._exampleDatabase.getMiners();
    return Observable.merge(...displayDataChanges).map(() => {
      setTimeout(() => { this._component.noMiners = this._exampleDatabase.data.length === 0; });

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
