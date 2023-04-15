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

function LoadFileIntoHolder(file_name, custom, id, index, receive_holder, receive_response, todo_list) {
  const file_div = document.createElement("div")
  file_div.className = "file_div"
  const file_nm = document.createElement("p")
  const file_download = document.createElement("button")
  file_nm.innerHTML = file_name
  file_nm.className = "file_name"
  file_download.className = "file_download_button"
  file_download.innerHTML = "Download"
  file_download.onclick = function() {
    const link = document.createElement("a")
    link.href = `http://${server_host}/download/${id}/${index}`
    link.download = file_name.replace(/^.*[\\\/]/, '')
    link.click()
  }

  const file_placeholder = document.createElement("div")
  const styleWidth = 750
  file_placeholder.innerHTML = '<p style="text-align:center;">Loading ...</p>'
  if (custom != null) {
    let width, height
    if (custom.width !== undefined) {
      width = custom.width
    }
    if (custom.height !== undefined) {
      height = custom.height
    }

    const styleHeight = height / width * styleWidth

    file_placeholder.style.width = styleWidth
    file_placeholder.style.height = styleHeight
    file_placeholder.style.backgroundColor = "#a83232"
  }

  file_div.appendChild(file_placeholder)
  file_div.appendChild(document.createElement("div"))
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
      file_placeholder.style.width = "auto"
      file_placeholder.style.height = "auto"
    })
    .catch(error => {
      file_placeholder.innerHTML = `<p>Promise Error: ${error}</p>`
    })
    .finally(() => {
      todo_list.splice(todo_list.indexOf(index), 1)
      if (todo_list.length == 0) {
        request_button.disabled = false
        clear_button.disabled = false
      }
    })
  })
  receive_response.catch(error => {
    todo_list.splice(todo_list.indexOf(index), 1)
    if (todo_list.length == 0) {
      request_button.disabled = false
      clear_button.disabled = false
    }
    file_placeholder.innerHTML = `<p>Promise Error: ${error}</p>`
  })
}