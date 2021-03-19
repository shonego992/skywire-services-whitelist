import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {Component, Inject, OnInit} from '@angular/core';
import {DataService} from '../../../services/data.service';
import {FormControl, Validators} from '@angular/forms';
import {Issue} from '../../../models/issue';

@Component({
  selector: 'app-add.dialog',
  templateUrl: './add.dialog.component.html',
  styleUrls: ['./add.component.scss']
})
export class AddDialogComponent implements OnInit {
  constructor(public dialogRef: MatDialogRef<AddDialogComponent>,
                @Inject(MAT_DIALOG_DATA) public data: Issue,
                public dataService: DataService) { }

    formControl = new FormControl('', [
      Validators.required
      // Validators.email,
    ]);

    getErrorMessage() {
      return this.formControl.hasError('required') ? 'Required field' :
        this.formControl.hasError('email') ? 'Not a valid email' :
          '';
    }

    submit() {
    // emppty stuff
    }

    onNoClick(): void {
      this.dialogRef.close();
    }

    public confirmAdd(): void {
      this.dataService.addIssue(this.data);
    }

    ngOnInit() {}
}
