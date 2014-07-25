angular.module('myApp.controllers').controller("footerController", ['$scope', 'mapService', function($scope, mapService){

	$scope.Math = window.Math; //Inject Math object in so we can use round in binding eval
	$scope.Date = window.Date;
	$scope.VehicleCount = 0;

	$scope.DiagnosticData = {};

	$scope.$on('LegendChange', function(Event, Data) {
		$scope.VehicleCount = Data.Count;
		$scope.Vehicles = Data.Vehicles;
	});


	$scope.$on('diagnosticChange', function(Event, Data){
		$scope.DiagnosticData[Data.ID] = {
			Data: Data
		}
	})


    $scope.$on("ConfigChanged", function (Event, Data) {
        $scope.IsLogged = (Data.User != undefined && Data.User != "" && ("data" in $.cookie()))

		//bind panning event
		mapService.Map.PanMap(Data.Settings.Map.Marker.FollowVehicleTrigger);
    });


	$scope.SelectVehicle = function(Vehicle) {
		mapService.SelectVehicle(Vehicle.Ref);

	}




}]);
