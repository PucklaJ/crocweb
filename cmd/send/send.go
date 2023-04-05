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

	sendFileName := os.Args[1]
	sharedSecret := "crocweb-test"

	sender, err := croc.New(croc.Options{
		IsSender:       true,
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

	fi, ef, num, err := croc.GetFilesInfo([]string{sendFileName}, false)
	if err != nil {
		panic(err)
	}

	err = sender.Send(fi, ef, num)
	if err != nil {
		panic(err)
	}
}
