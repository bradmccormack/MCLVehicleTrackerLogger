/*TODO add the http-auth-interceptor as a dependency and inject in the authService so we can confirm login if we need to */

angular.module('myApp.controllers').controller("mainController", ['$scope', '$http', 'networkService', 'shellService', 'authService' ,
	function ($scope, $http, networkService, shellService, authService) {


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


		//send over post to login as a urlencoded form
		$http({method: 'POST', url : '/system/login', headers: {'Content-Type': 'application/x-www-form-urlencoded'}, withCredentials: true, data: $.param({name: "", password: ""})}).
			success(function(data, status, headers, config){
				authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc
				networkService.Init();
				//check that the session cookie is set

			}).
			error(function(data, status, headers, config){

			});

		/*
		var cookies = $.cookie();
		if ("session" in $.cookie()) {

			//set the properties on the user object etc in the shellService

		}
		else
		{
			//$http({method: "POST"})
		}
		*/

		/*
		var myModal = $("#myModal");
		$("#myModal").toggle();
		*/

	}

	$scope.SystemInit = function () {
		BindAnim();
		Login();
	}


}]);
