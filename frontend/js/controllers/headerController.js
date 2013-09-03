myApp.controller("headerController", ['$scope', function($scope){
   
    $scope.clock = {
        interval: 1000,
        time: ""
   }
   
    var logout = function() {
        
    }
    
    var edit = function() {
        
    }
   
    var updateClock = function() {
        $scope.clock.time = new Date().toTimeString();
    }
    
 
   var timer = setInterval(function() {
        $scope.$apply(updateClock);
   }, $scope.clock.interval);
   updateClock();
   
}]);
