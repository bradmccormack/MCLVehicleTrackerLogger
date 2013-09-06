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

    $scope.Print = function() {
        //var elm = $("#MapCanvas");
        var elm = document.getElementById('MapCanvas');
        var width = elm.Width;
        var height = elm.Height;
        var dataUrl = "TODO.. figure out how to get the canvas ID of the embedded map. This will be tricky as it depends on what map API vendor is being used"

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
        alert("Received " + Data);

         /*
         1)Update LastPosition received information
         2) Update Legend if necessary
         3) Update Marker
         4) Snap to the vehicle if it is selceted and the trigger count has been set
         5) Draw line if Draw line functionality is set
          */
        mapService.LastPosition = {
            Time: new Date(),
            Position: Data.position
        }


    });


   Init();
   
   
   
}]);
