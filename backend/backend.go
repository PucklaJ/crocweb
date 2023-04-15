package backend

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	ReceiveDir = "receive_tmp"
)

var Server struct {
	RootDir string

	CurrentReceiveID uint64
	ReceiveData      map[uint64]Receive
	receiveDataMtx   sync.Mutex
}

func StartServer() {
	if wd, err := os.Getwd(); err == nil {
		Server.RootDir = wd
	} else {
		Server.RootDir = "."
	}
	Server.RootDir = filepath.Join(Server.RootDir, "frontend")
	Server.ReceiveData = make(map[uint64]Receive)

	// Change into temporary file receive directory
	err := os.RemoveAll(ReceiveDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to clean receive dir:", err)
	}
	err = os.Mkdir(ReceiveDir, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create receive dir:", err)
	}
	err = os.Chdir(ReceiveDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to change directory to receive dir:", err)
	}

	http.HandleFunc("/", root)
	http.HandleFunc("/code/", code)
	http.HandleFunc("/receive/", receive)
	http.HandleFunc("/download/", download)

	cleanUpTicker := time.NewTicker(time.Second)
	cleanUpExit := make(chan bool)
	go func() {
		for {
			select {
			case <-cleanUpTicker.C:
				Server.receiveDataMtx.Lock()

				var recvToDelete []uint64
				for id, recv := range Server.ReceiveData {
					recv.time -= time.Second
					if recv.time <= 0 {
						recvToDelete = append(recvToDelete, id)
						// Delete all files
						for _, f := range recv.Files {
							err := os.Remove(f.Name)
							if err != nil {
								fmt.Fprintf(os.Stderr, "Failed to delete \"%s\": %s", f.Name, err)
							}
						}
					}
					Server.ReceiveData[id] = recv
				}
				for _, id := range recvToDelete {
					delete(Server.ReceiveData, id)
				}

				Server.receiveDataMtx.Unlock()
			case <-cleanUpExit:
				return
			}
		}
	}()

	address := "0.0.0.0:8080"
	fmt.Println("Listening on", address)

	err = http.ListenAndServe(address, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Closing Server ...")
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		cleanUpExit <- false
		os.Exit(1)
	}

	cleanUpExit <- true
}
