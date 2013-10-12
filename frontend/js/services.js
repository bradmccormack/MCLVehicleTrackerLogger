'use strict';
var genericCallback;
var directionsService;

/* Services */
//http://www.ng-newsletter.com/posts/beginner2expert-services.html
angular.module('myApp.services', [])
    .factory("shellService", ['$rootScope', function ($rootScope) {


    var serviceInstance = {
        User:{
            First:"",
            Last:"",
            Access:0,
	        Password: "",
	        Email: ""
        },
        Company:{
            Name:"",
            MaxUsers:1,
            Expiry:new Date(),
            Logo:""
        },
        Settings:{
            Network:{
                EnableRF: false,
                Enable3G: false
            },
            Security:{
                RemoteSupport:false,
                SystemConsoleAccess:false,
                AdminPasswordResetOnly:false
            },
            Mobile:{
                AllowSmartPhone: false,
                ShowSmartPhoneLocation:false
            },
            Map:{
                API:"", //Contains reference to the current MapAPI in use. The MAP API is a facade over specific concrete implementations
                Marker:{
                    Smooth:false,
	                SnaptoRoad: false,
                    FollowVehicleTrigger:0 //Default Vehicle camera trigger. 10 would represent a camera pan after every 10 movement updates
                },
	            Camera : {
		            //Keeps track of settings per vehicle in the system
	            }
            }
        },
	    Messages: [],
	    ClearConfig: function() {
			this.User = "";

		  
	    },
	    LoadConfig: function(data) {
			this.User = {
				First: data.user.Firstname,
				Last: data.user.Lastname,
				Password: data.user.Password, //TODO encrypt serverside
				Access: data.user.AccessLevel,
				Email: data.user.Email,
				LoggedIn: true
			};

			this.Company = {
				Name: data.company.Name,
				MaxUsers: data.company.Maxusers,
				Expiry: data.company.Expiry,
				Logo: data.company.LogoPath
			}
			this.Settings = {
				Network: {
					EnableRF: data.settings.RadioCommunication,
					Enable3G: data.settings.DataCommunication
				},
				Security: {
					RemoteSupport: data.settings.SecurityRemoteAdmin,
					SystemConsoleAccess: data.settings.SecurityRemoteConsoleAccess,
					AdminPasswordResetOnly: data.settings.SecurityAdminPasswordReset
				},
				Mobile: {
					AllowSmartPhone: data.settings.MobileSmartPhoneAccess,
					ShowSmartPhoneLocation: data.settings.MobileShowBusLocation
				},
				Map: {
					API: data.settings.MapAPI,
					Marker: {
						Smooth: data.settings.Interpolate,
						SnaptoRoad: data.settings.SnaptoRoad,
						FollowVehicleTrigger: data.settings.CameraPanTrigger
					},
					Camera: {
					}
				},
				
			}
		    $rootScope.$broadcast("ConfigChanged", serviceInstance);
		}
    };

    return serviceInstance;
}])
    .factory("utilityService", [function () {
    return {
        RandomColor:function () {
            return (function lol(m, s, c) {
                return s[m.floor(m.random() * s.length)] +
                    (c && lol(m, s, c - 1));
            })(Math, '0123456789ABCDEF', 4)
        }

    }

}])
    .factory("mapService", ['shellService', 'utilityService','$rootScope', function (shellService, utilityService, $rootScope) {

    var LastPosition = {
        Time: new Date(),
        Position: {
            Latitude: undefined,
            Longitude: undefined
        }
    };

    var MapQuest = (function (Latitude, Longitude, Zoom, DivID) {
        throw "Not Implemented";
    });

    var BingMaps = (function (Latitude, Longitude, Zoom, DivID) {
        throw "Not Implemented";
    });


    var Leaflet = (function (Latitude, Longitude, Zoom, DivID) {

        var map;
        var divid;
        var latlng;
        var marker;
        var zoom;

        function setview(Latitude, Longitude, Zoom) {
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
            zoomIn:function () {

            },
            zoomOut:function () {

            },
            setView:function (Latitude, Longitude, Zoom) {
                setview(Latitude, Longitude, Zoom);
            },
            clearRoutes:function () {
                for (var Route in Routes) {
                    delete Routes[Route];
                }
            },
            clearRoute:function (Route) {
                if (Route in routes)
                    delete routes[Route];

            },
            addtoRoute:function (Route, Point) {
                ["Latitude", "Longitude", "Speed", "Heading", "Fix", "DateTime"].forEach(function () {
                    if (!(this in Point)) {
                        throw "missing params for addtoRoute"
                    }
                })
            },
            setMarker:function (Latitude, Longitude, Text, Color) {
                latlng = new L.LatLng(Latitude, Longitude);
                if (!marker) {
                    marker = L.marker(latlng).addTo(map);
                }
                else {
                    marker.setLatLng(latlng);
                }
                if (Text)
                    marker.bindPopup(Text).openPopup();
            },
            onClick:function (funct) {
                map.on("click", function (e) {
                    funct({Location:e.latlng});
                })

            }
        }
    });

    var GoogleMaps = (function (Latitude, Longitude, Zoom, DivID) {

        //var mapTypes = { MapTypeId.ROADMAP, MapTypeId.SATELLITE, MapTypeId.HYBRID, MapTypeId.TERRAIN }
        var map;
        var divid;
        var apiKey = "AIzaSyC5wXV9B15WaWQ08qMDD-0O-ZihSnbpi48"; //TODO let the backend send this over in an encrypted cookie
        var latlng;
        var markers = {};
        var zoom;
        var routes = {}; //used to keep track of all polyline and route information at those points

        function init() {
            var mapProp = {
                center:new google.maps.LatLng(Latitude, Longitude),
                zoom:Zoom || 15,
                mapTypeId:google.maps.MapTypeId.ROADMAP
            };

            zoom = Zoom;
            map = new google.maps.Map(document.getElementById(DivID || "map")
                , mapProp);
	        directionsService =  new google.maps.DirectionsService();

        }

        function setview(Latitude, Longitude, Zoom) {
            zoom = Zoom;
            Latitude = Latitude;
            Longitude = Longitude;

            //TODO restrict the zoom level and lat long boundary
        }

        if (!("google" in window)) {

            genericCallback = init;
            var url = "https://maps.googleapis.com/maps/api/js?v=3.exp&sensor=true&libraries=geometry&callback=genericCallback";
            $.getScript(url, function () {
                if ("google" in window) {
                    google.maps.visualRefresh = true;

                }
            });
        }
        return {

            zoomIn:function () {
                zoom++;
                map.setZoom(zoom);
            },
            zoomOut:function () {
                zoom--;
                map.setZoom(zoom);
            },

            setView:function (Latitude, Longitude, Zoom) {
                setview(Latitude, Longitude, Zoom);
            },
            panTo:function (Latitude, Longitude) {
                map.panTo(new google.maps.LatLng(Latitude, Longitude));
            },
            centerView:function (Latitude, Longitude) {
                map.setCenter(new google.maps.LatLng(Latitude, Longitude));
            },

            clearRoutes:function () {
                for (var Route in routes) {
                    delete routes[Route];
                }
            },
            clearRoute:function (Route) {
                if (Route in routes)
                    delete routes[Route];

            },
            addtoRoute:function (Vehicle, Point) {

                if (!(Vehicle in routes)) {
                    var polyOptions = {
                        strokeColor: Vehicles[Vehicle].Color,
                        strokeOpacity:1.0,
                        strokeWeight:3
                    }

                    routes[Vehicle] = {polyline:new google.maps.Polyline(polyOptions)};

                    routes[Vehicle].polyline.setMap(map);
                    routes[Vehicle].metadata = {}; //used for looking up date at this time.

                    google.maps.event.addListener(routes[Vehicle].polyline, 'mouseover', function (event) {
                        /* TODO this is tricky. I'm not sure if the co-ordinates that Google is returning are guaranteed to exist in the polyline ..
                         * I might have to find the nearest co-ordinate to get the metadata

                         var key = event.latLng.jb.toString().substring(0,9) + "," + event.latLng.kb.toString().substring(0,9); //Google gives back more detailed co-ords than they were originally stored with.
                         var DateTime = routes[Vehicle].metadata[key];
                         */

                    });

                }
                /*
                 * ["Latitude", "Longitude", "Speed", "Heading", "Fix", "DateTime"].forEach(function() {
                 */
                var path = routes[Vehicle].polyline.getPath();
                path.push(new google.maps.LatLng(Point.Latitude, Point.Longitude));
                //use the lat, long as the key for looking up meta data
                routes[Vehicle].metadata[Point.Latitude + "," + Point.Longitude] = { Lat:Point.Latitude, Long:Point.Longitude, Speed:Point.Speed, Heading:Point.Heading, DateTime:Point.DateTime};


            },
            //ID is the vehicle ID
            setMarker:function (ID, Latitude, Longitude, Text, Color) {

                //http://chart.apis.google.com/chart?chst=d_map_pin_letter&chld=%E2%80%A2|A37870
                if (!markers[ID]) {

                    var pinImage = new google.maps.MarkerImage("http://chart.apis.google.com/chart?chst=d_map_pin_letter_withshadow&chld=%E2%80%A2|" + Color,
                        new google.maps.Size(21, 34),
                        new google.maps.Point(0, 0),
                        new google.maps.Point(10, 34)
                    );

	                //Create new marker for this vehicle
                    markers[ID] = new google.maps.Marker({
                        icon:pinImage,
                        position:new google.maps.LatLng(Latitude, Longitude),
                        map:map,
                        animation:google.maps.Animation.DROP
                    });
                }

                if (Text) {
                    markers[ID].Text = Text;
                    markers[ID].setTitle(Text);
                }


                else {
	                //attempt to interpolate movement
                    if (shellService.Settings.Map.Marker.Smooth) {
                        var startLatLng = markers[ID].position;
                        var endLatLng = new google.maps.LatLng(Latitude, Longitude);
                        for (var i = 0; i < 1; i += 0.1) { //10 intermediary points
                            var intermediaryPoint = google.maps.geometry.spherical.interpolate(startLatLng, endLatLng, i);
                            markers[ID].setPosition(intermediaryPoint);
                        }
                    }
                    else {
                        markers[ID].setPosition(new google.maps.LatLng(Latitude, Longitude));
                    }

				    //Get camera settings for current vehicle
				    var Camera = shellService.Settings.Map.Camera[ID];

				    if (Camera.Snap) {
					    Camera.SnapCount++;
					    if (Camera.SnapCount == shellService.Settings.Map.Marker.FollowVehicleTrigger) {
						    map.panTo(new google.maps.LatLng(Latitude, Longitude));
						    Camera.SnapCount = 0;
					    }

				    }
                }
            },
            onClick:function (funct) {
                map.on("click", function (e) {
                    funct({Location:e.latlng});
                })

            },
            refresh:function () {
                google.maps.event.trigger(map, 'resize');
                map.setCenter(new google.maps.LatLng(Latitude, Longitude));
            },
            reload:function () {
                init();
            }

        }
    });

    var CurrentMapAPI;
    var Vendors = {
        "Leaflet":Leaflet,
        "GoogleMaps":GoogleMaps,
        "BingMaps":BingMaps,
        "MapQuest":MapQuest
    };

    var defaultLocation = {
        Latitude:-34.50118,
        Longitude:150.81071,
        Zoom:16};

	var Vehicles = {
	};
	var VehiclesCount = 0;
	var Messages = [];

		return {

        //Facade
        Map:{
            SetAPI:function (API) {
                var matchingAPI = Vendors[API];
                if (matchingAPI) {
                    if (CurrentMapAPI)
                        CurrentMapAPI = undefined;
                    CurrentMapAPI = new matchingAPI(defaultLocation.Latitude, defaultLocation.Longitude, defaultLocation.Zoom, "MapCanvas");
                }
            },
            Refresh:function () {
                CurrentMapAPI.refresh();
            },
            ReLoad:function () {
                if (!CurrentMapAPI)
                    CurrentMapAPI = new Vendors[shellService.Settings.Map.API.replace(" ", "")](defaultLocation.Latitude, defaultLocation.Longitude, defaultLocation.Zoom, "MapCanvas");
                else
                    CurrentMapAPI.reload();
            },
            ZoomIn:function () {
                CurrentMapAPI.zoomIn();
            },
            ZoomOut:function () {
                CurrentMapAPI.zoomOut();
            },
	        SetMarker: function(ID, Latitude, Longitude, Text) {
		        if(!(ID in Vehicles)) {
			        Vehicles[ID] = {
				        Ref: ID,
				        Latitude: Latitude,
				        Longitude: Longitude,
				        Color: utilityService.RandomColor(),
				        Selected: false
			        };
			        VehiclesCount++;

			        //Create a camera object for this vehicle
			        shellService.Settings.Map.Camera[ID] = {
				        SnapCount: 0,
				        Snap: true
			        }
			        $rootScope.$broadcast("LegendChange", {Count: VehiclesCount, Vehicles: Vehicles});
		        }
		        var Src = new google.maps.LatLng(Latitude, Longitude);

		        //Note - SnaptoRoad chokes if you pound the system with updates (eg 20 per second versus 1 because of call stack - TODO block on vehicle updates if waiting for snap
		        //to road for this vehicle
		        /*
			if(shellService.Settings.Map.Marker.SnaptoRoad){
					directionsService.route(
						{   origin: Src,
							destination: Src,
							travelMode: google.maps.DirectionsTravelMode.DRIVING
					}, function(response,status){
							if(status == google.maps.DirectionsStatus.OK){
								var pos = response.routes[0].legs[0].start_location;
								CurrentMapAPI.setMarker(ID, pos.pb, pos.qb, Text, Vehicles[ID].Color);
							}
						});
		        }
				else {
		      
		            CurrentMapAPI.setMarker(ID, Latitude, Longitude, Text, Vehicles[ID].Color);
		        }
			    */
			CurrentMapAPI.setMarker(ID, Latitude, Longitude, Text, Vehicles[ID].Color);
	        },
	        AddtoRoute:function (Vehicle, Point)
	        {
		        if(!(Vehicle in Vehicles)){
			        Vehicles[Vehicle] = {
				        Ref: Vehicle,
				        Latitude: Point.Latitude,
				        Longitude: Point.Longitude,
				        Color: utilityService.RandomColor(),
				        Selected: false
			        };
			        VehiclesCount++;

			        //Create a camera object for this vehicle
			        shellService.Settings.Map.Camera[Vehicle] = {
				        SnapCount: 0,
				        Snap: false
			        }
			        $rootScope.$broadcast("LegendChange", {Count: VehiclesCount, Vehicles: Vehicles});
		        }
		        CurrentMapAPI.addtoRoute(Vehicle, Point, Vehicles[Vehicle].Color);
	        }
        },
	    GetVehicles: function() {
		    return Vehicles;
	    },
		GetVehicleCount: function() {
			return VehiclesCount;
		},
        GetLastPosition:function () {
            return LastPosition;
        },
        UpdateLastPosition:function (Position) {
            LastPosition.Time = new Date();
            LastPosition.Position = Position;
        }

    }

}]).factory("networkService", ['mapService', 'utilityService', '$rootScope', function (mapService, utilityService, $rootScope) {

    return {
        Init: function () {
            var Con;

            if (window["WebSocket"]) {
                Con = new WebSocket("ws://dev.myclublink.com.au:8080/ws");

                Con.onopen = function () {
	                $rootScope.$broadcast("systemMessage", "Connected to server");
                };

                Con.onclose = function (evt) {
	                $rootScope.$broadcast("systemMessage", "Server connection closed");
                }
                Con.onmessage = function (evt) {
                    var data = JSON.parse(evt.data).Entry;
                    $rootScope.$broadcast("positionChange", data); //send the data out to any listeners
                }
            } else {
                alert("Your browser does not support WebSockets. You cannot use myClubLink until you upgrade to a modern browser");
            }
        }
    }


}]);
