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

    $scope.MapRefresh = function() {
        mapService.Map.Refresh();
    }

    $scope.MapZoomIn = function() {
        mapService.Map.ZoomIn();
    }

    $scope.MapZoomOut = function() {
        mapService.Map.ZoomOut();
    }

    $scope.MapReload = function() {
        mapService.Map.ReLoad();
    };

    //We want to watch for changes on the model that the service will initiate via $broadcast
    $scope.$on('positionChange', function(event, vehicle){
        alert(vehicle);
    });


   Init();
   
   
   
}]);
