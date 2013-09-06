angular.module('myApp.controllers').controller("mainController", ['$scope', 'networkService', function($scope, networkService){
   
   function Login() {
     
     var cookies = $.cookie();
        if("session" in $.cookie())
        {
                if("success" in cbobj && typeof cbobj.success == "function") {
                        cbobj.success("Craig Smith"); //TODO pull the name from the session cookie
                }
                return;
        }
   }
   

   
   $scope.Camera = {
    
   };
   
   $scope.Vehicles = [];





   $scope.SystemInit = function() {
       Login();
       alert("lol");

   }
   
   
   
}]);
