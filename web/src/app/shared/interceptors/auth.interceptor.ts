import {HttpErrorResponse, HttpHandler, HttpInterceptor, HttpRequest, HttpUserEvent} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/operator/do';
import {Injectable, Injector} from '@angular/core';
import {Router} from '@angular/router';
import {AuthService} from '../../services/auth.service';

const TOKEN_HEDER_KEY = 'Authorization';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

       constructor(private injector: Injector, private authService: AuthService, private router: Router) {}

     intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpUserEvent<any>> {
      //TODO check this
      //  if (req.headers.get('No-Auth') == "True")
      //    return next.handle(req.clone()).do(
      //      (err: any) => {
      //      }
      //    );

       let authReq = req;
      if (this.authService.getToken() != null) {
         authReq = req.clone({ headers: req.headers.set(TOKEN_HEDER_KEY, 'Bearer ' + this.authService.getToken())});
     }

      return next.handle(authReq).do(
        (err: any) => {
          if (err instanceof HttpErrorResponse) {
            if (err.status === 401) {
              this.router.navigate(['/']);
            }
          }
        }
      );

     }
}
