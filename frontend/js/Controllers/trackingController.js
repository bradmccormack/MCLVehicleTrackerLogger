angular.module('myApp.controllers').controller("trackingController", ['$scope', 'shellService', 'mapService', function($scope, shellService, mapService){
   
   function Init() {
       updateLiveInformation();
   }

    var updateLiveInformation = function() {
        $scope.VehiclesCount = mapService.Vehicles.length;
        var Delta = Math.round(Math.abs((new Date() - mapService.GetLastPosition().Time)) / 1000);
        var Minutes = Math.floor(Delta/ 60);
        var Seconds = Delta - Minutes * 60;
        $scope.LastUpdate = Minutes + " (Min) " + Seconds + " (Sec)";
    };

    var timer = setInterval(function() {
        $scope.$apply(updateLiveInformation);
    }, 1000);
   
   Init();
   
   
   
}]);
