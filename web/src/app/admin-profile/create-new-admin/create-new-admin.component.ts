import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-create-new-admin',
  templateUrl: './create-new-admin.component.html',
  styleUrls: ['./create-new-admin.component.scss']
})
export class CreateNewAdminComponent implements OnInit {
  createUser: FormGroup;
  permissions = new FormControl();
  permissionsList: string[] = ['Add Another Admin', 'Activate', 'Deactivate', 'Review DIY'];

  constructor() { }

  ngOnInit() {
    this.createUser = new FormGroup({
      email: new FormControl('', {validators: [Validators.required, Validators.email]}),
    });
  }

  onSubmit() {

  }

}
