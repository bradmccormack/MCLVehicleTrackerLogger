angular.module('myApp.controllers').controller("reportController", ['$scope', '$http', function($scope, $http){
   
   $scope.Init = function() {
       $http({method: 'GET', url: '/system/report'}).
           success(function(data, status, headers, config) {

               //Main.html(data.HTML);
               //add the data to the charts

               var barreport = $("#barreport")[0];

               //TODO think about moving more of the presentation to the server (eg the axis labels and colours)
               var barChartData = { labels : ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]}
               barChartData.datasets = [];
               barChartData.datasets.push({
                   fillColor : "rgba(151,187,205,0.5)",
                   strokeColor : "rgba(151,187,205,1)",
                   data : []
               });

               data.KMPerDay.forEach(function(value){
                   barChartData.datasets[0].data.push(value);
               });


               var piereport = $("#piereport")[0];
               var pieChartData = [
                   {
                       value : data.Availability[0],
                       color: "#009933"
                   },
                   {
                       value : data.Availability[1],
                       color: "#FF0000"
                   }
               ];



               var kmChart = new Chart(barreport.getContext("2d")).Bar(barChartData);
               var availabilityChart = new Chart(piereport.getContext("2d")).Pie(pieChartData);

           }).
           error(function(data, status, headers, config) {
               // called asynchronously if an error occurs
               // or server returns response with an error status.
           });
   }
   

   
   
   
}]);
