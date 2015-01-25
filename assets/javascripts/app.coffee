$ ->
  searchInput = $('input[type=search]')
  searchForm = $("#search-form")
  resultsTableBody = $('#results table tbody')
  indicesTableBody = $('#indices table tbody')
  urlInput = $('input[type=url]')
  addURLForm = $('#add-url-form')
  listURLsButton = $('#list-url')

  # Hide table
  indicesTable = $('#indices')
  resultsTable = $('#results')

  # Vertically center header
  header = $('#header')
  header.css('margin-top', $(window).height() / 2 - header.height() / 1.5)

  # Listen to input event and send search query
  searchInput.on "input", (e) ->
    value = e.target.value
    sendSearch value if value

  searchForm.on "submit", (e) ->
    e.preventDefault()

  addURLForm.on "submit", (e) ->
    addIndexOf urlInput.val() if urlInput.val()
    e.preventDefault()

  insertIntoTableBody = (tbody, indexItem, fields) ->
    row = $("<tr>")
    fields.forEach (field) ->
      td = $("<td>")
      if field == "last_modified_at"
        td.html moment(indexItem[field]).fromNow()
      else if field == "url"
        item = $("<a>").attr("href", indexItem[field]).text(indexItem[field])
        td.append item
      else if field == "name"
        item = $("<a>").attr("href", indexItem["url"]).text(indexItem[field])
        td.append item
      else
        td.html indexItem[field]
      row.append td
    tbody.append row

  listURLsButton.on "click", (e) ->
    $.ajax
      type: "GET"
      url: '/api/indices'
      dataType: 'json'
      success: (data) ->
        indicesTableBody.empty()
        data.forEach (elem) ->
          url = elem.scheme + "://" + elem.host + elem.path
          insertIntoTableBody indicesTableBody, {url: url}, ["url"]
        header.animate
          'margin-top': 0
          'slow'
        resultsTable.fadeOut ->
          indicesTable.fadeIn()

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
            resultsTableBody.empty()
            data.forEach (elem) ->
              insertIntoTableBody resultsTableBody, elem, ['name', 'last_modified_at', 'size']
            header.animate
              'margin-top': 0
              'slow'
            indicesTable.fadeOut ->
              resultsTable.fadeIn()


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
        $("#btn-row .error").addClass("hidden")
        $("#btn-row .error").val("")
        urlInput.val("")
      error: (response, error) ->
        errSpan = $("#btn-row .error")
        setTimeout () ->
          errSpan.empty()
        , 3000
        if response.status == 500
          errSpan.html("Internal Error, sorry :-/")
        else if response.status == 422
          errors = response.responseJSON.errors
          out = "<pre>"
          Object.keys(errors).forEach (key) ->
            errors[key].forEach (err) ->
              out += key + ": " + err + "\n"
          errSpan.html(out + "</pre>")
