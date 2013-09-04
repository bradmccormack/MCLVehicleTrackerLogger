'use strict';


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
    .factory("mapService", [function() {
        var LastPosition = {
            Time: new Date(),
            Position: ""
        };

        return {
            Map : {},
            Vehicles: [],

            GetLastPosition: function() {
                return LastPosition;
            },
            UpdateLastPosition: function(Position) {
                LastPosition.Time = new Date();
                LastPosition.Position = Position;
            }
        }
    
}]);
