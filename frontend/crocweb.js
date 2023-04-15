let code_text, receive_holder, request_button

function init() {
  code_text = document.getElementById("code_text")
  receive_holder = document.getElementById("receive_holder")
  request_button = document.getElementById("request_button")
  server_host = document.location.host
}

function request() {
  request_button.disabled = true

  shared_secret = code_text.value
  console.log(`Requesting Code "${shared_secret}" ...`)

  CodeRequest(shared_secret)
  .then(code_response => {
    if (code_response.ok) {
      code_response.json()
      .then(receive_data => {
        console.log(receive_data)

        receive_holder.innerHTML = null
        let todo_list = Array.from(Array(receive_data.files.length).keys())
        for (let i = 0; i < receive_data.files.length; i++) {
          const recv_file = receive_data.files[i].name;
          const receive_response = ReceiveRequest(receive_data.id, i)
          LoadFileIntoHolder(recv_file, receive_data.id, i, receive_holder, receive_response, todo_list)
        }
      })
      .catch(error => {
        console.error(`CodeRequest Failed to parse JSON: ${error}`)
      })
    } else {
      code_response.text()
      .then(error_msg => {
        console.error(`CodeRequest Response Error: ${code_response.status}/${error_msg}`)
      })
      .catch(error => {
        console.error(`CodeRequest Response Error Message Error: ${error}`)
      })
    }
  })
  .catch(error => {
    console.error(`CodeRequest Promise Error: ${error}`)
  })
}
