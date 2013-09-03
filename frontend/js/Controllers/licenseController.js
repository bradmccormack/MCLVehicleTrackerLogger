angular.module('myApp.controllers').controller("licenseController", ['$scope', function($scope){
   
   function Init() {
       
   }
   
   $scope.Company = {
        Name: "Test Company",
        MaxUsers: 1,
        Expiry: new Date()
   }
   
   $scope.User = {
        First : "Brad",
        Last : "Mccormack"
   }
  
   Init();
   
   
   
}]);
