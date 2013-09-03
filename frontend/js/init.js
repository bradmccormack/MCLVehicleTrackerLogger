'use strict';


// Declare app level module which depends on filters, and services
var myApp = angular.module('myClubLink', ['myClubLink.filters', 'myClubLink.services', 'myClubLink.directives', 'myClubLink.controllers']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/', {templateUrl: 'index.html', controller: 'headerController'})
    $routeProvider.when('/#tabMap', {templateUrl: 'partials/test.html', controller: 'mapController'});
    $routeProvider.when('/#tabSupport', {templateUrl: 'partials/test2.html', controller: 'supportController'});
    $routeProvider.otherwise({redirectTo: '/'});
  }]);




