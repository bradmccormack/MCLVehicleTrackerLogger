angular.module('myApp.controllers').controller("navController", ['$scope', '$location', 'shellService', function($scope, $location, shellService){

    $scope.$on("ConfigChanged", function (Event, Data) {
        $scope.IsLogged = Data.User != "" && $.cookie("data") != undefined;
    });

    //For setting the active attribute on the navbar tabs for the main menu
	$scope.isActive = function(viewLocation) {
		var active =  (viewLocation == $location.path());
		var locationpath = $location.path();

		return active;
	}
}]);
