let server_host

function CodeRequest(shared_secret) {
  const url = `http://${server_host}/code/${shared_secret}`
  return fetch(url)
}

function ReceiveRequest(id, index) {
  const url = `http://${server_host}/receive/${id}/${index}`
  return fetch(url)
}

function DownloadRequest(id, index) {
  const url = `http://${server_host}/download/${id}/${index}`
  return fetch(url)
}

function LoadFileIntoHolder(file_name, id, index, receive_holder, receive_response, todo_list) {
  console.log(`Receiving ${file_name} ...`)

  const file_div = document.createElement("div")
  const file_nm = document.createElement("p")
  const file_download = document.createElement("button")
  file_nm.innerHTML = file_name
  file_download.innerHTML = "Download"
  file_download.onclick = function() {
    const link = document.createElement("a")
    link.href = `http://${server_host}/download/${id}/${index}`
    link.download = file_name
    link.click()
    setTimeout(link.remove, 100)
  }

  const file_placeholder = document.createElement("div")
  file_placeholder.innerHTML = "<p>Loading ...</p>"

  file_div.appendChild(file_placeholder)
  file_div.appendChild(file_nm)
  file_div.appendChild(file_download)
  receive_holder.appendChild(file_div)

  receive_response.then(recv_resp => {
    recv_resp.text()
    .then(recv_resp_text => {
      if (recv_resp.ok) {
        file_placeholder.innerHTML = recv_resp_text
      } else {
        file_placeholder.innerHTML = `<p>Response Error: ${recv_resp.status}/${recv_resp_text}</p>`
      }
    })
    .catch(error => {
      file_placeholder.innerHTML = `<p>Promise Error: ${error}</p>`
    })
    .finally(() => {
      todo_list.splice(todo_list.indexOf(index), 1)
      if (todo_list.length == 0) {
        request_button.disabled = false
      }
    })
  })
  receive_response.catch(error => {
    todo_list.splice(todo_list.indexOf(index), 1)
    if (todo_list.length == 0) {
      request_button.disabled = false
    }
    file_placeholder.innerHTML = `<p>Promise Error: ${error}</p>`
  })
}