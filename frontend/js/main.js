


var System = (function(){

	var Self = this;
	
    var Con;
	var Colours = [];
    var Vehicles = [];

    //System wide settings
    var mapAPI = {};
    var Settings = { Marker: { InterpolateCount : 10}};
    var Position = { Last: {}};

	function showLostConnection() {
		$('#systemError').modal('toggle');  
	}

    function updateLegend(JSON) {
        //We need to know if the GPS signal is correct or not (Fix status is true)
        if(JSON.Vehicles.length > 0) {
            var Legend = $("div#Mainlegend div#Vehiclelegend ul");

            //Add any Vehicles to the Legend that are not there .. If GPS fix is not true then show warning etc
            //if there is no information for Vechile BLAH
            JSON.Vehicles.forEach(function(Vehicle){
                Legend.append('<li><a href="#"><i class="icon-truck"></i> ' + Vehicle +'</a></li>');
            });

        }
        else {
            //remove everything from the legend if there has been no contact in over 1 hour or whatever
        }
    }

	return {
		showLostConnection: function() {
			return this.showLostConnection;
		},
		setMapAPI: function(mapAPI) {
			if(mapAPI) {
				mapAPI = mapAPI;
			}
		},
        getMapAPI: function() {
          return mapAPI;
        },
        login: function(cbobj) {
        	var cookies = $.cookie();
       		if("session" in $.cookie())
       		{
       			if("success" in cbobj && typeof cbobj.success == "function") {
       				cbobj.success();
       			}
       			return;
       		}
       		
       		//get the login view
       		$.ajax(
       		{
       			type: "GET",
       			url: "/system/login",
       			error: function(jqXHR, textStatus, errorThrown) {showLostConnection();},
       			success: function(HTML) 
       			{
       				$('body').append(HTML);
       				bindHandlers();
       				
       				$('#myModal form.login-form').submit(function(e) 
       				{
       					
	       				e.preventDefault();
	       				var name = $(this).find("input[name=username]").val();
	       				var pass = $(this).find("input[name=password]").val();
	       				
	       				$.ajax({
	       					type: "POST",
	       					url: "/system/login",
	       					dataType: "JSON",
	       					xhrFields: {
     							withCredentials: true
   							},
	       					data: {name: name, password: pass},
	       					error: function(jqXHR, textStatus, errorThrown) {showLostConnection();},
       						success: function(JSON) 
       						{
       							var logSuccess = ("success" in JSON) && JSON.success;
       							//&& "Session" in $.cookie()
	       						if(logSuccess) 
	       						{
       								if("success" in cbobj && typeof cbobj.success == "function")
        									cbobj.success();
        							return;
	       						}
	       						else if("retries" in JSON) 
	       						{
	       							var retries = parseInt(JSON.retries);
	       							if(retries == 0) 
	       							{
	       								if("fail" in cbobj && typeof cbobj.fail == "function")
	       									cbobj.fail();
	       							}
	       							else
	       							{
	       								
	       								if("retry" in cbojb && typeof cbojb.retry == "function")
	       									cbobj.retry();
	       							}
	       						}
	       						
	       					}
       						
       					});
       					
       				});
       				$('#myModal').modal('toggle');  
       			}
       		});
        },
        
		init: function() {

			mapAPI = map; //set reference
			$("#tabMap").click();
			mapAPI.SetAPI("GoogleMaps");
			
            //add a couple of vehicles in hard coded for now
            //system.updateLegend({Vehicles: ["Mitsubishi Bus", "Izusu Bus"]});

            if (window["WebSocket"]) {
                    //alert("Browser supports Web Sockets. Yay");
                    //Con = new WebSocket("ws://dev.myclublink.com.au/ws");
                    Con = new WebSocket("ws://dev.myclublink.com.au:8080/ws");
                 
                    Con.onopen = function() {
                    	//log("Web Socket connection opened");
                        //Con.send("test message");
                    };

                    Con.onclose = function(evt) {
                        //alert("Closing web socket");
                        //appendLog($("<div><b>Connection closed.</b></div>"))
                    }
                    Con.onmessage = function(evt) {
                        var cords = evt.data.split(",");

                        mapAPI.Active.setMarker(cords[0],cords[1]);
                        /*
                        //TODO interpolate between cords
                        var X = cords[0];
                        var y = cords[1];
                        if(Position.Last.Latitude) {

                            var m = (Y - Position.Last.Longitude) / ( X - Position.Last.Latitude); //gradient
                            var b = 0; //figure out y intercept
                            var increment = X - Position.Last.Latitude / Settings.Marker.InerpolateCount;

                            for(var i = 0; i < Settings.Marker.InterpolateCount; i++) {
                                //y = mx+b
                                X += increment;
                                Y = m * X + b;
                                mapAPI.setMarker(X, Y);
                            }

                        }
                        else {
                            mapAPI.setMarker(cords[0],cords[1]); // remove this ?
                        }

                        Position.Last.Latitude = X;
                        Position.Last.Longitude = Y;
                        */

                        //alert("Message received " + evt.data);
                        //appendLog($("<div/>").text(evt.data))
                    }
            } else {
                alert("Your browser does not support WebSockets. You cannot use myClubLink until you upgrade to a modern browser");
            }
        }
	}
})();


//All this disgusting crap needs ripping out and implementing angularJS
function bindHandlers() {
    
    
   $("#settings").load(function() {
   		$(this).find("")
   }) 
    
    
    
   //systemError
   $("#systemError").submit(function(){
   	//redirect back to the home page for now
   	//window.location = "dev.myclublink.com.au";
  	window.location = "dev.myclublink.com.au";
   	
   })
    
    //Login
    $('#myModal').on("shown", function()
    {
        $(".username").focus(function() {
            $(".user-icon").css("left","-48px");
        });
        $(".username").blur(function() {
            $(".user-icon").css("left","0px");
        });

        $(".password").focus(function() {
            $(".pass-icon").css("left","-48px");
        });
        $(".password").blur(function() {
            $(".pass-icon").css("left","0px");
        });
    });


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
                    success: function(JSON) {	
                    	Main.html(JSON.HTML);
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
						
                    	JSON.KMPerDay.forEach(function(value){
                    		barChartData.datasets[0].data.push(value);
                    	});
                    	
                    	
                    	var piereport = $("#piereport")[0];
                    	var pieChartData = [
                    		{
                    			value : JSON.Availability[0],
                    			color: "#009933"
                    		},
                    		{
                    			value : JSON.Availability[1],
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

$(document).ready(function() {

    //Startup the main system
    bindHandlers();
    //attempt to login
    var Loginmodal = $("#myModal");
    System.login(
    	{ 
    		success: function() {
    			  $("#myModal").modal("toggle");
    			  System.init();

    			  //display success message or something
    		},
    		fail: function() {
    			modal.find("div.modal-body").html("<p>Sorry you do not have access to the system</p>");
    		},
    		retry: function() {
    			modalbody.append("<span class='text-error'> Incorrect details </span>");
    		}
	});
    	
    
    /*
	mapAPI.onClick(function(e){
		alert("Clicked at " + e.Location);
	})
	*/

});
