$ ->
  searchInput = $('input[type=search]')
  searchForm = $("#search-form")
  tableBody = $('table tbody')
  urlInput = $('input[type=url]')
  addURL = $('#add-url')

  # Hide table
  table = $('table')
  table.toggle()

  # Vertically center header
  header = $('#header')
  header.css('margin-top', $(window).height() / 2 - header.height() / 2)
  $('body').css('background-color', '#eee')

  # Listen to input event and send search query
  searchInput.on "input", (e) ->
    value = e.target.value
    sendSearch value

  searchForm.on "submit", (e) ->
    e.preventDefault()

  addURL.on "submit", (e) ->
    addIndex urlInput.val()
    e.preventDefault()

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
            header.animate
              'margin-top': 0
              'slow'
            $('body').animate 'background-color': 'white'
            table.slideDown()

        @timeoutHandle = 0
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
        $("#add-url .error").addClass("hidden")
        $("#add-url .error").val("")
      error: (response, error) ->
        errSpan = $("#add-url .error")
        if response.status == 500
          errSpan.html("Internal Error, sorry :-/")
        else if response.status == 422
          errors = response.responseJSON.errors
          out = "<ul>"
          Object.keys(errors).forEach (key) ->
            errors[key].forEach (err) ->
              out += "<li>" + key + " -> " + err + "</li>"
          out += "<ul>"
          errSpan.html(out)
