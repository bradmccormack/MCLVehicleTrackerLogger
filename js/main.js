/*
 * TODO figure out a nice abstraction over the Map API vendors so I can swap around simply by changing the active API in the backend.
 * Some clients may want functionality that is available in paid Google Maps API so we can swap over to that easily
 * 
 * http://www.bing.com/maps/ I don't mind the slider (map expand widget) add that 
 */

var genericCallback;

var MapQuest = (function(Latitude, Longitude, Zoom, DivID) {
	throw "Not Implemented";
});

var BingMaps = (function(Latitude, Longitude, Zoom, DivID) {
	throw "Not Implemented";
});

var GoogleMaps = (function(Latitude, Longitude, Zoom, DivID) {
    var map;
    var divid;
    var apiKey = "AIzaSyC5wXV9B15WaWQ08qMDD-0O-ZihSnbpi48"; //todo find a way to make this more hidden
    var latlng;
    var marker;


    function init()
    {
        var mapProp = {
            center:new google.maps.LatLng(Latitude, Longitude),
            zoom: Zoom || 15,
            mapTypeId:google.maps.MapTypeId.ROADMAP
        };

        map = new google.maps.Map(document.getElementById(DivID || "map")
            ,mapProp);
    }

    function setview(Latitude, Longitude, Zoom)
    {
        this.Latitude = Latitude;
        this.Longitude = Longitude;
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
    {
        init();
    }



    setview(Latitude || 51.505, Longitude || -0.09, Zoom || 18);

    return {
        setView: function(Latitude, Longitude, Zoom) {
            setview(Latitude, Longitude, Zoom);
        },
        setMarker: function(Latitude, Longitude, Text) {


            //if(Text)
            //    .bindPopup(Text).openPopup();
        },
        onClick: function(funct) {
            map.on("click", function(e) {
                funct({Location: e.latlng});
            })

        }

    }
});

var Leaflet = (function(Latitude,Longitude, Zoom, DivID) {
	
	var map;
	var divid;
    var latlng;
    var marker;


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
  
	divid = DivID || "map";
	setview(Latitude || 51.505, Longitude || -0.09, Zoom || 18);



	return {
		setView: function(Latitude, Longitude, Zoom) {
			setview(Latitude, Longitude, Zoom);
		},
		setMarker: function(Latitude, Longitude, Text) {
            latlng = new L.LatLng(Latitude, Longitude);
            if(!marker) {
               marker = L.marker(latlng).addTo(map);
            }
            else {
                marker.setLatLng(latlng);
            }


            //if(Text)
            //    .bindPopup(Text).openPopup();
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
	//Active : Leaflet,
    Active: GoogleMaps,
	Vendors: [Leaflet, GoogleMaps, BingMaps, MapQuest]
	
};


var System = (function(){

    var Con;
	var Colours = [];
	
	return {
		init: function() {
            if (window["WebSocket"]) {
                    //alert("Browser supports Web Sockets. Yay");
                    Con = new WebSocket("ws://myclublink.com.au:8080/ws");
                    //if(Con)
                    //    Con.send("test message");


                    Con.onopen = function() {
                        //Con.send("test message");
                    };

                    Con.onclose = function(evt) {
                        //alert("Closing web socket");
                        //appendLog($("<div><b>Connection closed.</b></div>"))
                    }
                    Con.onmessage = function(evt) {
                        var cords = evt.data.split(",");

                        mapAPI.setMarker(cords[0],cords[1]);

                        //alert("Message received " + evt.data);
                        //appendLog($("<div/>").text(evt.data))
                    }
            } else {
                alert("Your browser does not support WebSockets. You cannot use myClublink until you upgrade to a modern browser");
            }
        },

		
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

var mapAPI = {};

$(document).ready(function() {

    var system = new System();
    system.init();


    //OpText
	var defaultLocation = { Latitude: -34.50118, Longitude: 150.81071 };

    mapAPI = new MapAPI.Active(defaultLocation.Latitude, defaultLocation.Longitude, 16, "Mainmap");


	//add a couple of vehicles in hard coded for now
	system.updateLegend({Vehicles: ["Mitsubishi Bus", "Izusu Bus"]});
	
	
	
    /*
	mapAPI.onClick(function(e){
		alert("Clicked at " + e.Location);
	})
	*/




	
});
