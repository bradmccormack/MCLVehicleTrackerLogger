angular.module('myApp.controllers').controller("settingsController", ['$scope', 'shellService', '$http', function($scope, shellService, $http){

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
                if(Level >=8 && Level <= 10)
                    return "Admin";

                else return "Unknown Access Level"

            }
        }
        
    })();

   
    $scope.Password = 
    {
        ToggleResetWidget: function()
        {
            $scope.Password.Hidden = !$scope.Password.Hidden;
        },
        Update: function()
        {
            if ($scope.Password.New != $scope.Password.NewConfirm) {
              $scope.Password.Error = "Confirmation password doesn't match new password";
              $scope.Password.Old = $scope.Password.New = $scope.Password.NewConfirm = "";
              return;
            }
            else
                $scope.Password.Error = false;
            //TODO Check the Current password against the one stored in the cookie
            $scope.Password.Old = $scope.Password.New = $scope.Password.NewConfirm = "";
            $scope.Password.Hidden = !$scope.Password.Hidden;
        },
        Hidden: true,
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
                    "MobileShowBusLocation" : e.Mobile.ShowSmartPhoneLocation
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



    $scope.User = shellService.User;
    $scope.User.AccessWord = Helper.AccessLeveltoWord($scope.User.Access);
    $scope.Map = shellService.Map;
    $scope.Settings = shellService.Settings;

}]);