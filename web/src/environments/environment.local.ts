// This file can be replaced during build by using the `fileReplacements` array.
// `ng build ---prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  userService: 'http://localhost:8080/api/v1',
  whitelistService: 'http://localhost:8081/api/v1',
  signUpURL: 'http://localhost:4200/signup',
  resetPasswordURL: 'http://localhost:4201/reset-password',
  imageBaseURL: 'https://s3.eu-central-1.amazonaws.com/cikaradule/',
  uptimeService: 'http://localhost:8085/api/v1'
};

/*
 * In development mode, to ignore zone related error stack frames such as
 * `zone.run`, `zoneDelegate.invokeTask` for easier debugging, you can
 * import the following file, but please comment it out in production mode
 * because it will have performance impact when throw error
 */
// import 'zone.js/dist/zone-error';  // Included with Angular CLI.
