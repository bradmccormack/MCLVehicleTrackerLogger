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

					while(true) {
						
					}
					
					/*
					 * until all done
					 * for each vehicle pop off the next co-ordinate and add it to the polygon
					 * 		
					 */

					/*
					var flightPlanCoordinates = [new google.maps.LatLng(37.772323, -122.214897), new google.maps.LatLng(21.291982, -157.821856), new google.maps.LatLng(-18.142599, 178.431), new google.maps.LatLng(-27.46758, 153.027892)];
					var flightPath = new google.maps.Polyline({
						path : flightPlanCoordinates,
						strokeColor : '#FF0000',
						strokeOpacity : 1.0,
						strokeWeight : 2
					});
					*/

					flightPath.setMap(map);

					//mapAPI.Current().setMarker(data.Latitude, data.Longitude,"", Vehicles[data.ID].Color);
				}

			},
			error : function(a, b, c) {
				System.showLostConnection();
			}
		})
	});

})(); 