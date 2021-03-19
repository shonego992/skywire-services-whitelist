import {DataSource} from '@angular/cdk/collections';
import {Observable} from 'rxjs/Observable';

export class CaseListDatasource extends DataSource<any> {


    constructor(private _list$: Observable<any[]>) {
        super();
    }

    connect(): Observable<any[]> {
        return this._list$;
    }

    disconnect() {
    }
}
