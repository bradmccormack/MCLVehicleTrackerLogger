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
				//once Angular is started, remove class:
				elem.removeClass('waiting-for-angular');

				scope.$on('event:auth-loginRequired', function() {
					//main.hide();
					//login.toggle();

				});
				scope.$on('event:auth-loginConfirmed', function() {
					//login.toggle();
					//main.show();

				});

                scope.$on('event:server-lostContact', function() {
                    elem.addClass("waiting-for-angular");
                })

            }
		}
	});
