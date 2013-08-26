app.controller("mapControler", ['$scope', function($scope){
   
    $scope.MapAPI = {};
    $scope.Camera = { Snap: true, SnapCount: 0, SnapTrigger: 10 };
   
    $scope.updateLegend = function(VehicleID, Color)
    {
  
        //We need to know if the GPS signal is correct or not (Fix status is true)
      	
        var Legend = $("div#Mainlegend div#Vehiclelegend");
        Legend.find("span.text-error").remove();
      
        Vehicles[VehicleID] = {
        	DateTime: new Date(),
        	Color: Color
        }
        Legend.find("ul li#vehicle_" + VehicleID).remove();
        Legend.find("ul").append('<li id=vehicle_' + VehicleID + '><a href="#" style="color: #' + Vehicles[VehicleID].Color + '"><i class="icon-truck"></i> ' + VehicleID +'</a></li>');
          //mapAPI.Current().setMarkerColor(Vehicles[VehicleID].Color);
        //remove everything from the legend if there has been no contact in over 1 hour or whatever
        
    }
    
    /*
     if(Camera.Snap)
                    {
                            Camera.SnapCount++;
                            if(Camera.SnapCount == Camera.SnapTrigger)
                            {
                                    mapAPI.Current().panTo(data.Latitude, data.Longitude);
                                    Camera.SnapCount = 0;
                            }
                    
                    }
                    */
    
}]);

