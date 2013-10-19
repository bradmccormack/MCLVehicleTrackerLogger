var LoginCtrl = angular.module('myApp.controllers').controller("loginController", ['$scope', '$cookieStore', '$http', 'shellService', 'authService', 'networkService', '$location', '$timeout',
	function ($scope, $cookieStore, $http, shellService, authService, networkService, $location, $timeout) {


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
            $scope.LoginProgress = true;
            $http({method: 'POST', url: '/system/login', headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                withCredentials: true, data: $.param({name: $scope.Username, password: $scope.Password})}).
                success(function (data, status, headers, config) {
                    if(data.success) {

                        $timeout(function(){
                            $scope.LoginProgress = false;
                            shellService.LoadConfig(data);
                            networkService.Init();
                            $location.path("/tracking");
                            authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc
                        }, 1000)

                    }
                    else {
                        $scope.Error = data.error;
                    }

                }).
                error(function (data, status, headers, config) {
                    var error = data;
                });
        }




	}]);


//the $q library is the promise library
LoginCtrl.Login = function($q, $http, $location, shellService, networkService) {

    var defer = $q.defer();

    //The following needs refining. If a malicious user was to inject a cookie called data it would log in. It would be erroneous, but it would still try.
    //Instead we still need to hit the server with the Post login regardless to confirm the data cookie is valid
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


                //authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc

            }).
            error(function (data, status, headers, config) {
                var error = data;
            });
        defer.resolve();

    }
    return defer.promise;
}
