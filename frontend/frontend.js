function CodeRequest(shared_secret) {
  const url = `http://localhost:8080/code/${shared_secret}`
  return fetch(url)
}

function RecieveRequest(id) {
  const url = `http://localhost:8080/receive/${id}`
  return fetch(url)
}