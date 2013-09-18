angular.module('myApp.controllers').controller("loginController", ['$scope', '$cookieStore', '$http', 'shellService', 'authService', function($scope, $cookieStore, $http, shellService, authService){

	$scope.Login = function() {
		$http({method: 'POST', url : '/system/login', headers: {'Content-Type': 'application/x-www-form-urlencoded'},
			withCredentials: true, data: $.param({name: $scope.Username, password: $scope.Password})}).
			success(function(data, status, headers, config){
				authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc

				shellService.User = {
					First: $cookieStore.get("User")["Firstname"],
					Last: $cookieStore.get("User")["Lastname"],
					Password: $cookieStore.get("User")["Password"],
					Access: $cookieStore.get("User")["AccessLevel"],
					Email: $cookieStore.get("User")["Email"]
				}

				//Start receiving Web-socket traffic
				networkService.Init();

			}).
			error(function(data, status, headers, config){

			});
	}




}]);
