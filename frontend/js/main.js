

//All this disgusting crap needs ripping out and implementing angularJS
function bindHandlers() {
	//Menu Nav
	$("#tabMap,#tabProfile, #tabSupport, #tabReports, #tabSettings,#tabLicense").click(function(){
        var Self = $(this);
        var Main = $("#Maintab");
        var actions = {
            tabMap: function() {
                 $.ajax({
                    type: "GET",
                    url: "/system/map",
                    dataType: "html",
                    success: function(HTML) {
                    	Main.html(HTML);
                    	var MapAPI = System.getMapAPI(); //probably need a refresh method on the MapAPI
                    	MapAPI.SetAPI("GoogleMaps");
                    },
                    error: function(a,b,c) {
                      System.showLostConnection();
                    }
                })
            },
            tabSupport: function() {
                  $.ajax({
                    type: "GET",
                    url: "/system/support",
                    dataType: "html",
                    success: function(HTML) {
                    	Main.html(HTML);
                    	tinymce.init({
					    selector: "textarea",
					    plugins: [
					        "advlist autolink lists link image charmap print preview anchor",
					        "searchreplace visualblocks code fullscreen",
					        "insertdatetime media table contextmenu paste moxiemanager"
					    ],
					    height: 300,
					    toolbar: "insertfile undo redo | styleselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | link image"
			});
                    },
                    error: function(a,b,c) {
                      System.showLostConnection();
                    }
                })
            },
            tabReports: function() {
                $.ajax({
                    type: "GET",
                    url: "/system/report",
                    dataType: "json",
                    success: function(data) {	
                    	Main.html(data.HTML);
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
                    	
                    	
                    },
                    error: function(a,b,c) {
                      System.showLostConnection();
                    }
                })
            },
            tabSettings: function() {
                $.ajax({
                    type: "GET",
                    url: "/system/settings",
                    dataType: "html",
                    success: function(HTML) {
                    	Main.html(HTML);
                    },
                    error: function(a,b,c) {
                      System.showLostConnection();
                    }
                })
            },
            tabLicense: function() {
            	$.ajax({
            		type: "GET",
            		url: "/system/license",
            		dataType: "html",
            		success: function(HTML) {
            			Main.html(HTML);
            		}
            	});
            }
        }[Self.attr("id")]();
       
	});
}
