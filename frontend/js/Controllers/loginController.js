var LoginCtrl = angular.module('myApp.controllers').controller("loginController", ['$scope', '$cookieStore', '$http', 'shellService', 'authService', 'networkService', '$location',
	function ($scope, $cookieStore, $http, shellService, authService, networkService, $location) {


        /*bind event handlers using jquery to do the animation */
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



        $scope.Login = function() {
            $http({method: 'POST', url: '/system/login', headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                withCredentials: true, data: $.param({name: $scope.Username, password: $scope.Password})}).
                success(function (data, status, headers, config) {
                    shellService.LoadConfig(data);
                    networkService.Init();
                    $location.path("/tracking");
                    authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc

                }).
                error(function (data, status, headers, config) {

                });
        }




	}]);


//the $q library is the promise library
LoginCtrl.Login = function($q, $http, $location, shellService, networkService) {

    var defer = $q.defer();
    if ("data" in $.cookie()) {

        defer.reject("User logged in already");

        $http({method: "GET", url: "/system/login", withCredentials: true}).
            success(function (data, status, headers, config) {
                shellService.LoadConfig(data);
                networkService.Init();
                $location.path("/tracking");
            })
    }
    else
    {
        $http({method: 'POST', url: '/system/login', headers: {'Content-Type': 'application/x-www-form-urlencoded'},
            withCredentials: true, data: $.param({name: '', password: ''})}).
            success(function (data, status, headers, config) {
                shellService.LoadConfig(data);
                networkService.Init();
                $location.path("/tracking");
                authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc

            }).
            error(function (data, status, headers, config) {

            });
        defer.resolve();

    }
    return defer.promise;
}
