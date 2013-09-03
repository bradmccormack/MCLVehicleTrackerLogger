angular.module('myApp.controllers').controller("settingsController", ['$scope', function($scope){
   
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
   
   
    $scope.User =
    {
        First: "Brad",
        Last: "Mccormack",
        Access: "10"
    };
    
    $scope.Map =
    {
        Marker :
        {
            Smooth: true,
            SmoothCount: 10
        },
        Active: "Google Maps",

    }
    
    $scope.Network =
    {
        EnableRF: true,
        Enable3G: true
    }
    
    $scope.Security =
    {
        RemoteSupport: true,
        SystemConsoleAccess: true,
        AdminPasswordResetOnly: true
    }
    
    $scope.Mobile =
    {
        AllowSmartPhone: true,
        ShowSmartPhoneLocation: true
    }
}]);