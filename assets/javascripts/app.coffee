$ ->
  searchInput = $('input[type=search]')
  tableBody = $('table tbody')
  urlInput = $('input[type=url]')
  addButton = $('button')

  console.log addButton

  # Listen to input event and send search query
  searchInput.on "input", (e) ->
    value = e.target.value
    sendSearch value

  addButton.on "click", (e) ->
    addIndex urlInput.val()

  insertIntoTableBody = (data) ->
    row = $("<tr>")
    fields = ['Name', 'Last modified', 'Size', 'Description']
    fields.forEach (field) ->
      td = $("<td>")
      td.text data[field]
      row.appendChild td

    tableBody.appendChild row

  # Send search query with a delay
  sendSearch = (() ->
    @timeoutHandle = 0

    (query) =>
      # Cancel waiting search query
      if timeoutHandle != 0
        clearTimeout @timeoutHandle
        timeoutHandle = 0

      # Setup the new search query
      @timeoutHandle = setTimeout () ->
        $.ajax
          type: "GET"
          url: '/api/search'
          data: {search: query}
          success: (data) ->
            data.forEach insertIntoTableBody if data instanceof Array

        @timeoutHandle = 0
        console.log 'SENT !'
      , 300
  )()

  addIndex = (url) ->
    $.ajax
      type: "POST"
      contentType: "application/json; charset=utf-8"
      dataType: 'json'
      url: '/api/indices'
      data: JSON.stringify url: url
      success: (data) ->
        # Do something
