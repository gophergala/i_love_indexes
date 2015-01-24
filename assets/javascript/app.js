(function () {
	var searchInput = document.querySelector('input[type=search]')
	var tableBody = document.querySelector('table tbody')

	var urlInput = document.querySelector('input[type=url]')
	var addButton = document.querySelector('button')

	console.log(addButton)

	// Listen to input event and send search query
	searchInput.addEventListener('input', function (e) {
		var value = e.target.value

		sendSearch(value)
	})

	addButton.addEventListener('click', function (e) {
		console.log('CLICKED')
	})

	function insertIntoTableBody (data) {
		var row = document.createElement('tr')

		var fields = ['Name', 'Last modified', 'Size', 'Description']
		fields.forEach(function (field) {
			var td = document.createElement('td')
			td.innerText = data[field]

			row.appendChild(td)
		})

		tableBody.appendChild(row)
	}

	// Send search query with a delay
	var sendSearch = (function () {
		var timeoutHandle = 0

		return function (query) {
			// Cancel waiting search query
			if (timeoutHandle != 0) {
				clearTimeout(timeoutHandle)
				timeoutHandle = 0
			}

			// Setup the new search query
			timeoutHandle = setTimeout(function () {
				$.get('/api/search', {search: query}, function (data) {
					if (data instanceof Array) {
						data.forEach(insertIntoTableBody)
					}
				})

				timeoutHandle = 0
				console.log('SENT !')
			}, 300)
		}
	})()

	function addIndex (url) {
		$.post('/api/indices', {url: url}, function (data) {
			// body...
		})
	}

	addIndex('http://itinuae.com/torrents/movies/')

})()