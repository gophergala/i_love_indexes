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

  # Listen to input event and send search query
  searchInput.on "input", (e) ->
    value = e.target.value
    sendSearch value if value

  searchForm.on "submit", (e) ->
    e.preventDefault()

  addURL.on "submit", (e) ->
    addIndexOf urlInput.val()
    e.preventDefault()

  insertIntoTableBody = (indexItem) ->
    row = $("<tr>")
    fields = ['name', 'last_modified_at', 'size']
    fields.forEach (field, index) ->
      td = $("<td>")
      itemText = indexItem[field]
      if index == 1
        itemText = moment(itemText).fromNow()
      td.text itemText
      row.append td

    tableBody.append row

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
          dataType: 'json'
          data: {search: query}
          success: (data) ->
            console.log data
            tableBody.empty()
            data.forEach insertIntoTableBody if data instanceof Array
            header.animate
              'margin-top': 0
              'slow'
            table.fadeIn()

        @timeoutHandle = 0
      , 300
  )()

  addIndexOf = (url) ->
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
