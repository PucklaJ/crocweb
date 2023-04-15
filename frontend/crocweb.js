let code_text, recieve_holder, request_button

function init() {
  code_text = document.getElementById("code_text")
  recieve_holder = document.getElementById("receive_holder")
  request_button = document.getElementById("request_button")

  console.log("Frontend Loaded")
}

async function code_text_enter() {
  request_button.disabled = true

  shared_secret = code_text.value
  console.log(`Requesting Code "${shared_secret}" ...`)

  const code_response = await CodeRequest(shared_secret)
  if (code_response.ok) {
    const recieve_data = await code_response.json()
    console.log(recieve_data)

    const recieve_response = await RecieveRequest(recieve_data.id)
    if (recieve_response.ok) {
      const recieve_element = await recieve_response.text()

      recieve_holder.innerHTML = recieve_element
    } else {
      const error_msg = await recieve_response.text()
      console.error(`Failed Request: ${recieve_response.status}/${error_msg}`)
    }
  } else {
    const error_msg = await code_response.text()
    console.error(`Failed Code: ${code_response.status}/${error_msg}`)
  }

  request_button.disabled = false
}