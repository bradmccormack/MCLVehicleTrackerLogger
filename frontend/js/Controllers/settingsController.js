angular.module('myApp.controllers').controller("settingsController", ['$scope', 'shellService', function($scope, shellService){

    var Helper = (function(){
        return {
            AccessLeveltoWord: function(Level) {
                Levels = {1: "Guest", 5: "User", 8: "Senior", 10: "Admin"};
                if(Level in Levels)
                    return Levels[Level];
                else
                    return "Unknown Access Level";
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
   
    $scope.User = shellService.User;
    $scope.User.Access = Helper.AccessLeveltoWord($scope.User.Access);

    $scope.Map = shellService.Map;
    $scope.Network = shellService.Network;
    $scope.Security = shellService.Security;
    $scope.Mobile = shellService.Mobile;
}]);