angular.module('myApp.controllers',[]).controller("mainController", ['$scope', function($scope){
   
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
   
   function Init() {
       Login();
   }
   
   
   $scope.Camera = {
    
   };
   
   $scope.Vehicles = [];
   
   Init();
   
   
   
   
}]);
