'use strict';


// Declare app level module which depends on filters, and services
angular.module('myApp', ['myApp.filters', 'myApp.services', 'myApp.directives', 'myApp.controllers']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/view1', {templateUrl: 'partials/partial1.html', controller: 'MyCtrl1'});
    $routeProvider.when('/view2', {templateUrl: 'partials/partial2.html', controller: 'MyCtrl2'});
    $routeProvider.when('/view3', {templateUrl: 'partials/partial3.html', controller: 'MyCtrl2'});
    $routeProvider.when('/settings', {templateUrl: 'partials/settings.html', controller: 'settingsController'})
    $routeProvider.when('/support', {templateUrl: 'partials/support.html', controller: 'supportController'})
    $routeProvider.when('/license', {templateUrl: 'partials/license.html', controller: 'licenseController'})
    $routeProvider.when('/reports', {templateUrl: 'partials/report.html', controller: 'reportController'})
    $routeProvider.otherwise({redirectTo: '/'});
  }]);
