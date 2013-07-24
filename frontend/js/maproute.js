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
		System.getMapAPI().Current().clearRoutes();
		
		$.ajax({
			type : "POST",
			url : "/system/historicalroute",
			datatype : "json",
			data : {
				dateFrom : from,
				dateTo : to
			},
			success : function(result) {
				if (result.success) {
					var vehicles = result.data;
					var vl = Object.keys(vehicles).length;
					if(vl == 0) {
						alert("no vehicle data for that time period");
						return;
					}
						
					var Colors = {};
					for(vehicle in vehicles) {
						Colors[vehicle] = Utility.RandomColor();
						System.updateLegend(vehicle, Colors[vehicle]);
					}
				
					//debugger;
					while(true) {
						for(i = 0; i < vl; i++) {
							var currentvehicle = Object.keys(vehicles)[i];
							var currentpositions = vehicles[currentvehicle];
							
							if(currentpositions.length > 0) {
								var point = vehicles[currentvehicle].shift();
								//Lat, Long, Speed, Fix, Heading, Date
								System.getMapAPI().Current().addtoRoute(currentvehicle, 
									{Latitude: point[0], Longitude: point[1], Speed: point[2], Fix: point[3], Heading: point[4], DateTime: point[5]}, Colors[currentvehicle]);
							} else {
								delete vehicles[currentvehicle];
								break;
							}
						}
						var vl = Object.keys(vehicles).length;
						if(vl == 0) {
							break;
						}
						
					}
					//Note will need to update the Legend too
				}

			},
			error : function(a, b, c) {
				System.showLostConnection();
			}
		})
	});

})(); 