var genericCallback;

var map = (function(){
	
	var current = {};
	
	var MapQuest = (function(Latitude, Longitude, Zoom, DivID) {
		throw "Not Implemented";
	});
	
	var BingMaps = (function(Latitude, Longitude, Zoom, DivID) {
		throw "Not Implemented";
	});

	var Leaflet = (function(Latitude,Longitude, Zoom, DivID) {
		
		var map;
		var divid;
	    var latlng;
	    var marker;
	    var zoom;
	
		function setview(Latitude, Longitude, Zoom)
		{
	        this.zoom = Zoom;
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
	  
		divid = DivID || "map";
		setview(Latitude || 51.505, Longitude || -0.09, Zoom || 18);
	
		return {
			getMapAPI: function() {
				return mapAPI;
			},
			zoomIn: function() {
	
	        },
	        zoomOut: function() {
	
	        },
	        setView: function(Latitude, Longitude, Zoom) {
				setview(Latitude, Longitude, Zoom);
			},
			addtoRoute: function(Route, Point) {
	        	["Latitude", "Longitude", "Speed", "Heading", "Fix", "DateTime"].forEach(function() {
	        		if(!(this in Point)){
	        			throw "missing params for addtoRoute"
	        		}
	        			
	        		
	        	})
	        },
			setMarker: function(Latitude, Longitude, Text, Color) {
	            latlng = new L.LatLng(Latitude, Longitude);
	            if(!marker) {
	               marker = L.marker(latlng).addTo(map);
	            }
	            else {
	                marker.setLatLng(latlng);
	            }
	            if(Text)
	               marker.bindPopup(Text).openPopup();
			},
			onClick: function(funct) {
				map.on("click", function(e) {
					funct({Location: e.latlng});
				})
				
			}
			
		}
	});
	
	
	var GoogleMaps = (function(Latitude, Longitude, Zoom, DivID) {
	
	    //var mapTypes = { MapTypeId.ROADMAP, MapTypeId.SATELLITE, MapTypeId.HYBRID, MapTypeId.TERRAIN }
	    var map;
	    var divid;
	    var apiKey = "AIzaSyC5wXV9B15WaWQ08qMDD-0O-ZihSnbpi48"; //todo find a way to make this more hidden
	    var latlng;
	    var marker;
	    var zoom;
		var routes = {}; //used to keep track of all polyline routes
	    function init()
	    {
	        var mapProp = {
	            center: new google.maps.LatLng(Latitude, Longitude),
	            zoom: Zoom || 15,
	            mapTypeId: google.maps.MapTypeId.ROADMAP
	        };
	
	        map = new google.maps.Map(document.getElementById(DivID || "map")
	            ,mapProp);
	    }
	
	    function setview(Latitude, Longitude, Zoom)
	    {
	        this.zoom = Zoom;
	        this.Latitude = Latitude;
	        this.Longitude = Longitude;
	        
	        //TODO restrict the zoom level and lat long boundary
	    }
	
	    if(!("google" in window)) {
	
	        genericCallback = init;
	        var url = "https://maps.googleapis.com/maps/api/js?v=3.exp&sensor=true&callback=genericCallback";
	        $.getScript(url, function() {
	            if("google" in window)  {
	                google.maps.visualRefresh = true;
	            }
	        });
	    }
	    else
	    	init();
	    
	
	    return {
	
	        zoomIn : function() {
	            zoom++;
	            map.setZoom(this.zoom);
	        },
	        zoomOut: function() {
	            zoom--;
	            map.setZoom(this.zoom);
	        },
	
	        setView: function(Latitude, Longitude, Zoom) {
	            setview(Latitude, Longitude, Zoom);
	        },
	        addtoRoute: function(Route, Point, Color) {
        		/*
        		["Latitude", "Longitude", "Speed", "Heading", "Fix", "DateTime"].forEach(function() {
        		if(!(this in Point)){
        			throw "missing params for addtoRoute";
        		}
        		*/
	        		
        		if(!(Route in routes)) {
        			var polyOptions = {
					    strokeColor: Color || '#000000',
					    strokeOpacity: 1.0,
					    strokeWeight: 3
					}
				  routes[Route] = new google.maps.Polyline(polyOptions);
				  routes[Route].setMap(map);
        		}
        		
        		var path = routes[Route].getPath();
        		path.push(new google.maps.LatLng(Point.Latitude,Point.Longitude));
	        		
	        	
	        },
	        //ID is the vehicle ID
	        setMarker: function(Latitude, Longitude, Text, Color) {
	         	

//http://chart.apis.google.com/chart?chst=d_map_pin_letter&chld=%E2%80%A2|A37870
	            if(!marker) {
	            	
    				var pinImage = new google.maps.MarkerImage("http://chart.apis.google.com/chart?chst=d_map_pin_letter_withshadow&chld=%E2%80%A2|" + Color,
			       		new google.maps.Size(21, 34),
			        	new google.maps.Point(0,0),
			        	new google.maps.Point(10, 34)
			        );
			    		
	                marker = new google.maps.Marker({
	              	              
				    icon: pinImage,
	                position: new google.maps.LatLng(Latitude, Longitude),
	                map: map});
	                if(Text)
	                    marker.Text = Text;
	                }
	            else {
	                marker.setPosition(new google.maps.LatLng(Latitude, Longitude));
	                if(Text)
	                    marker.setTitle(Text);
	            }
	
	        },
	        onClick: function(funct) {
	            map.on("click", function(e) {
	                funct({Location: e.latlng});
	            })
	
	        }
	
	    }
	});

	
	
	var defaultLocation = { Latitude: -34.50118, Longitude: 150.81071 };
	
	var Settings =
    {
        Vendors: {
        	"Leaflet" : Leaflet,
        	"GoogleMaps": GoogleMaps,
        	"BingMaps": BingMaps,
        	"MapQuest": MapQuest
        }
    };
	
	//setview(Latitude || 51.505, Longitude || -0.09, Zoom || 18);
	//This all probably needs to be cleaned up to be a facade. Instead of returning the Current Map implementation and working with it directly all the methods
	//should be below and interect with the current map vendor selected
	return {
		Current: function() {
			return current;
		},
		SetAPI: function(API) {
			matchingAPI = Settings.Vendors[API];
			if(matchingAPI) {
				current = new matchingAPI(defaultLocation.Latitude, defaultLocation.Longitude, 16, "Mainmap");
			}
		},
		
		
	}
	
   

})();


//Map controls
$("div#Mainmapcontrols button").click(function(){
var Self = $(this);
function notAvailable() {

    Self.popover({
        title: "Not Implemented",
        content: "Sorry, this functionality is not yet available",
        placement: "bottom",
        container: "body"
    }).show();

}

var actions = {
    "mapRefresh" : function() {
      notAvailable();
    },
    "mapZoomIn" : function() {
        notAvailable();
        //System.getMapAPI().zoomIn();
    },
    "mapZoomOut" : function() {
        notAvailable();
        //System.getMapAPI().zoomOut();
    },
    "mapMarker" : function() {
        notAvailable();
    },
    "mapFollow" : function() {
        notAvailable();
    },
    "mapRoute" : function() {
        notAvailable();
    },
    "mapPrint" : function() {
        notAvailable();
    }
}[Self.attr("id")]();
});




