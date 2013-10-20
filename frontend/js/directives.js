'use strict';

/* Directives */


angular.module('myApp.directives', []).
  directive('appVersion', ['version', function(version) {
    return function(scope, elm, attrs) {
      elm.text(version);
    };
  }]).directive('authapplication', function() {
		return {
			restrict: 'C',
			link: function(scope, elem, attrs) {

                scope.$on('event:server-lostContact', function() {
                    elem.addClass("waiting-for-angular");
                })

            }
		}
	});
