package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/PucklaJ/crocweb/global"
	"github.com/schollz/croc/v9/src/croc"
)

func root(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(Server.RootDir, "index.html"))
}

func code(w http.ResponseWriter, r *http.Request) {
	sharedSecret := filepath.Base(r.URL.Path)

	receiver, err := croc.New(croc.Options{
		IsSender:       false,
		SharedSecret:   sharedSecret,
		Debug:          false,
		RelayAddress:   global.DefaultRelayAddress,
		RelayPorts:     global.DefaultRelayPorts,
		RelayPassword:  global.DefaultRelayPassword,
		Stdout:         false,
		NoPrompt:       true,
		DisableLocal:   false,
		NoMultiplexing: false,
		OnlyLocal:      false,
		NoCompress:     false,
		Curve:          global.DefaultCurve,
		HashAlgorithm:  global.DefaultHash,
		ThrottleUpload: global.DefaultUploadThrottle,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var recv *Receive

	doneChan := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Millisecond * 10)
		for {
			select {
			case <-ticker.C:
				if receiver.Step2FileInfoTransferred {
					ticker.Stop()

					r := FromCroc(receiver)
					recv = &r

					encoder := json.NewEncoder(w)
					err = encoder.Encode(recv)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			case success := <-doneChan:
				if success && recv == nil {
					r := FromCroc(receiver)
					recv = &r

					encoder := json.NewEncoder(w)
					err = encoder.Encode(recv)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				return
			}
		}
	}()

	err = receiver.Receive()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		doneChan <- false
		return
	}
	doneChan <- true
	wg.Wait()
}

func receive(w http.ResponseWriter, r *http.Request) {
	idStr := filepath.Base(r.URL.Path)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprint("Invalid Receive ID ", idStr), http.StatusBadRequest)
		return
	}

	recv, ok := Server.ReceiveData[id]
	if !ok {
		http.Error(w, fmt.Sprint("Receive with ID ", id, " not found"), http.StatusBadRequest)
	}

	if len(recv.Files) == 0 {
		http.Error(w, "Receive does not have any files", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, recv.Files[0].Name)
}
