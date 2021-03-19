import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';

import {LoginComponent} from './auth/login/login.component';
import {AuthGuard} from './services/auth.guard';
import {RoleGuard} from './services/role.guard';
import {AdminProfileComponent} from './admin-profile/admin-profile.component';
import {UsersListComponent} from './admin-profile/users-list/users-list.component';
import {WhitelistFormComponent} from './whitelist-form/whitelist-form.component';
import {WhitelistAppComponent} from './admin-profile/whitelist-app/whitelist-app.component';
import {AdminEditWhitelistComponent} from './shared/dialogs/edit/edit.dialog.component';
import {EditAdminComponent} from './shared/dialogs/edit-admin/edit-admin.component';
import {EditUserComponent} from './shared/dialogs/edit-user/edit-user.component';
import {UserMinerOverviewComponent} from './user-miner-overview/user-miner-overview.component';
import {UploadUserListComponent} from './admin-profile/upload-user-list/upload-user-list.component';
import {UserMinerDetailsComponent} from './user-miner-details/user-miner-details.component';
import {LoginGuard} from './services/login.guard';
import {PageNotFoundComponent} from './page-not-found/page-not-found.component';
import {AdminAllMinersViewComponent} from './admin-profile/admin-all-miners-view/admin-all-miners-view.component';
import {AccountInfoComponent} from './account-info/account-info.component';
import {AdminMinerDetailsComponent} from './admin-miner-details/admin-miner-details.component';
import { LayoutComponent } from './layout/layout.component';

const routes: Routes = [
  {path: '', component: LoginComponent, canActivate: [LoginGuard]},
  {
    path: '',
    component: LayoutComponent,
    children: [
      {path: 'whitelist-form', component: WhitelistFormComponent, canActivate: [AuthGuard]},
      {path: 'user-miners', component: UserMinerOverviewComponent, canActivate: [AuthGuard]},
      {path: 'miner-details', component: UserMinerDetailsComponent, canActivate: [AuthGuard]},
      {path: 'users-list', component: UsersListComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'whitelist-app', component: WhitelistAppComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'edit', component: AdminEditWhitelistComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'edit-admin', component: EditAdminComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'edit-user', component: EditUserComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'admin', component: AdminProfileComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'import-miners', component: UploadUserListComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'view-miners', component: AdminAllMinersViewComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'admin-view-miner', component: AdminMinerDetailsComponent, canActivate: [AuthGuard], data: {expectedRole: 'is_admin'}},
      {path: 'account-info', component: AccountInfoComponent, canActivate: [AuthGuard]},
    ]
  },
  {path: '**', component: PageNotFoundComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
  providers: [AuthGuard, RoleGuard, LoginGuard]
})
export class AppRoutingModule {
}
