(function() {
	var dateFrom, dateTo;
	
	var datepickerFrom = $('#routeDateFrom');
	   datepickerFrom.datetimepicker({
	   	language: 'pt-BR',
	   	pick12HourFormat: true,
	   	format: 'dd/MM/yyyy hh:mm:ss'
	   });
   
    var datepickerTo = $('#routeDateTo');
	   datepickerTo.datetimepicker({
	   	language: 'pt-BR',
	   	pick12HourFormat: true,
	   	format: 'dd/MM/yyyy hh:mm:ss'
	   });
  
  	datepickerFrom.on('changeDate', function(e) {
  		dateFrom = e.date.toString(); //use UTC as the logger records everything 
  		//localDate can also be used
	});
	
	datepickerTo.on('changeDate', function(e) {
		dateTo = e.date.toString();
	});
	
	var showRoute = $("#showRoute");
	//Routes Buttons
	$("#showRoute").click(function() {
		
		var from = $("div#routeDateFrom > input").val();
		var to = $("div#routeDateTo > input").val();
	
		if(!from || !to || !dateFrom || !dateTo) {
			//invalid dates
			alert("dates invalid - TODO bootstrap validation around inputs");
			return;
		}

		if(dateFrom > dateTo ) {
			alert("To date must be > from date")	
		}
		

		$.ajax({
			type: "POST",
			url: "/system/historicalroute",
			datatype: "json",
			data: {dateFrom: from, dateTo: to},
			success: function(JSON) {
				/*
				 * [{Vehicle: "ISUSU", GPSSentences: ["sentence1", "sentence2"}]
				 */
			},
			error: function(a, b, c) {
					System.showLostConnection();
			}
		})
	});
	
	
})();