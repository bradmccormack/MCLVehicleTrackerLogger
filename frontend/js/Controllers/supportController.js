angular.module('myApp.controllers').controller("supportController", ['$scope', function($scope){
   
   function Init() {
        $("div#editor").wysiwyg();
   }
   
   $scope.User =
   {
        First : "Brad",
        Last: "Mccormack"
   }
   
   $scope.Company = "Test Company";
   $scope.Send = function() {
 
   }
   
   Init();
   
   
   
}]);
