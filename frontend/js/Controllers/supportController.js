angular.module('myApp.controllers').controller("supportController", ['$scope', 'shellService', function($scope, shellService){
   
   $scope.Init = function(){
	   $('#editor').wysiwyg({
		   hotKeys: {
			   'ctrl+b meta+b': 'bold',
			   'ctrl+i meta+i': 'italic',
			   'ctrl+u meta+u': 'underline',
			   'ctrl+z meta+z': 'undo',
			   'ctrl+y meta+y meta+shift+z': 'redo'
		   }
	   });
	   $scope.Subject="";
	   $scope.Body = "";

	   $scope.User = {
		   First: shellService.User.First,
		   Last: shellService.User.Last
	   };

	   $scope.Company = shellService.Company.Name;
   };



   $scope.Send = function() {

   }

   


   
}]);
