var genericCallback;


var MapQuest = (function(Latitude, Longitude, Zoom, DivID) {
	throw "Not Implemented";
});

var BingMaps = (function(Latitude, Longitude, Zoom, DivID) {
	throw "Not Implemented";
});

var GoogleMaps = (function(Latitude, Longitude, Zoom, DivID) {

    //var mapTypes = { MapTypeId.ROADMAP, MapTypeId.SATELLITE, MapTypeId.HYBRID, MapTypeId.TERRAIN }
    var map;
    var divid;
    var apiKey = "AIzaSyC5wXV9B15WaWQ08qMDD-0O-ZihSnbpi48"; //todo find a way to make this more hidden
    var latlng;
    var marker;
    var zoom;

    function init()
    {
        var mapProp = {
            center: new google.maps.LatLng(Latitude, Longitude),
            zoom: Zoom || 15,
            mapTypeId: google.maps.MapTypeId.HYBRID
        };

        map = new google.maps.Map(document.getElementById(DivID || "map")
            ,mapProp);
    }

    function setview(Latitude, Longitude, Zoom)
    {
        this.zoom = Zoom;
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
        setMarker: function(Latitude, Longitude, Text) {
            if(!marker) {
                marker = new google.maps.Marker({
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
		zoomIn: function() {

        },
        zoomOut: function() {

        },
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


var System = (function(){

	var Self = this;
    var MapAPI =
    {
        Active: GoogleMaps,
        Vendors: [Leaflet, GoogleMaps, BingMaps, MapQuest]
    };


    var mapAPI;
    var Con;
	var Colours = [];
    var Vehicles = [];

    //System wide settings
    var Settings = { Marker: { InterpolateCount : 10}};
    var Position = { Last: {}};

    function updateLegend(JSON) {
        //We need to know if the GPS signal is correct or not (Fix status is true)
        if(JSON.Vehicles.length > 0) {
            var Legend = $("div#Mainlegend div#Vehiclelegend ul");

            //Add any Vehicles to the Legend that are not there .. If GPS fix is not true then show warning etc
            //if there is no information for Vechile BLAH
            JSON.Vehicles.forEach(function(Vehicle){
                Legend.append('<li><a href="#"><i class="icon-truck"></i> ' + Vehicle +'</a></li>');
            });

        }
        else {
            //remove everything from the legend if there has been no contact in over 1 hour or whatever
        }
    }

	return {
        getMapAPI: function() {
          return this.mapAPI;
        },
		init: function() {

            var defaultLocation = { Latitude: -34.50118, Longitude: 150.81071 };
            Self.mapAPI = new MapAPI.Active(defaultLocation.Latitude, defaultLocation.Longitude, 16, "Mainmap");


            //add a couple of vehicles in hard coded for now
            //system.updateLegend({Vehicles: ["Mitsubishi Bus", "Izusu Bus"]});

            if (window["WebSocket"]) {
                    //alert("Browser supports Web Sockets. Yay");
                    Con = new WebSocket("ws://dev.myclublink.com.au/ws");
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

                        Self.mapAPI.setMarker(cords[0],cords[1]);
                        /*
                        //TODO interpolate between cords
                        var X = cords[0];
                        var y = cords[1];
                        if(Position.Last.Latitude) {

                            var m = (Y - Position.Last.Longitude) / ( X - Position.Last.Latitude); //gradient
                            var b = 0; //figure out y intercept
                            var increment = X - Position.Last.Latitude / Settings.Marker.InerpolateCount;

                            for(var i = 0; i < Settings.Marker.InterpolateCount; i++) {
                                //y = mx+b
                                X += increment;
                                Y = m * X + b;
                                mapAPI.setMarker(X, Y);
                            }

                        }
                        else {
                            mapAPI.setMarker(cords[0],cords[1]); // remove this ?
                        }

                        Position.Last.Latitude = X;
                        Position.Last.Longitude = Y;
                        */

                        //alert("Message received " + evt.data);
                        //appendLog($("<div/>").text(evt.data))
                    }
            } else {
                alert("Your browser does not support WebSockets. You cannot use myClubLink until you upgrade to a modern browser");
            }
        }
	}
});

function bindHandlers() {
    
    //Login
    $('#myModal').on("shown", function()
    {
        $(".username").focus(function() {
            $(".user-icon").css("left","-48px");
        });
        $(".username").blur(function() {
            $(".user-icon").css("left","0px");
        });

        $(".password").focus(function() {
            $(".pass-icon").css("left","-48px");
        });
        $(".password").blur(function() {
            $(".pass-icon").css("left","0px");
        });
    });


	//Menu Nav
	$("#tabMap,#tabProfile, #tabSupport, #tabReports, #tabSettings").click(function(){
        var Self = $(this);
        var Main = $("#Mainmap");
        var actions = {
            tabMap: function() {
                throw "Not Implemented";
            },
            tabProfile: function() {
                throw "Not Implemented";

            },
            tabSupport: function() {
                throw "Not Implemented";
            },
            tabReports: function() {
                throw "Not Implemented";
            },
            tabSettings: function() {
                $.ajax({
                    type: "GET",
                    url: "/system/settings",
                    dataType: "HTML",
                    success: function(HTML) {
                        Main.replaceWith(HTML);
                    },
                    error: function(a,b,c) {
                      //TODO
                    }
                })
                Main.html("<p>Lol</p>");

            }
        }[Self.attr("id")]();
	});


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


}

$(document).ready(function() {

    bindHandlers();

    //Startup the main system
    var system = new System();
    system.init();

 	$('#myModal').modal('toggle');  //Perform login

    /*
	mapAPI.onClick(function(e){
		alert("Clicked at " + e.Location);
	})
	*/




	
});
