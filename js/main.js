/*
 * TODO figure out a nice abstraction over the Map API vendors so I can swap around simply by changing the active API in the backend.
 * Some clients may want functionality that is available in paid Google Maps API so we can swap over to that easily
 * 
 * http://www.bing.com/maps/ I don't mind the slider (map expand widget) add that 
 */



var MapQuest = (function(Latitude, Longitude, Zoom, DivID) {
	throw "Not Implemented";
});

var Bing = (function(Latitude, Longitude, Zoom, DivID) {
	throw "Not Implemented";
});

var Google = (function(Latitude, Longitude, Zoom, DivID) {
	throw "Not Implemented";
});

var Leaflet = (function(Latitude,Longitude, Zoom, DivID) {
	
	var map;
	var divid;
	
	function setview(Latitude, Longitude, Zoom)
	{
		this.Latitude = Latitude;
		this.Longitude = Longitude;
		map = L.map(divid).setView([Latitude, Longitude], Zoom);
	    
	    /*
		L.tileLayer('http://{s}.tile.cloudmade.com/6c45216fc160453e9dec40f0f9cd1312/997/256/{z}/{x}/{y}.png', {
	    attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://cloudmade.com">CloudMade</a>',
	    maxZoom: 13
		}).addTo(map);
	    */
	    
	   
	    L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
	    }).addTo(map);
	    
	}
    
    /*
    if(typeof(L) == "undefined") {
    	Console.log("The leaflet JS resource has not loaded!");
    	return;
    }
    */
  
	divid = DivID || "map";
	setview(Latitude || 51.505, Longitude || -0.09, Zoom || 18);
    

	
	return {
		setView: function(Latitude, Longitude, Zoom) {
			setview(Latitude, Longitude, Zoom);
		},
		setMarker: function(Latitude, Longitude, Text) {
			L.marker([Latitude, Longitude]).addTo(map).bindPopup(Text).openPopup();
		},
		onClick: function(funct) {
			map.on("click", function(e) {
				funct({Location: e.latlng});
			})
			
		}
		
	}
});

var MapAPI = 
{
	Active : Leaflet,
	Vendors: [Leaflet, Google, Bing, MapQuest]
	
};


var System = (function(){
	
	var Colours = [];
	
	return {
		

		
		updateLegend: function(JSON) {
			//We need to know if the GPS signal is correct or not (Fix status is true)
			if(JSON.Vehicles.length > 0) {
				var Legend = $("div#Mainlegend div#Vehiclelegend ul");
			
	
				JSON.Vehicles.forEach(function(Vehicle){
					Legend.append('<li><a href="#"><i class="icon-truck"></i> ' + Vehicle +'</a></li>');
				});
			
			}
			else {
				
			}
		}
	}
});


$(document).ready(function() {
	
	var defaultLocation = { Latitude: 34.50094, Longitude: 150.81060 };
	var mapAPI = new MapAPI.Active(51.505, -0.09, 18, "Mainmap");
	var system = new System();
	
	//add a couple of vehicles in hard coded for now
	system.updateLegend({Vehicles: ["Mitsubishi Bus", "Izusu Bus"]});
	
	
	
	/*
	mapAPI.onClick(function(e){
		alert("Clicked at " + e.Location);
	})
	*/


	
});
