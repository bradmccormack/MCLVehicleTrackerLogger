(function() {
	var dateFrom, dateTo;

	var datepickerFrom = $('#routeDateFrom');
	datepickerFrom.datetimepicker({
		language : 'pt-BR',
		pick12HourFormat : true,
		format : 'dd/MM/yyyy hh:mm:ss'
	});

	var datepickerTo = $('#routeDateTo');
	datepickerTo.datetimepicker({
		language : 'pt-BR',
		pick12HourFormat : true,
		format : 'dd/MM/yyyy hh:mm:ss'
	});

	datepickerFrom.on('changeDate', function(e) {
		//TODO the UTC date in the date picker appears to be wrong
		dateFrom = e.date.toString();
		//use UTC as the logger records everything
		//localDate can also be used
	});

	datepickerTo.on('changeDate', function(e) {
		//TODO the UTC date in the date picker appears to be wrong
		dateTo = e.date.toString();

	});

	var showRoute = $("#showRoute");
	//Routes Buttons
	$("#showRoute").click(function() {

		var from = $("div#routeDateFrom > input").val();
		var to = $("div#routeDateTo > input").val();

		if (!from || !to || !dateFrom || !dateTo) {
			//invalid dates
			alert("dates invalid - TODO bootstrap validation around inputs");
			return;
		}

		if (dateFrom > dateTo) {
			alert("To date must be > from date")
		}

		$.ajax({
			type : "POST",
			url : "/system/historicalroute",
			datatype : "json",
			data : {
				dateFrom : from,
				dateTo : to
			},
			success : function(JSON) {
				if (JSON.success) {
					var vehicles = JSON.data;
					var vl = Object.keys(vehicles).length;
					if(vl == 0) {
						alert("no vehicle data for that time period");
						return;
					}
						
					/*
					 *   addtoRoute: function(Route, Point, Color) {
	        	["Latitude", "Longitude", "Speed", "Heading", "Fix", "DateTime"].forEach(function() {
					 */

					debugger;
					while(true) {
						for(i = 0; i < vl; i++) {
							var current = Object.keys(vehicles)[i];
							
							if(current.length > 0) {
								var point = vehicles[current].shift();
								//Lat, Long, Speed, Fix, Heading, Date
								System.getMapAPI().Current().addtoRoute(current, 
									{Latitude: point[0], Longitude: point[1], Speed: point[2], Fix: point[3], Heading: point[4], DateTime: point[5]});
							} else {
								//done for this vehicle remove it from the list
								break;//hack job as I know there is only one vehicle in the system now
							}
							
						}
						
					}
					//Note will need to update the Legend too
					
				
					//mapAPI.Current().setMarker(data.Latitude, data.Longitude,"", Vehicles[data.ID].Color);
				}

			},
			error : function(a, b, c) {
				System.showLostConnection();
			}
		})
	});

})(); 