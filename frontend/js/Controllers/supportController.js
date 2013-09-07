angular.module('myApp.controllers').controller("supportController", ['$scope', function($scope){
   
   $scope.Init = function(){
        $("div#editor").wysihtml5();
   }
   
   $scope.User =
   {
        First : "Brad",
        Last: "Mccormack"
   }
   
   $scope.Company = "Test Company";
   $scope.Send = function() {
 
   }
   


   
}]);
