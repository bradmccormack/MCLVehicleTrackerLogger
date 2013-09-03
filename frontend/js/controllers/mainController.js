myApp.controller("mainController", [function($scope) {
{
    $scope.Con = {};
    $scope.Colours = [];
    $scope.Vehicles = {};
    $scope.Camera = {};
    $scope.Map = {};
    $scope.Messages = [];
    
    var Utility = (function(){
         return {
                 RandomColor: function() {
                         return (function lol(m,s,c){return s[m.floor(m.random() * s.length)] +
                         (c && lol(m,s,c-1));})(Math,'0123456789ABCDEF',4)
                 }
         }	
     })();
   
    $scope.showLostConnection = function()
    {
	$('#systemError').modal('toggle');  
    }
    
    $scope.systemMessage = function(Message)
    {
	$("div#SystemMessages > ul#Messages").append("<li class='text-info'>" + new Date().toTimeString() + Message + "</li>");
    }
    
    $scope.init = function()
    {
	if (window["WebSocket"])
	{
		//Con = new WebSocket("ws://dev.myclublink.com.au/ws");
		$scope.Con = new WebSocket("ws://dev.myclublink.com.au:8080/ws");
	     
		$scope.Con.onopen = function()
		{
		    $scope.systemMessage("Connected to server");
		};

		$scope.Con.onclose = function(evt)
		{
		    $scope.systemMessage("Server connection closed");
		}
		$scope.Con.onmessage = function(evt)
		{
		    var data = JSON.parse(evt.data).Entry;
		  
                    //add vehicle to Legend if it is not there   
                    if(!(data.ID in Vehicles))
                    {
                        mapController.updateLegend(data.ID, Utility.RandomColor());
                    }
                    //TODO remove vehicle if no contact for X minutes
                    
                    $scope.Map.setMarker(data.ID, data.Latitude, data.Longitude,"", Vehicles[data.ID].Color,Settings.Marker.Interpolate);
		    
                    
                    
                    $(tabVehicles).find("span.text-error").remove();
                    var VehicleInfo = $(tabVehicles).find("span[data-vehicle = '" + data.ID + "']");
                    VehicleInfo.remove();

                    var html = "<span data-vehicle='" + data.ID + "'> <i class='icon-truck'></i> " + data.ID + "  <strong>Speed(KM/Hr)</strong> " + data.Speed 
                    + " <strong>Heading Degrees)</strong> " + Math.round(data.Heading) + " <strong>Time</strong> " + data.Date + "</span>"
                    $(tabVehicles).append(html);

		}
	}
        else
        {
		alert("Your browser does not support WebSockets. You cannot use myClubLink until you upgrade to a modern browser");
	}
    };
    }
}]);

