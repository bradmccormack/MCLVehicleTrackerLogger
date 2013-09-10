angular.module('myApp.controllers').controller("trackingController", ['$scope', 'shellService', 'mapService', function($scope, shellService, mapService){
   
   function Init() {
       updateLiveInformation();
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

    //We want to watch for changes on the model that the service will initiate via $broadcast
    $scope.$on('positionChange', function(Event, Data){


        //1)Update LastPosition received information
		mapService.UpdateLastPosition({Latitude: Data.Latitude, Longitude: Data.Longitude});
	    // 2) Update Marker and if there is no marker already set then add to the legend , generate colour etc
	    mapService.Map.SetMarker(Data.ID, Data.Latitude, Data.Longitude);

		/*

         4) Snap to the vehicle if it is selceted and the trigger count has been set
         5) Draw line if Draw line functionality is set
          */


    });


   Init();
   
   
   
}]);
