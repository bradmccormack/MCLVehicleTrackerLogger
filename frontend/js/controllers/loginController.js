 app.controller("loginControler", ['$scope', function($scope){
   
   $scope.login =  function(cbobj)
	{
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
        }
   
    /*
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
    */
   
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
    
}]);
   