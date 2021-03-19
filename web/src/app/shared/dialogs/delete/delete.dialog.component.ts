import {Component, Inject} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';

@Component({
  selector: 'app-delete',
  templateUrl: './delete.component.html',
  styleUrls: ['./delete.component.scss']
})
export class DeleteDialogComponent  {
  constructor(public dialogRef: MatDialogRef<DeleteDialogComponent>,  @Inject(MAT_DIALOG_DATA) public data: any) {
  }

    onNoClick(): void {
      this.dialogRef.close();
    }

    confirmDelete(): void {
       return;
    }
}
