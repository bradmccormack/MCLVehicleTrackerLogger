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

				//show waiting for server lost contact view if there was a connection drop with the server.
                scope.$on('event:server-lostContact', function() {
                    elem.addClass("waiting-for-angular");
					var timeout = setInterval(function(){
						window.location="/login";
					}, 5000);
				})

            }
		}
	});
