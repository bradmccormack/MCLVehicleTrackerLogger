angular.module('myApp.controllers').controller("loginController", ['$scope', '$cookieStore', '$http', 'shellService', 'authService', 'networkService', '$location',
	function ($scope, $cookieStore, $http, shellService, authService, networkService, $location) {

		$scope.Login = function () {

			if ("data" in $.cookie()) {
				$http({method: "GET", url: "/system/login", withCredentials: true}).
					success(function (data, status, headers, config) {
						//set the shellService values
						shellService.User = {
							First: data.user.Firstname,
							Last: data.user.Lastname,
							Password: data.user.Password, //TODO encrypt serverside
							Access: data.user.AccessLevel,
							Email: data.user.Email
						};

						shellService.Company = {
							Name: data.company.Name,
							MaxUsers: data.company.Maxusers,
							Expiry: data.company.Expiry,
							Logo: data.company.LogoPath
						}
						shellService.Settings = {
							Network: {
								EnableRF: data.settings.RadioCommunication,
								Enable3G: data.settings.DataCommunication
							},
							Security: {
								RemoteSupport: data.settings.SecurityRemoteAdmin,
								SystemConsoleAccess: data.settings.SecurityRemoteConsoleAccess,
								AdminPasswordResetOnly: data.settings.SecurityAdminPasswordReset
							},
							Mobile: {
								AllowSmartPhone: data.settings.MobileSmartPhoneAccess,
								ShowSmartPhoneLocation: data.settings.MobileShowBusLocation
							},
							Map: {
								API: data.settings.MapAPI,
								Marker: {
									Smooth: data.settings.Interpolate,
									SnaptoRoad: data.settings.SnaptoRoad,
									FollowVehicleTrigger: data.settings.CameraPanTrigger
								}
							}}
						$location.path("/tracking");
					})
			}
			else {
				$http({method: 'POST', url: '/system/login', headers: {'Content-Type': 'application/x-www-form-urlencoded'},
					withCredentials: true, data: $.param({name: $scope.Username, password: $scope.Password})}).
					success(function (data, status, headers, config) {

						/*
						 //Do I need to just grab the JSON and assign it to the cookiestore ? I think so.
						 //OR do I grab the values out of the session cookie
						 if("data" in $.cookie()) {
						 console.log($.cookie("session"));
						 var First = $.cookie("session")["User"];
						 }
						 shellService.User = {
						 First: $cookieStore.get("User")["Firstname"],
						 Last: $cookieStore.get("User")["Lastname"],
						 Password: $cookieStore.get("User")["Password"],
						 Access: $cookieStore.get("User")["AccessLevel"],
						 Email: $cookieStore.get("User")["Email"]
						 }
						 */
						//Start receiving Web-socket traffic
						networkService.Init();
						authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc


					}).
					error(function (data, status, headers, config) {

					});
			}
		}


	}]);
