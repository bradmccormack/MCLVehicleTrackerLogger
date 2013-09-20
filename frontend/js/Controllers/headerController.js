/*Pass in the shellService via Dependency injection rather than rely on prototypical inheritance between controllers to access User information */
/*Pass in http-auth-interceptor and ngCookies modules as module dependencies too */

angular.module('myApp.controllers',['http-auth-interceptor', 'ngCookies']).controller("headerController", ['$scope', 'shellService', function($scope, shellService){
   
    $scope.clock = {
        interval: 1000,
        time: ""
   }
   
   $scope.User = {
    First: shellService.User.First,
    Last: shellService.User.Last
   }
   
    var logout = function() {
        //do HTTP to clear cookie then redirect them back to login
    }
    
    var edit = function() {
        
    }
   
    var updateClock = function() {
        $scope.clock.time = new Date().toLocaleString();
    }
    
 
   var timer = setInterval(function() {
        $scope.$apply(updateClock);
   }, $scope.clock.interval);
   updateClock();
   
}]);
