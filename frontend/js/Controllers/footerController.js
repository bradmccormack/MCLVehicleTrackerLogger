angular.module('myApp.controllers').controller("footerController", ['$scope', 'mapService', function($scope, mapService){

	var mapLoaded = false;
	$scope.VehicleData = {};

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


	$scope.$on("mapLoaded", function(Event, Data) {
		mapLoaded = true;
	});


    $scope.$on("ConfigChanged", function (Event, Data) {
        $scope.IsLogged = (Data.User != undefined && Data.User != "" && ("data" in $.cookie()))

		//bind panning event
		mapService.Map.PanMap(Data.Settings.Map.Marker.FollowVehicleTrigger);
    });

	//We want to watch for changes on the model that the service will initiate via $broadcast
	$scope.$on('positionChange', function(Event, Data){

		if(mapLoaded) {

			$scope.VehicleData[Data.ID] = {
				Data: Data
			}

			//1)Update LastPosition received information
			mapService.UpdateLastPosition({Latitude: Data.Latitude, Longitude: Data.Longitude});
			// 2) Update Marker and if there is no marker already set then add to the legend , generate colour etc
			mapService.Map.SetMarker(Data.ID, Data.Latitude, Data.Longitude);
		}

		//) Draw line if Draw line functionality is set
	});



	$scope.SelectVehicle = function(Vehicle) {
		mapService.SelectVehicle(Vehicle.Ref);

	}




}]);
