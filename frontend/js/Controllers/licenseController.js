angular.module('myApp.controllers').controller("licenseController", ['$scope', 'shellService', function($scope, shellService){
   
   function Init() {
       
   }
   
   $scope.Company = {
        Name: "Test Company",
        MaxUsers: 1,
        Expiry: new Date()
   }
   
   $scope.User = shellService.User;
  
   Init();
   
   
   
}]);
