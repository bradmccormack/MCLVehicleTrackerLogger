
var Utility = (function(){
	return {
		RandomColor: function() {
			return (function lol(m,s,c){return s[m.floor(m.random() * s.length)] +
  			(c && lol(m,s,c-1));})(Math,'0123456789ABCDEF',4)
		}
	}	
})();


var System = (function(){

	var Self = this;
	
    var Con;
	var Colours = [];
    var Vehicles = [];

    //System wide settings
    var mapAPI = {};
    var Settings = { Marker: { InterpolateCount : 10}};
    var Position = { Last: {}};

	//Dom elements that we will write to 
	var tabVehicles = $("#tab2");


	function showLostConnection() {
		$('#systemError').modal('toggle');  
	}

	function systemMessage(Message) {
		$("div#SystemMessages > ul#Messages").append("<li class='text-info'>" + new Date().toTimeString() + Message + "</li>");
	}

    function updateLegend(VehicleID) {
  
        //We need to know if the GPS signal is correct or not (Fix status is true)
      	
        var Legend = $("div#Mainlegend div#Vehiclelegend");
        Legend.find("span.text-error").remove();
      
        Vehicles[VehicleID] = {
        	DateTime: new Date(),
        	Color: Utility.RandomColor()
        }
        Legend.find("ul").append('<li><a href="#" style="color: #' + Vehicles[VehicleID].Color + '"><i class="icon-truck"></i> ' + VehicleID +'</a></li>');
          //mapAPI.Current().setMarkerColor(Vehicles[VehicleID].Color);
        //remove everything from the legend if there has been no contact in over 1 hour or whatever
        
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
        	//TODO Add the spinner to the Login button
        	
        	var cookies = $.cookie();
       		if("session" in $.cookie())
       		{
       			if("success" in cbobj && typeof cbobj.success == "function") {
       				cbobj.success("Craig Smith"); //TODO pull the name from the session cookie
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
       				//bindHandlers();
       				
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
       							var logSuccess = ("success" in JSON) && JSON.success && "session" in $.cookie()
	       						if(logSuccess) 
	       						{
       								if("success" in cbobj && typeof cbobj.success == "function")
        									cbobj.success(JSON.user); //callback with logged in user name
        							return;
	       						}
	       						else if("retries" in JSON) 
	       						{
	       							var retries = parseInt(JSON.retries);
	       							if(retries == 0) 
	       							{
	       								if("fail" in cbobj && typeof cbobj.fail == "function")
	       									cbobj.fail(JSON.message);
	       							}
	       							else
	       							{
	       								
	       								if("retry" in cbobj && typeof cbobj.retry == "function")
	       									cbobj.retry(JSON.message);
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
			//mapAPI.SetAPI("GoogleMaps");
			
            if (window["WebSocket"]) {
                    //Con = new WebSocket("ws://dev.myclublink.com.au/ws");
                    Con = new WebSocket("ws://dev.myclublink.com.au:8080/ws");
                 
                    Con.onopen = function() {
                    	systemMessage("Connected to server");
                    };

                    Con.onclose = function(evt) {
                        systemMessage("Server connection closed");
                    }
                    Con.onmessage = function(evt) {
                    	var data = JSON.parse(evt.data).Entry;
                      
                      	//add vehicle to Legend if it is not there   
                        if(!(data.ID in Vehicles)) {
                        	updateLegend(data.ID);
                        }
                        //TODO remove vehicle if no contact for X minutes
                        
                        mapAPI.Current().setMarker(data.Latitude, data.Longitude,"", Vehicles[data.ID].Color);
                      	
						$(tabVehicles).find("span.text-error").remove();
                      	var VehicleInfo = $(tabVehicles).find("span[data-vehicle = '" + data.ID + "']");
                      	VehicleInfo.remove();

                      	var html = "<span data-vehicle='" + data.ID + "'> <i class='icon-truck'></i> " + data.ID + "  <strong>Speed(KM/Hr)</strong> " + data.Speed 
                      	+ " <strong>Heading Degrees)</strong> " + Math.round(data.Heading) + " <strong>Time</strong> " + data.Date + "</span>"
                      	$(tabVehicles).append(html);
                                     
                        /*
                       


                        mapAPI.Current().setMarker(cords[0],cords[1]);
                        */
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
   
   //TODO set the datepicker from to localtime - 1 day and to to localtime current date
   
  
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
                    	var MapAPI = System.getMapAPI(); //probably need a refresh method on the MapAPI
                    	MapAPI.SetAPI("GoogleMaps");
                    	bindHandlers();
                    
                    
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

$(document).ready(function() {

    //Startup the main system
    //
    //attempt to login
	bindHandlers();
    System.login(
    	{ 
    		success: function(user) {
    			  $("#myModal").modal("toggle");
    			  //update the user drop down
    			  $("span#username").text(user);
    			  System.init();

    			  //display success message or something
    		},
    		fail: function(msg) {
    			var Loginmodal = $("#myModal div.modal-footer");
    			$(Loginmodal).find("strong.text-error").remove();
    		$("<strong class='text-error pull-left'> <i class='icon-warning-sign icon-white'></i> " + msg + "</strong>").insertBefore($(Loginmodal).find("strong"));
    		},
    		retry: function(msg) {
    			var Loginmodal = $("#myModal div.modal-footer");
    			$(Loginmodal).find("strong.text-warning").remove();
    		$("<strong class='text-warning pull-left'> <i class='icon-warning-sign icon-white'></i> " + msg + "</strong>").insertBefore($(Loginmodal).find("strong"));
    		}
	});
    	
    
    /*
	mapAPI.onClick(function(e){
		alert("Clicked at " + e.Location);
	})
	*/

});
