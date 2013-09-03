'use strict';


/* Services */
//http://www.ng-newsletter.com/posts/beginner2expert-services.html
angular.module('myApp.services', [])
    .factory("shellService", [function(){
    
        var serviceInstance = {
            Map: {
                API: "Google Maps", //Contains reference to the current MapAPI in use. The MAP API is a facade over specific concrete implementations
                Marker: {
                    Smooth: false,
                    FollowCarTrigger: 10 //Every 10 updates the Map system will pan to the selected car
                }
            },
            Network: {
                EnableRF: true,
                Enable3G: true  
            },
            User: {
                First: "Brad",
                Last: "McCormack",
                Access: 10
            },
            Company: {
               Name: "Test Company",
               MaxUsers: 1,
               Expiry: new Date()  
            },
            Security: {
                RemoteSupport: false,
                SystemConsoleAccess: false,
                AdminPasswordResetOnly: false
            },
            Mobile: {
                 AllowSmartPhone: true,
                 ShowSmartPhoneLocation: false
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
        return {
            Map : {}
        }
    
}]);
