import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {HttpClient} from '@angular/common/http';
import {MatDialog, MatPaginator, MatSort} from '@angular/material';
import {Issue} from '../../models/issue';
import {Observable} from 'rxjs/Observable';
import {BehaviorSubject} from 'rxjs/BehaviorSubject';
import {DataSource} from '@angular/cdk/collections';
import 'rxjs/add/observable/merge';
import 'rxjs/add/observable/fromEvent';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import {AddDialogComponent} from '../../shared/dialogs/add/add.dialog.component';
import {DeleteDialogComponent} from '../../shared/dialogs/delete/delete.dialog.component';
import {SharedService} from '../../services/shared.service';
import {ImportedUsersService} from '../../services/imported-users.service';
import {environment} from '../../../environments/environment';
import {ApiRoutes} from '../../shared/routes';
import {HttpService} from '../../services/http.service';

@Component({
  selector: 'app-admin-miner-overview',
  templateUrl: './admin-imported-miners-overview.component.html',
  styleUrls: ['./admin-imported-miners-overview.component.scss']
})
export class AdminMinerOverviewComponent implements OnInit {
  displayedColumns = ['id', 'numberOfMiners', 'actions'];
  exampleDatabase: ImportedUsersService | null;
  dataSource: ExampleDataSource | null;
  id: number;
  userId: number;

  constructor(public httpClient: HttpClient,
              public dialog: MatDialog,
              public dataService: ImportedUsersService,
              public router: Router,
              public sharedService: SharedService,
              public activeRoute: ActivatedRoute,
              public httpService: HttpService) {}

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;
  @ViewChild('filter') filter: ElementRef;

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

  deleteItem(username: string): void {
    const dialogRef = this.dialog.open(DeleteDialogComponent, {
      data: {username: username, message: 'MINER.DELETE_MINER'}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result === 1) {
        const foundIndex = this.exampleDatabase.dataChange.value.findIndex(x => x.username === username);
        // for delete we use splice in order to remove single object from ImportedUsersService
        this.exampleDatabase.dataChange.value.splice(foundIndex, 1);
        this.refreshTable();
      }
    });
  }

  public setInputValue(row: any) {
    if(!row.numberOfMiners) {
      row.numberOfMiners = 0;
    }
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
    this.exampleDatabase = new ImportedUsersService(this.httpClient, this.sharedService);
    this.dataSource = new ExampleDataSource(this.exampleDatabase, this.paginator, this.sort);
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

  public saveChangesOnMiners() {
    const data = {data: this.dataSource.filteredData};
    this.httpService.postToUrl(environment.whitelistService + ApiRoutes.WHITELIST.Import, data).subscribe((res: any) => {
        this.sharedService.showSuccess('List uploaded to backend successfully');
      },
      (err: any) => {
        this.sharedService.showError('Can\'t upload list to server', err.split(': ')[1]);
      }
    );
  }


  public importChangesOnMiners() {
    const data = {data: this.dataSource.filteredData};
    this.httpService.postToUrl(environment.whitelistService + ApiRoutes.WHITELIST.ProcessImport, data)
      .subscribe((res: any) => {
          this.sharedService.showSuccess('Users imported to system sucesfully');
          window.location.reload();
        },
        (err: any) => {
          this.sharedService.showError('Can\'t import users to system', err.split(': ')[1]);
        }
      );
  }
}

export class ExampleDataSource extends DataSource<any> {
  _filterChange = new BehaviorSubject('');

  get filter(): string {
    return this._filterChange.value;
  }

  set filter(filter: string) {
    this._filterChange.next(filter);
  }

  filteredData: any[] = [];
  renderedData: any[] = [];

  constructor(public _exampleDatabase: ImportedUsersService,
              public _paginator: MatPaginator,
              public _sort: MatSort) {
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
      this._paginator.page
    ];

    this._exampleDatabase.getImportedUsers();;
    return Observable.merge(...displayDataChanges).map(() => {
      // Filter data
      this.filteredData = this._exampleDatabase.data.slice().filter((user: any) => {
        const searchStr = (user.username.toString()).toLowerCase();
        return searchStr.indexOf(this.filter.toLowerCase()) !== -1;
      });
      console.log(this.filteredData);

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
