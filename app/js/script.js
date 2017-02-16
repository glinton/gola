function httpGetAsync(uriString, callback)
{
	var xmlHttp = new XMLHttpRequest()
	xmlHttp.onload = function() {
		callback(xmlHttp.responseText)
	}

	xmlHttp.open("GET", uriString, true) // true for asynchronous
	xmlHttp.send(null)
}
