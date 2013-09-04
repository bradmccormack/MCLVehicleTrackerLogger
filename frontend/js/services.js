'use strict';
var genericCallback;

/* Services */
//http://www.ng-newsletter.com/posts/beginner2expert-services.html
angular.module('myApp.services', [])
    .factory("shellService", [function(){

        //TODO pull all the data out from the cookie that is returned via Login process
        var serviceInstance = {
            User: {
                First: "Brad",
                Last: "McCormack",
                Access: 10
            },
            Company: {
               Name: "Test Company",
               MaxUsers: 1,
               Expiry: new Date(),
               Logo: "img/sussex_logo.PNG"
            },
            Settings: {
                Network: {
                    EnableRF: true,
                    Enable3G: true
                },
                Security: {
                    RemoteSupport: false,
                    SystemConsoleAccess: false,
                    AdminPasswordResetOnly: false
                },
                Mobile: {
                    AllowSmartPhone: true,
                    ShowSmartPhoneLocation: false
                },
                Map: {
                    API: "Google Maps", //Contains reference to the current MapAPI in use. The MAP API is a facade over specific concrete implementations
                    Marker: {
                        Smooth: false,
                        FollowCarTrigger: 10 //Every 10 updates the Map system will pan to the selected car
                    }
                }
            }
          
        };
        
        return serviceInstance;
  }])
    .factory("utilityService", [function(){
        return {
                RandomColor: function() {
                        return (function lol(m,s,c){return s[m.floor(m.random() * s.length)] +
                        (c && lol(m,s,c-1));})(Math,'0123456789ABCDEF',4)
        }
              
    }
    
}])
    .factory("mapService", ['shellService', function(shellService) {

        var LastPosition = {
            Time: new Date(),
            Position: ""
        };

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
            zoomIn: function() {

            },
            zoomOut: function() {

            },
            setView: function(Latitude, Longitude, Zoom) {
                setview(Latitude, Longitude, Zoom);
            },
            clearRoutes: function() {
                for(var Route in Routes) {
                    delete Routes[Route];
                }
            },
            clearRoute: function(Route) {
                if(Route in routes)
                    delete routes[Route];

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
        var apiKey = "AIzaSyC5wXV9B15WaWQ08qMDD-0O-ZihSnbpi48"; //TODO let the backend send this over in an encrypted cookie
        var latlng;
        var markers = {};
        var zoom;
        var routes = {}; //used to keep track of all polyline and route information at those points

        function init()
        {
            var mapProp = {
                center: new google.maps.LatLng(Latitude, Longitude),
                zoom: Zoom || 15,
                mapTypeId: google.maps.MapTypeId.ROADMAP
            };

            zoom = Zoom;
            map = new google.maps.Map(document.getElementById(DivID || "map")
                , mapProp);
        }

        function setview(Latitude, Longitude, Zoom)
        {
            zoom = Zoom;
            Latitude = Latitude;
            Longitude = Longitude;

            //TODO restrict the zoom level and lat long boundary
        }

        if(!("google" in window)) {

            genericCallback = init;
            var url = "https://maps.googleapis.com/maps/api/js?v=3.exp&sensor=true&libraries=geometry&callback=genericCallback";
            $.getScript(url, function() {
                if("google" in window)  {
                    google.maps.visualRefresh = true;
                }
            });
        }


        return {

            zoomIn : function() {
                zoom++;
                map.setZoom(zoom);
            },
            zoomOut: function() {
                zoom--;
                map.setZoom(zoom);
            },

            setView: function(Latitude, Longitude, Zoom) {
                setview(Latitude, Longitude, Zoom);
            },
            panTo: function(Latitude, Longitude) {
                map.panTo(new google.maps.LatLng(Latitude, Longitude));
            },
            centerView: function(Latitude, Longitude) {
                map.setCenter(new google.maps.LatLng(Latitude,Longitude));
            },

            clearRoutes: function() {
                for(var Route in routes) {
                    delete routes[Route];
                }
            },
            clearRoute: function(Route) {
                if(Route in routes)
                    delete routes[Route];

            },
            addtoRoute: function(Route, Point, Color) {

                if(!(Route in routes)) {
                    var polyOptions = {
                        strokeColor: Color || Utility.RandomColor(),
                        strokeOpacity: 1.0,
                        strokeWeight: 3
                    }

                    routes[Route]= {polyline : new google.maps.Polyline(polyOptions)};

                    routes[Route].polyline.setMap(map);
                    routes[Route].metadata = {}; //used for looking up date at this time.

                    google.maps.event.addListener(routes[Route].polyline, 'mouseover', function (event) {
                        /* TODO this is tricky. I'm not sure if the co-ordinates that Google is returning are guaranteed to exist in the polyline ..
                         * I might have to find the nearest co-ordinate to get the metadata

                         var key = event.latLng.jb.toString().substring(0,9) + "," + event.latLng.kb.toString().substring(0,9); //Google gives back more detailed co-ords than they were originally stored with.
                         var DateTime = routes[Route].metadata[key];
                         */

                    });

                }
                /*
                 * ["Latitude", "Longitude", "Speed", "Heading", "Fix", "DateTime"].forEach(function() {
                 */
                var path = routes[Route].polyline.getPath();
                path.push(new google.maps.LatLng(Point.Latitude,Point.Longitude));
                //use the lat, long as the key for looking up meta data
                routes[Route].metadata[Point.Latitude + "," + Point.Longitude] = { Lat: Point.Latitude, Long: Point.Longitude, Speed: Point.Speed, Heading: Point.Heading, DateTime: Point.DateTime};


            },
            //ID is the vehicle ID
            setMarker: function(ID, Latitude, Longitude, Text, Color, isInterpolate) {

                //http://chart.apis.google.com/chart?chst=d_map_pin_letter&chld=%E2%80%A2|A37870
                if(!markers[ID]) {

                    var pinImage = new google.maps.MarkerImage("http://chart.apis.google.com/chart?chst=d_map_pin_letter_withshadow&chld=%E2%80%A2|" + Color,
                        new google.maps.Size(21, 34),
                        new google.maps.Point(0,0),
                        new google.maps.Point(10, 34)
                    );

                    markers[ID] = new google.maps.Marker({
                        icon: pinImage,
                        position: new google.maps.LatLng(Latitude, Longitude),
                        map: map,
                        animation: google.maps.Animation.DROP
                    });
                    panTo(markers[ID].position);

                    if(Text) {
                        markers[ID].Text = Text;
                        markers[ID].setTitle(Text);
                    }

                }
                else {
                    if(isInterpolate) {
                        var startLatLng = markers[ID].position;
                        var endLatLng = new google.maps.LatLng(Latitude, Longitude);
                        for(var i = 0 ; i < 1 ; i+=0.1) {
                            var intermediaryPoint =  google.maps.geometry.spherical.interpolate(startLatLng, endLatLng, i);
                            markers[ID].setPosition(intermediaryPoint);
                        }
                    }
                    else {
                        markers[ID].setPosition(new google.maps.LatLng(Latitude, Longitude));

                    }
                }
            },
            onClick: function(funct) {
                map.on("click", function(e) {
                    funct({Location: e.latlng});
                })

            },
            refresh: function() {
                init();
                //google.maps.event.trigger(map, 'resize');
                //map.setCenter(new google.maps.LatLng(Latitude,Longitude));
            }

        }
    });



        var CurrentMapAPI;
        var Vendors =  {
            "Leaflet" : Leaflet,
            "GoogleMaps": GoogleMaps,
            "BingMaps": BingMaps,
            "MapQuest": MapQuest
        };

        var defaultLocation = {
            Latitude: -34.50118,
            Longitude: 150.81071,
            Zoom: 16};
        return {

            //Facade
            Map : {
                SetAPI: function(API) {
                    var matchingAPI = Vendors[API];
                    if(matchingAPI) {
                        if(CurrentMapAPI)
                            CurrentMapAPI = undefined;
                        CurrentMapAPI = new matchingAPI(defaultLocation.Latitude, defaultLocation.Longitude, defaultLocation.Zoom, "MapCanvas");
                    }
                },
                Refresh: function() {
                    if(!CurrentMapAPI)
                        CurrentMapAPI = new Vendors[shellService.Settings.Map.API.replace(" ","")](defaultLocation.Latitude, defaultLocation.Longitude, defaultLocation.Zoom, "MapCanvas");
                    else
                        CurrentMapAPI.refresh();
                },
                ZoomIn: function(){
                    CurrentMapAPI.zoomIn();
                },
                ZoomOut: function(){
                    CurrentMapAPI.zoomOut();
                }
            },
            Vehicles: [
            ],
            GetLastPosition: function() {
                return LastPosition;
            },
            UpdateLastPosition: function(Position) {
                LastPosition.Time = new Date();
                LastPosition.Position = Position;
            }

        }

}]);
