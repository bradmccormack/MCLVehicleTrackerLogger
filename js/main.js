/*
 * TODO figure out a nice abstraction over the Map API vendors so I can swap around simply by changing the active API in the backend.
 * Some clients may want functionality that is available in paid Google Maps API so we can swap over to that easily
 */



var MapAPI = 
{
	Active : Leaflet,
	Vendors: [Leaflet, Google, Bing, MapQuest]
	
};

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
	    L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png').addTo(map);
	}
    
    /*
    if(typeof(L) == "undefined") {
    	Console.log("The leaflet JS resource has not loaded!");
    	return;
    }
    */
  
	divid = DivID || "map";
	setview(Latitude || 51.505, Longitude || -0.09, Zoom || 13);
    

	
	return {
		setView: function(Latitude, Longitude, Zoom) {
			setview(Latitude, Longitude, Zoom);
		},
		setMarker: function(Latitude, Longitude, Text) {
			L.marker([Latitude, Longitude]).addTo(map).bindPopup(Text).openPopup();
		}
		
		
	}
});





$(document).ready(function() {
	
	var MapAPI = new MapAPI.Active(51.505, -0.09, 13, "Mainmap");



	
});
