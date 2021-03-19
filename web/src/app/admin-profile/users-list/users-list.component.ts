import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {Router} from '@angular/router';
import {HttpClient} from '@angular/common/http';
import {MatDialog, MatPaginator, MatSort} from '@angular/material';
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
import {User} from '../../models/user.model';
import {UserDataService} from '../../services/userData.service';
import {AdminClaims} from '../../models/admin.claims';
import {AuthService} from '../../services/auth.service';
import { DropdownOption } from 'src/app/models/issue';
import {HttpService} from '../../services/http.service';
import {ApiRoutes} from '../../shared/routes';
import {environment} from '../../../environments/environment';
import {SharedService} from '../../services/shared.service';

@Component({
  selector: 'app-users-list',
  templateUrl: './users-list.component.html',
  styleUrls: ['./users-list.component.scss']
})

export class UsersListComponent implements OnInit {
    displayedColumns = ['username', 'createdAt', 'is_active', 'own_miner', 'actions']; //TODO add maybe 'access rights' or 'last_sign_in'
    dataSource: ExampleDataSource | null;
    index: number;
    id: number;
    adminClaims: AdminClaims;
    filterActiveState: string = '';
    activeStates: DropdownOption[] = [
        { value: '', viewValue: 'All' },
        { value: 'ACTIVE', viewValue: 'Active' },
        { value: 'DISABLED', viewValue: 'Disabled' }
    ];

    constructor(public httpClient: HttpClient,
        public httpService: HttpService,
        public dialog: MatDialog,
        public dataService: UserDataService,
        public router: Router,
        public sharedService: SharedService,
        public authService: AuthService) { }

    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;
    @ViewChild('filter') filter: ElementRef;

    ngOnInit() {
        this.loadData(this);
        this.adminClaims = this.authService.getAdminClaims() || new AdminClaims();
    }

    public ownsOfficialMiner(row: any) {
      if (row.miners) {
        for (const miner of row.miners) {
          if (miner.type == 0) {
            return true;
          }
        }
      }
      return false;
    }

    refresh() {
        this.loadData(this);
    }

    addNew(issue: User) {
        const dialogRef = this.dialog.open(AddDialogComponent, {
            data: { issue: issue }
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result === 1) {
                // After dialog is closed we're doing frontend updates
                // For add we're just pushing a new row inside DataService
                this.dataService.dataChange.value.push(this.dataService.getDialogData());
                this.refreshTable();
            }
        });
    }

    startEdit(i: number, state: string, url: string) {
        // index row is used just for debugging proposes and can be removed
        this.index = i;
        this.router.navigate(['edit-user']);
    }

    chooseActiveState(value) {
        this.dataSource.activeState = value;
        this.dataSource._paginator.firstPage();
    }

    deleteUser(row: any) {
        const dialogRef = this.dialog.open(DeleteDialogComponent, {
            data: { message: 'MODALS.DISABLE_USER'}
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result === 1) {
              this.httpService.getFromUrl(
                environment.whitelistService + ApiRoutes.ADMIN.DisableUser + '?username=' + row.username).subscribe(
                (data: any) => {
                  row.status = 2;
                },
                (err: any) => {
                  this.sharedService.showError('Can\'t disable user ', err.split(': ')[1]);
                }
              );
            }
        });
    }

    activateUser(row: any) {
      const dialogRef = this.dialog.open(DeleteDialogComponent, {
        data: { message: 'MODALS.ENABLE_USER'}
      });

      dialogRef.afterClosed().subscribe(result => {
        if (result === 1) {
          this.httpService.getFromUrl(
            environment.whitelistService + ApiRoutes.ADMIN.EnableUser + '?username=' + row.username).subscribe(
            (data: any) => {
              row.status = 0;
            },
            (err: any) => {
              this.sharedService.showError('Can\'t enable user ', err.split(': ')[1]);
            }
          );
        }
      });
    }

    public getDateValue(value: string): string {
        const time = new Date(value);
        return time.toUTCString();
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

    public loadData(_this) {
        _this.dataService = new UserDataService(_this.httpClient);
        _this.dataSource = new ExampleDataSource(_this.dataService, _this.paginator, _this.sort);
        Observable.fromEvent(_this.filter.nativeElement, 'keyup')
            .debounceTime(150)
            .distinctUntilChanged()
            .subscribe(() => {
                if (!_this.dataSource) {
                    return;
                }
                let stringValue = _this.filter.nativeElement.value.toLowerCase();
                if (stringValue.includes('@gmail.com')) {
                  var parts = stringValue.split('@');
                  parts[0] = parts[0].split('.').join('');
                  stringValue = parts[0] + '@' +  parts[1];
                }
                _this.dataSource.filter = stringValue;
            });
    }
}


export class ExampleDataSource extends DataSource<User> {
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

    filteredData: User[] = [];
    renderedData: User[] = [];

    constructor(public _exampleDatabase: UserDataService,
        public _paginator: MatPaginator,
        public _sort: MatSort) {
        super();
        // Reset to the first page when the user changes the filter.
        this._filterChange.subscribe(() => this._paginator.pageIndex = 0);
    }

    /** Connect function called by the table to retrieve one stream containing the data to render. */
    connect(): Observable<User[]> {
        // Listen for any changes in the base data, sorting, filtering, or pagination
        const displayDataChanges = [
            this._exampleDatabase.dataChange,
            this._sort.sortChange,
            this._filterChange,
            this._filterStatus,
            this._paginator.page
        ];

        this._exampleDatabase.getAllUsers();

        return Observable.merge(...displayDataChanges).map(() => {
            // Filter data
            this.filteredData = this._exampleDatabase.data.slice().filter((user: User) => {
                const searchStr = (user.id + user.username).toLowerCase();
                if (this.activeState && this.activeState.length > 0) {
                    if (this.activeState === "ACTIVE" && user.status === 2) {
                        return false;
                    } else if (this.activeState === "DISABLED" && user.status !== 2) {
                        return false;
                    }
                }
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
    sortData(data: User[]): User[] {
        if (!this._sort.active || this._sort.direction === '') {
            return data;
        }

        return data.sort((a, b) => {
            let propertyA: number | string = '';
            let propertyB: number | string = '';

            switch (this._sort.active) {
                case 'id': [propertyA, propertyB] = [a.id, b.id]; break;
                case 'username': [propertyA, propertyB] = [a.username, b.username]; break;
                case 'createdAt': [propertyA, propertyB] = [a.createdAt, b.createdAt]; break;
                case 'status': [propertyA, propertyB] = [a.status, b.status]; break;
                case 'skycoinAddress': [propertyA, propertyB] = [a.skycoinAddress, b.skycoinAddress]; break;
            }

            const valueA = isNaN(+propertyA) ? propertyA : +propertyA;
            const valueB = isNaN(+propertyB) ? propertyB : +propertyB;

            return (valueA < valueB ? -1 : 1) * (this._sort.direction === 'asc' ? 1 : -1);
        });
    }
}
