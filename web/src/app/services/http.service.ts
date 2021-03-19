import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class HttpService {

    public constructor(private client: HttpClient) {
    }

    public postToUrl<T>(url: string, body: T, httpOptions?: any): Observable<any> {
        return this.client.post(url, body, httpOptions);
    }


  public postToUrlFormData<T>(url: string, data?: FormData, httpOptions?: any): Observable<any> {
    return this.client.post(url, data);
  }

    public patchToUrl<T>(url: string, body: T, httpOptions?: any): Observable<any> {
        return this.client.patch(url, body, httpOptions);
    }

    public getFromUrl(url: string): Observable<any> {
        return this.client.get<any>(url);
    }

    public deleteFromUrl(url: string, httpOptions?: any): Observable<any> {
        return this.client.delete<any>(url, httpOptions);
    }

}
