/*Pass in the shellService via Dependency injection rather than rely on prototypical inheritance between controllers to access User information */
/*Pass in http-auth-interceptor and ngCookies modules as module dependencies too */

angular.module('myApp.controllers', ['http-auth-interceptor', 'ngCookies']).controller("headerController", ['$scope', 'shellService', '$location', '$http', 'networkService', '$cookieStore',
	function ($scope, shellService, $location, $http, networkService, $cookieStore) {



		$scope.Logout = function () {
			shellService.ClearConfig();
			networkService.Stop();


			$http({method: 'POST', url: '/system/logout',
				withCredentials: true}).
				success(function (data, status, headers, config) {
					$cookieStore.remove("data");
					$location.path("/login");
				}).
				error(function (data, status, headers, config) {
					var error = data;
				});
		}

		$scope.$on("ConfigChanged", function (Event, Data) {
			$scope.IsLogged = (Data.User != undefined && Data.User != "" && ("data" in $.cookie()))
			if ($scope.IsLogged) {
				$scope.User = {
					First: Data.User.First,
					Last: Data.User.Last
				}
			}
			else $scope.User = {};

		});

		$scope.$on('event:server-lostContact', function () {
			shellService.ClearConfig();
		});


		var edit = function () {

		}


	}])
;
