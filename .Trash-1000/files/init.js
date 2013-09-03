'use strict';


// Declare app level module which depends on filters, and services
var myApp = angular.module('myApp', ['myApp.filters', 'myApp.services', 'myApp.directives']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/', {templateUrl: 'index.html', controller: 'headerController'})
    $routeProvider.when('/#tabMap', {templateUrl: 'partials/test.html', controller: 'headerController'});
    $routeProvider.when('/#tabSupport', {templateUrl: 'partials/test2.html', controller: 'headerController'});
    $routeProvider.when('/#tabSettings', {templateUrl: 'partials/support.html', controller: 'headerController'});
    $routeProvider.otherwise({redirectTo: '/'});
  }]);




