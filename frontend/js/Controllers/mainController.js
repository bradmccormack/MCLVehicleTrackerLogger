/*TODO add the http-auth-interceptor as a dependency and inject in the authService so we can confirm login if we need to */

angular.module('myApp.controllers').controller("mainController", ['$scope', 'shellService', '$location',
	function ($scope, shellService, $location) {


	function BindAnim()
	{
		//Login
		$('#myModal').on("shown", function()
		{
			$(".username").focus(function() {
				$(".user-icon").css("left","-48px");
			});
			$(".username").blur(function() {
				$(".user-icon").css("left","0px");
			});

			$(".password").focus(function() {
				$(".pass-icon").css("left","-48px");
			});
			$(".password").blur(function() {
				$(".pass-icon").css("left","0px");
			});
		});
	}

	//This login is only called at the start of the application with empty name and password fields so the security will fail. This will result in 401 which will trigger login box to appear and real login can take over
	function Login() {
		$location.path("/login");
	}

	$scope.SystemInit = function () {
		BindAnim();
		Login();
	}


}]);
