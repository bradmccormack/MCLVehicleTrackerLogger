angular.module('myApp.controllers').controller("navController", ['$scope', '$location', '$cookies', function($scope, $location, $cookies){

    $scope.$on("ConfigChanged", function (Event, Data) {
        $scope.IsLogged = (Data.User != undefined && Data.User != "" && ("data" in $.cookie()))
    });

    //For setting the active attribute on the navbar tabs for the main menu
	$scope.isActive = function(viewLocation) {
		var active =  (viewLocation == $location.path());
		var locationpath = $location.path();

		return active;
	}
}]);
