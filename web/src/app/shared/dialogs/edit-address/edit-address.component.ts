import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material';

@Component({
  selector: 'app-edit-address',
  templateUrl: './edit-address.component.html',
  styleUrls: ['./edit-address.component.scss']
})
export class EditAddressComponent implements OnInit {
  form: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private dialogRef: MatDialogRef<EditAddressComponent>,
  ) { }

  ngOnInit() {
    this.form = this.formBuilder.group({
      'address': ['', Validators.required],
    });
  }

  save() {
    this.dialogRef.close(this.form.get('address').value);
  }

  close() {
    this.dialogRef.close();
  }
}
