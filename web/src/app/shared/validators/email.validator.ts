import { FormGroup } from '@angular/forms';

export class EmailValidator {
   static validate(updateAddress: FormGroup) {
       let newEmail = updateAddress.controls.newEmail.value;
       let confirmEmail = updateAddress.controls.confirmEmail.value;

       if (confirmEmail.length <= 0) {
           return null;
       }

       if (confirmEmail !== newEmail) {
           return {
               doesMatchPassword: true
           };
       }

       return null;

   }
}
