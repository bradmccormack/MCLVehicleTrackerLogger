/*Pass in the shellService via Dependency injection rather than rely on prototypical inheritence between controllers to access User information */

angular.module('myApp.controllers',['http-auth-interceptor']).controller("headerController", ['$scope', 'shellService', function($scope, shellService){
   
    $scope.clock = {
        interval: 1000,
        time: ""
   }
   
   $scope.User = {
    First: shellService.User.First,
    Last: shellService.User.Last
   }
   
    var logout = function() {
        
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
