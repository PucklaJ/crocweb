package backend

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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

	address := "0.0.0.0:8080"

	fmt.Println("Listening on", address)

	err = http.ListenAndServe(address, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Closing Server ...")
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
