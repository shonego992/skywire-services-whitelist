import { FormGroup } from '@angular/forms';

export class PasswordValidator {

   static validate(passwordFormGroup: FormGroup) {
       let newPassword = passwordFormGroup.controls.newPassword.value;
       let confirmPassword = passwordFormGroup.controls.confirmPassword.value;

       if (confirmPassword.length <= 0) {
           return null;
       }

       if (confirmPassword !== newPassword) {
           return {
               doesMatchPassword: true
           };
       }

       return null;

   }
}
