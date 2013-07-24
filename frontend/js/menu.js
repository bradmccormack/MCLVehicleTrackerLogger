debugger;
//Map controls
$("div#Mainmapcontrols button").click(function() {
	var Self = $(this);
	
	function notAvailable() {
		alert("Not available");
	}
	

	var actions = {
	    "mapRefresh" : function() {
	    	notAvailable();
	    },
	    "mapZoomIn" : function() {
	    	System.getMapAPI().Current().zoomIn();
	        //System.getMapAPI().zoomIn();
	    },
	    "mapZoomOut" : function() {
	        System.getMapAPI().Current().zoomOut();
	        //System.getMapAPI().zoomOut();
	    },
	    "mapMarker" : function() {
	        notAvailable();
	    },
	    "mapFollow" : function() {
	        notAvailable();
	    },
	    "mapRoute" : function() {
	        notAvailable();
	    },
	    "mapPrint" : function() {
	        notAvailable();
	    }
	}[Self.attr("id")]();
});









