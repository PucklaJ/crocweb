package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gopherjs/gopherjs/js"
)

var Client struct {
	ReceiveData Receive
}

type ReceiveFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Receive struct {
	ID    uint64        `json:"id"`
	Files []ReceiveFile `json:"files"`
}

func CodeRequest(sharedSecret string) (Receive, error) {
	url := fmt.Sprint("http://localhost:8080/code/", sharedSecret)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Receive{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Receive{}, err
	}
	defer res.Body.Close()

	var recv Receive

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&recv)
	if err != nil {
		return recv, err
	}

	return recv, nil
}

func ReceiveRequest(id uint64) (string, error) {
	url := fmt.Sprint("http://localhost:8080/receive/", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	contents, err := io.ReadAll(res.Body)

	return string(contents), err
}

func (recv Receive) IntoHolder(receiveHolder *js.Object) {
	receiveHolder.Set("innerHTML", nil)
	document := js.Global.Get("document")

	for _, f := range recv.Files {
		fileName := document.Call("createElement", "p")
		fileName.Set("innerHTML", f.Name)
		receiveHolder.Call("appendChild", fileName)
	}
}
