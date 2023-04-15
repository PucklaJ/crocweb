let code_text, receive_holder, request_button, clear_button

function on_init() {
  code_text = document.getElementById("code_text")
  receive_holder = document.getElementById("receive_holder")
  request_button = document.getElementById("request_button")
  clear_button = document.getElementById("clear_button")
  server_host = document.location.host
}

function on_request() {
  request_button.disabled = true
  clear_button.disabled = true

  shared_secret = code_text.value

  receive_holder.innerHTML = `<p>Requesting Code "${shared_secret}" ...</p>`

  CodeRequest(shared_secret)
  .then(code_response => {
    if (code_response.ok) {
      code_response.json()
      .then(receive_data => {
        receive_holder.innerHTML = null
        let todo_list = Array.from(Array(receive_data.files.length).keys())
        for (let i = 0; i < receive_data.files.length; i++) {
          const recv_file = receive_data.files[i].name;
          const receive_response = ReceiveRequest(receive_data.id, i)

          LoadFileIntoHolder(
            recv_file,
            receive_data.id,
            i,
            receive_holder,
            receive_response,
            todo_list
          )
        }
      })
      .catch(error => {
        request_button.disabled = false
        receive_holder.innerHTML = `<p>CodeRequest Failed to parse JSON: ${error}</p>`
      })
    } else {
      request_button.disabled = false
      code_response.text()
      .then(error_msg => {
        receive_holder.innerHTML = `<p>CodeRequest Response Error: ${code_response.status}/${error_msg}</p>`
      })
      .catch(error => {
        receive_holder.innerHTML = `<p>CodeRequest Response Error Message Error: ${error}</p>`
      })
    }
  })
  .catch(error => {
    request_button.disabled = false
    receive_holder.innerHTML = `<p>CodeRequest Promise Error: ${error}</p>`
  })
}

function on_clear() {
  receive_holder.innerHTML = null
  clear_button.disabled = true
}

function on_code_text(event) {
  if (event.keyCode == 13) {
    on_request()
  }
}