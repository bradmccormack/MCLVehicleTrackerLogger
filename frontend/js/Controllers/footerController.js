angular.module('myApp.controllers').controller("footerController", ['$scope', 'shellService', 'mapService', function($scope, shellService, mapService){

	$scope.VehicleCount = 0;

	$scope.$on('LegendChange', function(Event, Data){
		$scope.VehicleCount = Data.Count;
		$scope.Vehicles = Data.Vehicles;
	});



}]);
