'use strict';


// Declare app level module which depends on filters, and services
angular.module('myClubLink', ['myClubLink.filters', 'myClubLink.services', 'myClubLink.directives', 'myClubLink.controllers']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/', {templateUrl: 'index.html', controller: 'headerController'})
    $routeProvider.when('/#tabMap', {templateUrl: '/system/map', controller: 'mapController'});
    $routeProvider.when('/#tabSupport', {templateUrl: '/system/support', controller: 'supportController'});
    $routeProvider.otherwise({redirectTo: '/'});
  }]);

/*
$(document).ready(function(){
    $("div#Mainnav ul li").mouseover(function(){
        $(this).addClass("bounceIn");
    });
})
*/


