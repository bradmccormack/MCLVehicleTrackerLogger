angular.module('myApp.controllers').controller("licenseController", ['$scope', 'shellService', function($scope, shellService){
   
   function Init() {
       
   }
   
   $scope.Company = shellService.Company;
   $scope.User = shellService.User;
  
   Init();
   
   
   
}]);
