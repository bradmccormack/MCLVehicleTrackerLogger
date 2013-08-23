


var System = (function(){
	

	var Self = this;
	
	//When SnapCount reaches SnapTrigger it will snap the view. Setting too low a value impacts performance quite a bit.
	var Camera = { Snap: true, SnapCount: 0, SnapTrigger: 10 };
	
	var Con;


	//System wide settings
	var mapAPI = {};
	
	var Position = { Last: {}};

	//Dom elements that we will write to 
	var tabVehicles = $("#tab2");


	

    function updateLegend(VehicleID, Color) {
  
        //We need to know if the GPS signal is correct or not (Fix status is true)
      	
        var Legend = $("div#Mainlegend div#Vehiclelegend");
        Legend.find("span.text-error").remove();
      
        Vehicles[VehicleID] = {
        	DateTime: new Date(),
        	Color: Color
        }
        Legend.find("ul li#vehicle_" + VehicleID).remove();
        Legend.find("ul").append('<li id=vehicle_' + VehicleID + '><a href="#" style="color: #' + Vehicles[VehicleID].Color + '"><i class="icon-truck"></i> ' + VehicleID +'</a></li>');
          //mapAPI.Current().setMarkerColor(Vehicles[VehicleID].Color);
        //remove everything from the legend if there has been no contact in over 1 hour or whatever
        
    }

	return {
		
		
		setMapAPI: function(mapAPI) {
			if(mapAPI) {
				mapAPI = mapAPI;
			}
		},
		setCameraSnap: function(enabled) {
			Camera.Snap = enabled;
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
