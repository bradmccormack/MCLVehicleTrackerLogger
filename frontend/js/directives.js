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

				var login = elem.find('#login-holder');
				var main = elem.find('#content');

				login.hide();

				scope.$on('event:auth-loginRequired', function() {
					login.slideDown('slow', function() {
						main.hide();
					});
				});
				scope.$on('event:auth-loginConfirmed', function() {
					main.show();
					login.slideUp();
				});
			}
		}
	});
