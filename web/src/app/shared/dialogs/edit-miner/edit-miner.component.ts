import { Component, Inject } from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';

@Component({
  selector: 'app-edit-miner',
  templateUrl: './edit-miner.component.html',
  styleUrls: ['./edit-miner.component.scss']
})
export class EditMinerComponent {
  constructor(public dialogRef: MatDialogRef<EditMinerComponent>,  @Inject(MAT_DIALOG_DATA) public data: any) {
    console.log(this.data);
  }

    onNoClick(): void {
      this.dialogRef.close();
    }

    confirmDelete(): void {
       return;
    }
}
