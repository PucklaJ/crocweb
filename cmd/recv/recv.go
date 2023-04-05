package main

import (
	"os"

	"github.com/PucklaJ/crocweb/global"
	"github.com/schollz/croc/v9/src/croc"
)

func main() {
	if len(os.Args) != 2 {
		panic("Invalid Arguments")
	}

	sharedSecret := os.Args[1]

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
		panic(err)
	}

	err = receiver.Receive()
	if err != nil {
		panic(err)
	}
}
