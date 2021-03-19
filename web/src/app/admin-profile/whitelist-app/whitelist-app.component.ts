import {ChangeDetectorRef, Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {DataService} from '../../services/data.service';
import {HttpClient} from '@angular/common/http';
import {MatDialog, MatPaginator, MatSort} from '@angular/material';
import {Issue, DropdownOption} from '../../models/issue';
import {Observable} from 'rxjs/Observable';
import {BehaviorSubject} from 'rxjs/BehaviorSubject';
import {DataSource} from '@angular/cdk/collections';
import 'rxjs/add/observable/merge';
import 'rxjs/add/observable/fromEvent';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import {AddDialogComponent} from '../../shared/dialogs/add/add.dialog.component';
import {WhitelistApplication} from '../../models/application-model';
import {SharedService} from '../../services/shared.service';


@Component({
  selector: 'app-whitelist-app',
  templateUrl: './whitelist-app.component.html',
  styleUrls: ['./whitelist-app.component.scss']
})
export class WhitelistAppComponent implements OnInit {
  startDate: Date = null;
  endDate: Date = null;
  displayedColumns = ['id', 'state', 'username', 'created_at', 'is_active', 'actions'];
  exampleDatabase: DataService | null;
  dataSource: ExampleDataSource | null;
  index: number;
  id: number;
  userId: string;
  filterValue: string = 'PENDING';
  statuses: DropdownOption[] = [
    { value: '', viewValue: 'All' },
    { value: 'PENDING', viewValue: 'Pending' },
    { value: 'APPROVED', viewValue: 'Approved' },
    { value: 'DECLINED', viewValue: 'Declined' },
    { value: 'DISABLED', viewValue: 'Disabled' }
  ];
  filterActiveState: string = '';
  activeStates: DropdownOption[] = [
    { value: '', viewValue: 'All' },
    { value: 'ACTIVE', viewValue: 'Active' },
    { value: 'DISABLED', viewValue: 'Disabled' }
  ];
    // this one exists but not used for now { value: 'CANCELED', viewValue: 'Canceled' }

  @Input()
  public username: string;

  constructor (public httpClient: HttpClient,
               public dialog: MatDialog,
               public dataService: DataService,
               public router: Router,
               public sharedService: SharedService,
               public activeRoute: ActivatedRoute,
               public changeDetector: ChangeDetectorRef) {}

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild('filter') filter: ElementRef;

  public filterDateChange() {
    this.dataSource.setDateFilter(this.startDate, this.endDate);
  }

  public getDateValue (value: string): string {
    const time = new Date(value);
    return time.toUTCString();
  }

  ngOnInit () {
    this.loadData();
    if (this.username) {
      this.dataSource.setUserFilter(this.username);
      this.filterValue = '';
    } else {
      this.filterValue = 'PENDING';
      this.dataSource.filter = 'PENDING';
    }
  }

  chooseActiveState(value) {
    this.dataSource.activeState = value;
    this.dataSource._paginator.firstPage();
  }

  addNew (issue: Issue) {
    const dialogRef = this.dialog.open(AddDialogComponent, {
      data: {issue: issue}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result === 1) {
        // After dialog is closed we're doing frontend updates
        // For add we're just pushing a new row inside DataService
        this.exampleDatabase.dataChange.value.push(this.dataService.getDialogData());
        this.refreshTable();
      }
    });
  }

  // If you don't need a filter or a pagination this can be simplified, you just use code from else block
  private refreshTable () {
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
      this.dataSource.filter = this.filterValue;
    }
  }

  chooseState(value) {
    this.dataSource.filter = value;
    this.dataSource._paginator.firstPage();
  }

  public loadData () {
    this.exampleDatabase = new DataService(this.httpClient, this.sharedService);
    this.dataSource = new ExampleDataSource(this.exampleDatabase, this.paginator, this.sort);
    this.changeDetector.detectChanges();
    }

    applyFilter(filterValue: string) {
      this.dataSource.filter = filterValue.trim().toLowerCase();
    }
}


export class ExampleDataSource extends DataSource<WhitelistApplication> {
  _filterChange = new BehaviorSubject('');
  _filterStatus = new BehaviorSubject('');
 private startDate: Date;
 private endDate: Date;
 private username: string;
 private id: number;

 public setDateFilter(startDate, endDate: Date) {
    this.startDate = startDate;
    this.endDate = endDate;
    this._filterChange.next(this.filter);
 }

 public setUserFilter(username: string) {
   this.username = username;
   this._filterChange.next(this.filter);
 }

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

 constructor(public _exampleDatabase: DataService,
             public _paginator: MatPaginator,
             public _sort: MatSort) {
   super();
   // Reset to the first page when the user changes the filter.
   this._filterChange.subscribe(() => this._paginator.pageIndex = 0);
 }

 /** Connect function called by the table to retrieve one stream containing the data to render. */
 connect(): Observable<WhitelistApplication[]> {
   // Listen for any changes in the base data, sorting, filtering, or pagination
   const displayDataChanges = [
     this._exampleDatabase.dataChange,
     this._sort.sortChange,
     this._filterChange,
     this._filterStatus,
     this._paginator.page
   ];

   this._exampleDatabase.getAllIssues();
   return Observable.merge(...displayDataChanges).map(() => {
     // Filter data
     this.filteredData = this._exampleDatabase.data.slice().filter((application: WhitelistApplication) => {
       const searchStr = (application.currentStatus.toString() + application.userId.toString() + application.id.toString()).toLowerCase();
       if (this.activeState && this.activeState.length > 0) {
         if (this.activeState === "ACTIVE" && application.disabled) {
           return false;
         } else if (this.activeState === "DISABLED" && !application.disabled) {
           return false;
         }
       }
       if (searchStr.indexOf(this.filter.toLowerCase()) !== -1 ) {
         const applicationDate = new Date(application.createdAt);
         if (this.startDate) {
           if (this.startDate > applicationDate) {
             return false;
           }
         }
         if (this.endDate) {
           if (this.endDate < applicationDate) {
             return false;
           }
         }
         if (this.username) {
          if (this.username != application.userId) {
            return false;
          }
         }
         return true;
       }
       return false;
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
      //  case 'username': [propertyA, propertyB] = [a.username, b.username]; break;
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
