angular.module('myApp.controllers').controller("trackingController", ['$scope', '$http', 'shellService', 'mapService', function($scope, $http, shellService, mapService){


	$scope.clock = {
		interval: 1000,
		time: ""
	}

	var updateClock = function () {
		$scope.clock.time = moment().format("Do MMM YYYY, h:mm:ss a")
	}


	var mapLoaded = false;

   function Init() {
       updateLiveInformation();

	   var timer = setInterval(function () {
		   $scope.$apply(updateClock);
	   }, $scope.clock.interval);

       //Set the time in Local time to be today
       var Dte = moment().local();
       var DteFrom = moment().local().subtract('days', 1);


       $scope.routeDateFrom = DteFrom.format("DD/MM/YYYY HH:MM:SS");
       $scope.routeDateTo = Dte.format("DD/MM/YYYY HH:MM:SS");

       var datepickerFrom = $('#routeDateFrom');
       datepickerFrom.datetimepicker({
           language : 'en-AU',
           pick12HourFormat : true,
           format : 'dd/MM/yyyy hh:mm:ss'
       });

       var datepickerTo = $('#routeDateTo');
       datepickerTo.datetimepicker({
           language : 'en-AU',
           pick12HourFormat : true,
           format : 'dd/MM/yyyy hh:mm:ss'
       });


       $('#routeDateFrom').data('datetimepicker').setLocalDate(DteFrom.toDate());
       $('#routeDateTo').data('datetimepicker').setLocalDate(Dte.toDate());
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


    $scope.LiveMode = function(value) {
        mapService.Map.SetMode(value);
    }



    $scope.Print = function() {
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


        //Parse current values and convert to a format that we can use in the SQL query
        var FromDate = moment($('#routeDateFrom').data('datetimepicker').getLocalDate()).format("YYYY-MM-DD HH:MM:SS");
        var ToDate =   moment($('#routeDateTo').data('datetimepicker').getLocalDate()).format("YYYY-MM-DD HH:MM:SS");


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

					//this needs to be rewrittent to build up a path object then just splat it over the otherside... it won't ever "replay" like originally intended
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

	$scope.$on("mapLoaded", function(Event, Data) {
		mapLoaded = true;
	});

	//something has changed (perhaps the camera snap so rebind panning event
	$scope.$on("ConfigChanged", function (Event, Data) {
		if(mapLoaded) {
			//bind panning event
			mapService.Map.PanMap(Data.Settings.Map.Marker.FollowVehicleTrigger);
		}

	});


	$scope.SystemMessages = function(){
		return shellService.Messages;
	}

	$scope.SystemMessageCount = function() {
		return Object.keys(shellService.Messages).length;
	}

	$scope.MarkSystemMessageRead = function(ID) {
		delete shellService.Messages[ID];
	}

	$scope.$on('systemMessage', function(Event, Data){

		var Message = {
			Text: Data.message, Read: false, MsgDateTime: new Date().toLocaleTimeString()
		}
		if(Data.warning)
			Message.Warning = true;

		if(Data.information)
			Message.Information = true;

		shellService.Messages[Message.MsgDateTime] = Message;
	});



   Init();
   
   
   
}]);
