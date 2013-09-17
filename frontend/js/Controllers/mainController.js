angular.module('myApp.controllers').controller("mainController", ['$scope', 'networkService', 'shellService', '$http', function ($scope, networkService, shellService, $http) {


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

	function Login() {

		//send over post to login
		$http({method: 'POST', url : '/system/login', data: {name: "brad", password: "fred"}}).
			success(function(data, status, headers, config){

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
		networkService.Init();
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
