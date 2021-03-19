# Web

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 6.1.3.

## Basic setup
Just run npm i, and run in one of the selected profiles(dev, production, mock)

## Development server and mock data

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

In case of running with mock data, you also need to run: json-server --watch db.json --routes routes.json in separate terminal to run the mock database.
To be able to run json-server you need to either install it with `npm i -g json-server` or use it from local node_modules using `node_modules/json-server/bin/index.js`
To run the server with mock data run ng s -c = mock.
db.json file contains mock data for return values, and routes contain the mapping for mock URL's return data.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.
Because of multiple modules inside the project, use ng g c new-component --module app (or name of other module where you want to export the component from).

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use the `--prod` flag for a production build.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).
