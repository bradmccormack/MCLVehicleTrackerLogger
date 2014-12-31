angular.module('myApp.controllers').controller("settingsController", ['$scope', 'shellService', '$http', '$timeout', function($scope, shellService, $http, $timeout){

    var Loaded = false;

    var Helper = (function(){
        return {
            AccessLeveltoWord: function(Level) {
                if(Level == 0)
                    return "Guest";

                if(Level >= 1 && Level <= 5)
                    return "User";

                if(Level >=5 && Level <= 8)
                    return "Advanced User"
                if(Level >=8 && Level < 10)
                    return "Senior User";

				else if(Level == 10)
					return "Admin";

                else return "Unknown Access Level"

            }
        }
        
    })();

	$scope.User = shellService.User;
	$scope.User.AccessWord = Helper.AccessLeveltoWord($scope.User.Access);
	$scope.Map = shellService.Map;
	$scope.Settings = shellService.Settings;

    $scope.Password = 
    {
		ToggleResetHidden: function()
		{
			$scope.Password.Hidden = $scope.Settings.Security.AdminPasswordResetOnly ? $scope.User.Access !=10 : false;
		},
        ToggleResetWidget: function()
        {
			$scope.Password.WidgetHidden = !$scope.Password.WidgetHidden;
        },
        Update: function()
        {
            if ($scope.Password.New != $scope.Password.NewConfirm) {

			  $scope.Password.Error = "Confirmation password doesn't match new password";
              $timeout(function() { $scope.Password.Error = ""}, 2000);

			  $scope.Password.Old = $scope.Password.New = $scope.Password.NewConfirm = "";
              return;
            }
            else {
				$http({method: 'POST', url: '/system/settings/password', headers: {'Content-Type': 'application/json'},
					withCredentials: true, data:
						JSON.stringify({
							"passwordold" : encodeURIComponent($scope.Password.Old), //TODO encrypt to stop intercepting proxy
							"password" : encodeURIComponent($scope.Password.New) //TODO encrypt to stop intercepting proxy
						})}).
					success(function (data, status, headers, config) {
						if("success" in data) {
							$scope.Password.Note = "Password updated successfully";
							$timeout(function() { $scope.Password.Note = "";}, 2000)
						} else if("error" in data) {
							$scope.Password.Error = data.error;
							$timeout(function() { $scope.Password.Error = "";}, 2000)
						}

					}).
					error(function (data, status, headers, config) {
						$timeout(function() {$scope.Password.Error = "Error updating password";}, 2000);

					});
			}

            $scope.Password.Old = $scope.Password.New = $scope.Password.NewConfirm = "";
            $scope.Password.Hidden = !$scope.Password.Hidden;
        },
        Hidden: $scope.Settings.Security.AdminPasswordResetOnly ? $scope.User.Access != 10 : false,
		WidgetHidden: $scope.Settings.Security.AdminPasswordResetOnly ? $scope.User.Access != 10 : false,
        Old: "",
        New: "",
        Confirm: "",
        Error: false
    };

    $scope.$watch('Settings', function(e) {
        //Don't attempt to Persist if the Settings are changing due to view load.
        if(Loaded)
        {
            $http({method: 'POST', url: '/system/settings', headers: {'Content-Type': 'application/json'},
                withCredentials: true, data:
                JSON.stringify({
                    "MapAPI" : e.Map.API,
                    "Interpolate" : e.Map.Marker.Smooth,
                    "SnaptoRoad" : e.Map.Marker.SnaptoRoad,
                    "CameraPanTrigger" : e.Map.Marker.FollowVehicleTrigger,
                    "RadioCommunication" : e.Network.EnableRF,
                    "DataCommunication" : e.Network.Enable3G,
                    "SecurityRemoteAdmin" : e.Security.RemoteSupport,
                    "SecurityConsoleAccess" : e.Security.SystemConsoleAccess,
                    "SecurityAdminPasswordReset" : e.Security.AdminPasswordResetOnly,
                    "MobileSmartPhoneAccess" : e.Mobile.AllowSmartPhone,
                    "MobileShowBusLocation" : e.Mobile.ShowSmartPhoneLocation,
					"MinZoom"	: e.Map.Boundary.MinZoom,
					"MaxZoom"	: e.Map.Boundary.MaxZoom,
					"ClubBoundaryKM" : e.Map.Boundary.ClubBoundaryKM

                })}).
                success(function (data, status, headers, config) {
                    //authService.loginConfirmed(); //Login confirmed so the authservice will broadcast auth event which the directive will take care of and close login etc
                }).
                error(function (data, status, headers, config) {
                    var error = data;
                });

        }
        Loaded = true;

    }, true);



}]);