angular.module('myApp.controllers').controller("footerController", ['$scope', 'shellService', 'mapService', function($scope, shellService, mapService){

	$scope.Math = window.Math; //Inject Math object in so we can use round in binding eval
	$scope.Date = window.Date;
	$scope.VehicleCount = 0;
	$scope.VehicleData = {};

	$scope.$on('LegendChange', function(Event, Data){
		$scope.VehicleCount = Data.Count;
		$scope.Vehicles = Data.Vehicles;
	});

	$scope.$on('positionChange', function(Event, Data){

		$scope.VehicleData[Data.ID] = {
			Data: Data
		}
	});

    $scope.$on("ConfigChanged", function (Event, Data) {
        $scope.IsLogged = Data.User != "" && $.cookie("data") != undefined;
    });




}]);
