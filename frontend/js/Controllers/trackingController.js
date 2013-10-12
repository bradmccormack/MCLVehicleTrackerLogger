angular.module('myApp.controllers').controller("trackingController", ['$scope', '$http', 'shellService', 'mapService', function($scope, $http, shellService, mapService){

   function formatDateSQL(Dte) {
	   var Months = Dte.getMonth() + 1 < 10 ? "0" + (Dte.getMonth() + 1) : Dte.getMonth() + 1;
	   var Days = Dte.getDate() < 10 ? "0" + Dte.getDate() : Dte.getDate();
	   var Hours = Dte.getHours() < 10 ? "0" + Dte.getHours() : Dte.getHours();
	   var Minutes = Dte.getMinutes() < 10 ? "0" + Dte.getMinutes() : Dte.getMinutes();
	   return [Dte.getFullYear(), Months, Days].join("-") + " " + [Hours, Minutes, Dte.getSeconds()].join(":");
   }

   function formatDate(Dte) {
	   var Hours = Dte.getHours() < 10 ? "0" + Dte.getHours() : Dte.getHours();
	   var Minutes = Dte.getMinutes() < 10 ? "0" + Dte.getMinutes() : Dte.getMinutes();
	   var Seconds = Dte.getSeconds() < 10 ? "0" + Dte.getSeconds() : Dte.getSeconds();
	   var Months = Dte.getMonth() + 1 < 10 ? "0" + (Dte.getMonth() + 1) : Dte.getMonth() + 1;
	   var Days = Dte.getDate() < 10 ? "0" + Dte.getDate() : Dte.getDate();
	   return [Days, Months, Dte.getFullYear()].join('/') + " " + [Hours, Minutes, Seconds].join(":");
   }

   function Init() {
       updateLiveInformation();

	   var Dte = new Date();
	   var DteFrom = new Date(Dte.getFullYear(), Dte.getMonth() - 1, Dte.getDay(), Dte.getHours(), Dte.getMinutes(), Dte.getSeconds());
	   $scope.routeDateFrom = formatDate(DteFrom);
	   $scope.routeDateTo = formatDate(Dte);

	   var datepickerFrom = $('#routeDateFrom');
	   datepickerFrom.datetimepicker({
		   language : 'en-AU',
		   pick12HourFormat : true,
		   format : 'dd/MM/yyyy hh:mm:ss',
		   startDate: $scope.routeDateFrom
	   });

	   var datepickerTo = $('#routeDateTo');
	   datepickerTo.datetimepicker({
		   language : 'en-AU',
		   pick12HourFormat : true,
		   format : 'dd/MM/yyyy hh:mm:ss',
		   startDate: $scope.routeDateTo
	   });

	   datepickerFrom.on('changeDate', function(e) {
		   $scope.routeDateFrom = formatDate(e.date);
	   });

	   datepickerTo.on('changeDate', function(e) {
		   $scope.routeDateTo = formatDate(e.date);
	   });

	   $('#routeDateFrom').data('datetimepicker').setLocalDate(DteFrom);
	   $('#routeDateTo').data('datetimepicker').setLocalDate(Dte);
   }

    var updateLiveInformation = function() {
        $scope.VehiclesCount = mapService.GetVehicleCount();
	    var LastTime = mapService.GetLastPosition().Time;
        var Delta = Math.round(Math.abs((new Date() - LastTime)) / 1000);
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

    $scope.Print = function() {
        //var elm = $("#MapCanvas");
        var elm = document.getElementById('MapCanvas');
        var width = elm.Width;
        var height = elm.Height;
        var dataUrl = "TODO.. figure out how to get the canvas ID of the embedded map. This will be tricky as it depends on what map API vendor is being used" // OR JUST grab the pixels by geting the pixel cords of the map div

        //var dataUrl = document.getElementById('MapCanvas').toDataURL(); //attempt to save base64 string to server using this var
        var windowContent = '<!DOCTYPE html>';
        windowContent += '<html>'
        windowContent += '<head><title>Print canvas</title></head>';
        windowContent += '<body>'
        windowContent += '<h1>TODO</h1>'
        windowContent += '<img src="' + dataUrl + '">';
        windowContent += '</body>';
        windowContent += '</html>';
        var printWin = window.open('','','width=' + width + ',height=' + height);
        printWin.document.open();
        printWin.document.write(windowContent);
        printWin.document.close();
        printWin.focus();
        printWin.print();
        printWin.close();
    }

	$scope.ShowRoute = function() {

		/*
		SELECT BusID, Latitude, Longitude, Speed, Heading, Fix, DateTime
		FROM GPSRecords
		WHERE datetime >='2013-07-2 22:14:45' AND datetime <='2013-10-12 12:14:45' GROUP BY id ORDER BY datetime asc
		*/

		var FromDate = formatDateSQL( $('#routeDateFrom').data('datetimepicker').getDate());
		var ToDate = formatDateSQL( $('#routeDateTo').data('datetimepicker').getDate());




		$http({method: 'POST', url: '/system/historicalroute', headers: {'Content-Type': 'application/x-www-form-urlencoded'},
			withCredentials: true, data: $.param({dateFrom: FromDate, dateTo: ToDate})}).
			success(function (result, status, headers, config) {
				if (result.success) {
					var vehicles = result.data;
					var vl = Object.keys(vehicles).length;
					if(vl == 0) {
						alert("No vehicle data for that time period");
						return;
					}

					while(true) {
						for(var i = 0; i < vl; i++) {
							var currentvehicle = Object.keys(vehicles)[i];
							var currentpositions = vehicles[currentvehicle];

							if(currentpositions.length > 0) {
								var point = vehicles[currentvehicle].shift();
								//Lat, Long, Speed, Fix, Heading, Date
								mapService.Map.AddtoRoute(currentvehicle,
									{Latitude: point[0], Longitude: point[1], Speed: point[2], Fix: point[3], Heading: point[4], DateTime: point[5]});
							} else {
								delete vehicles[currentvehicle];
								break;
							}
						}
						var vl = Object.keys(vehicles).length;
						if(vl == 0) {
							break;
						}

					}
				}

			}).
			error(function (data, status, headers, config) {
				alert(data + " - " + status + " - " + headers);
				//alert("Time to show error message as Couldn't get the route")
			});
	}

	$scope.SystemMessages = function(){
		return shellService.Messages;
	}

	$scope.$on('systemMessage', function(Event, Data){
		shellService.Messages.push({Message: Data, Read: false})
	});

    //We want to watch for changes on the model that the service will initiate via $broadcast
    $scope.$on('positionChange', function(Event, Data){
        //1)Update LastPosition received information
		mapService.UpdateLastPosition({Latitude: Data.Latitude, Longitude: Data.Longitude});
	    // 2) Update Marker and if there is no marker already set then add to the legend , generate colour etc
	    mapService.Map.SetMarker(Data.ID, Data.Latitude, Data.Longitude);


         //5) Draw line if Draw line functionality is set
    });


   Init();
   
   
   
}]);
