import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FlexLayoutModule} from '@angular/flex-layout';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HttpClientModule, HTTP_INTERCEPTORS, HttpClient} from '@angular/common/http';
import {ToastrModule} from 'ngx-toastr';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import { NgxLoadingModule } from 'ngx-loading';

import {AppComponent} from './app.component';
import {MaterialModule} from './material.module';
import {LoginComponent} from './auth/login/login.component';
import {AppRoutingModule} from './app-routing.module';
import {HeaderComponent} from './navigation/header/header.component';
import {SidenavListComponent} from './navigation/sidenav-list/sidenav-list.component';
import {AuthService} from './services/auth.service';
import {UserService} from './services/user.service';
import {DataService} from './services/data.service';
import {UserDataService} from './services/userData.service';
import {AuthInterceptor} from './shared/interceptors/auth.interceptor';
import {AdminProfileComponent} from './admin-profile/admin-profile.component';
import {HttpService} from './services/http.service';
import {UsersListComponent} from './admin-profile/users-list/users-list.component';
import {CreateNewAdminComponent} from './admin-profile/create-new-admin/create-new-admin.component';
import {AddDialogComponent} from './shared/dialogs/add/add.dialog.component';
import {AdminEditWhitelistComponent} from './shared/dialogs/edit/edit.dialog.component';
import {DeleteDialogComponent} from './shared/dialogs/delete/delete.dialog.component';
import {TranslateLoader, TranslateModule} from '@ngx-translate/core';
import {ErrorInterceptor} from './shared/interceptors/error.interceptor';
import {TranslateHttpLoader} from '@ngx-translate/http-loader';
import {SharedService} from './services/shared.service';
import {WhitelistFormComponent} from './whitelist-form/whitelist-form.component';
import {NgxUploaderModule} from 'ngx-uploader';
import {WhitelistAppComponent} from './admin-profile/whitelist-app/whitelist-app.component';
import {EditAdminComponent} from './shared/dialogs/edit-admin/edit-admin.component';
import {EditUserComponent} from './shared/dialogs/edit-user/edit-user.component';

// TODO: consider to drop components below this line
import {AlertComponent} from './shared/dialogs/alert/alert.component';
import {NgxImageGalleryModule} from 'ngx-image-gallery';
import {UserMinerOverviewComponent} from './user-miner-overview/user-miner-overview.component';
import {UploadUserListComponent} from './admin-profile/upload-user-list/upload-user-list.component';
import {UploadService} from './services/uploader.service';
import {MinersService} from './services/miners.service';
import {UserMinerDetailsComponent} from './user-miner-details/user-miner-details.component';
import {AdminMinerOverviewComponent} from './admin-profile/admin-imported-miners-overview/admin-imported-miners-overview.component';
import {ImportedUsersService} from './services/imported-users.service';
import {PageNotFoundComponent} from './page-not-found/page-not-found.component';
import {UserMinersForAdminComponent} from './user-miners-for-admin/user-miners-for-admin.component';
import {UserMinersForAdmin} from './services/miners-service-for-admin.service';
import {AdminAllMinersViewComponent} from './admin-profile/admin-all-miners-view/admin-all-miners-view.component';
import {AdminMinersService} from './services/admin-miners.service';
import {AccountInfoComponent} from './account-info/account-info.component';
import {EditMinerComponent} from './shared/dialogs/edit-miner/edit-miner.component';
import {AdminMinerDetailsComponent} from './admin-miner-details/admin-miner-details.component';
import { UptimePipe } from './shared/validators/uptime.pipe';
import { LayoutComponent } from './layout/layout.component';
import { EditAddressComponent } from './shared/dialogs/edit-address/edit-address.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    HeaderComponent,
    SidenavListComponent,
    AdminProfileComponent,
    UsersListComponent,
    CreateNewAdminComponent,
    AddDialogComponent,
    AdminEditWhitelistComponent,
    DeleteDialogComponent,
    AlertComponent,
    WhitelistFormComponent,
    WhitelistAppComponent,
    EditAdminComponent,
    EditUserComponent,
    UserMinerOverviewComponent,
    UploadUserListComponent,
    UserMinerDetailsComponent,
    AdminMinerOverviewComponent,
    PageNotFoundComponent,
    UserMinersForAdminComponent,
    AdminAllMinersViewComponent,
    AccountInfoComponent,
    EditMinerComponent,
    AdminMinerDetailsComponent,
    UptimePipe,
    LayoutComponent,
    EditAddressComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    MaterialModule,
    AppRoutingModule,
    FlexLayoutModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    NgxUploaderModule,
    NgxImageGalleryModule,
    ToastrModule.forRoot({
      timeOut: 10000,
      positionClass: 'toast-top-center',
      preventDuplicates: true
    }),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useFactory: (createTranslateLoader),
        deps: [HttpClient]
      }
    }),
    NgxLoadingModule.forRoot({
    })
  ],
  entryComponents: [
    AddDialogComponent,
    AdminEditWhitelistComponent,
    DeleteDialogComponent,
    EditMinerComponent,
    EditAddressComponent,
  ],
  providers: [
    AuthService,
    UserService,
    HttpService,
    DataService,
    MinersService,
    UserMinersForAdmin,
    UserDataService,
    ImportedUsersService,
    SharedService,
    UploadService,
    AdminMinersService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: ErrorInterceptor,
      multi: true
    },
    {
      provide: MAT_DIALOG_DATA,
      useValue: {}
    },
    {
      provide: MatDialogRef,
      useValue: {}
    }

  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}

export function createTranslateLoader (http: HttpClient) {
  return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}
