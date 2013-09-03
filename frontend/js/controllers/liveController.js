angular.module('myClubLink.controllers', []).controller("liveController", ['$scope', function($scope){
   
   $scope.systemMessages = [];
   $scope.Live = {
        LastPosition : {
            TimeDelta : 2,
            Position : "123,456"
        },
        Vehicles : []
   }
   $scope.GetVehicleCount = function() {
        return $scope.Live.Vehicles.length;
   }
   
}]);