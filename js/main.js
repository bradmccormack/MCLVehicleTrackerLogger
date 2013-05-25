$(document).ready(function() {
	var map = L.map('Mainmap').setView([51.505, -0.09], 13);


	// add an OpenStreetMap tile layer
	L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png').addTo(map);


	// add a marker in the given location, attach some popup content to it and open the popup
	L.marker([51.5, -0.09]).addTo(map) 
	    .bindPopup('A pretty CSS3 popup. <br> Easily customizable.')
	    .openPopup();
	    
	
});
