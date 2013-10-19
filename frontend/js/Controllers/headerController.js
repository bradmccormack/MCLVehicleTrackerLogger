/*Pass in the shellService via Dependency injection rather than rely on prototypical inheritance between controllers to access User information */
/*Pass in http-auth-interceptor and ngCookies modules as module dependencies too */

angular.module('myApp.controllers', ['http-auth-interceptor', 'ngCookies']).controller("headerController", ['$scope', 'shellService', '$location', '$http',
	function ($scope, shellService, $location, $http) {

		$scope.clock = {
			interval: 1000,
			time: ""
		}

		$scope.User = {
			First: shellService.User.First,
			Last: shellService.User.Last
		}


		$scope.Logout = function(){
            $.removeCookie("data");
            $http({method: 'POST', url: '/system/logout',
				withCredentials: true}).
				success(function (data, status, headers, config) {

					//TODO Need to stop the network service networkService.Init();
                    $scope.User = "";

					$location.path("/login");

				}).
				error(function (data, status, headers, config) {

				});
	}

	$scope.$on("ConfigChanged", function (Event, Data) {
		$scope.User.First = Data.User.First;
		$scope.User.Last = Data.User.Last;
        updateClock();
	});


var edit = function () {

}

var updateClock = function () {
	$scope.clock.time = new Date().toLocaleString();

}


var timer = setInterval(function () {
	$scope.$apply(updateClock);
}, $scope.clock.interval);


}])
;
