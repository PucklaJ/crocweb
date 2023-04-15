let server_host

function CodeRequest(shared_secret) {
  const url = `http://${server_host}/code/${shared_secret}`
  return fetch(url)
}

function RecieveRequest(id, index) {
  const url = `http://${server_host}/receive/${id}/${index}`
  return fetch(url)
}

function DownloadRequest(id, index) {
  const url = `http://${server_host}/download/${id}/${index}`
  return fetch(url)
}

function LoadFileIntoHolder(file_name, id, index, receive_holder, receive_response) {
  console.log(`Receiving ${file_name} ...`)

  const file_div = document.createElement("div")
  const file_nm = document.createElement("p")
  const file_download = document.createElement("button")
  file_nm.innerHTML = file_name
  file_download.innerHTML = "Download"
  file_download.onclick = async function() {
    const element = document.createElement("a")
    element.style.display = "none"
    element.setAttribute("download", file_name)

    const download_response = await DownloadRequest(id, index)
    if (download_response.ok) {
      const download_blob = await download_response.blob()

      if (window.navigator && window.navigator.msSaveOrOpenBlob) {
        window.navigator.msSaveOrOpenBlob(download_blob)
      } else {
        const data = window.URL.createObjectURL(download_blob)
        const link = document.createElement("a")
        link.href = data
        link.download = file_name
        link.dispatchEvent(
          new MouseEvent('click', {
            bubbles: true,
            cancelable: true,
            view: window
          })
        )

        setTimeout(() => {
          window.URL.revokeObjectURL(data)
          link.remove()
        }, 100)
      }
    } else {
      const error_msg = await download_response.text()
      console.error(`Failed to download ${file_name}: ${download_response.status}/${error_msg}`)
    }
  }

  const file_placeholder = document.createElement("div")
  file_placeholder.innerHTML = "<p>Loading ...</p>"

  file_div.appendChild(file_placeholder)
  file_div.appendChild(file_nm)
  file_div.appendChild(file_download)
  receive_holder.appendChild(file_div)

  receive_response.then(async recv_resp => {
    console.log(`Got Response from ${file_name}`)
    const recv_resp_text = await recv_resp.text()
    if (recv_resp.ok) {
      file_placeholder.innerHTML = recv_resp_text
    } else {
      file_placeholder.innerHTML = `<p>Response Error: ${recv_resp.status}/${recv_resp_text}</p>`
    }
  })
  receive_response.catch(error => {
    file_placeholder.innerHTML = `<p>Promise Error: ${error}</p>`
  })
}