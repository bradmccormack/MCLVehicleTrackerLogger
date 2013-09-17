'use strict';

// Declare app level module which depends on filters, and services
angular.module('myApp', ['myApp.filters', 'myApp.services', 'myApp.directives', 'myApp.controllers']).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {templateUrl: 'partials/login.html', controller: 'mainController'})
  $routeProvider.when('/settings', {templateUrl: 'partials/settings.html', controller: 'settingsController'})
  $routeProvider.when('/support', {templateUrl: 'partials/support.html', controller: 'supportController'})
  $routeProvider.when('/license', {templateUrl: 'partials/license.html', controller: 'licenseController'})
  $routeProvider.when('/reports', {templateUrl: 'partials/report.html', controller: 'reportController'})
  $routeProvider.when('/tracking', {templateUrl: 'partials/tracking.html', controller: 'trackingController'})
  $routeProvider.otherwise({redirectTo: '/tracking'});
}]);



  

