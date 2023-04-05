package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	var codeText, receiveHolder, requestButton *js.Object

	js.Global.Set("init", func() {
		document := js.Global.Get("document")
		codeText = document.Call("getElementById", "code_text")
		receiveHolder = document.Call("getElementById", "receive_holder")
		requestButton = document.Call("getElementById", "request_button")

		fmt.Println(codeText)
		fmt.Println("Frontend Loaded")
	})

	js.Global.Set("code_text_enter", func() {
		requestButton.Set("disabled", true)

		sharedSecret := codeText.Get("value").String()
		fmt.Printf("Requesting Code \"%s\" ...\n", sharedSecret)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			recv, err := CodeRequest(sharedSecret)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			Client.ReceiveData = recv

			recv.IntoHolder(receiveHolder)
			wg.Done()
		}()
		go func() {
			defer requestButton.Set("disabled", false)
			wg.Wait()

			fileContents, err := ReceiveRequest(Client.ReceiveData.ID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			ext := filepath.Ext(Client.ReceiveData.Files[0].Name)

			if strings.EqualFold(ext, ".txt") {
				element := js.Global.Get("document").Call("createElement", "p")
				element.Set("innerHTML", fileContents)
				receiveHolder.Call("appendChild", element)
			} else if strings.EqualFold(ext, ".png") || strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
				fmt.Println("Decoding Image ...")
				reader := strings.NewReader(fileContents)
				img, _, err := image.Decode(reader)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}

				width := img.Bounds().Size().X
				height := img.Bounds().Size().Y

				fmt.Println("Creating Array ...")
				arr := js.Global.Get("Uint8ClampedArray").New(width * height * 4)
				for x := 0; x < width; x++ {
					for y := 0; y < height; y++ {
						col := img.At(x, y)
						r, g, b, a := col.RGBA()

						arr.SetIndex(y*width*4+x*4+0, int(float64(r)/float64(0xffff)*255.0))
						arr.SetIndex(y*width*4+x*4+1, int(float64(g)/float64(0xffff)*255.0))
						arr.SetIndex(y*width*4+x*4+2, int(float64(b)/float64(0xffff)*255.0))
						arr.SetIndex(y*width*4+x*4+3, int(float64(a)/float64(0xffff)*255.0))
					}
				}

				fmt.Println("Creating ImageData ...")
				imgData := js.Global.Get("ImageData").New(arr, width, height)

				fmt.Println("Creating Canvas ...")
				canvas := js.Global.Get("document").Call("createElement", "canvas")
				ctx := canvas.Call("getContext", "2d")
				canvas.Set("width", width)
				canvas.Set("height", height)
				ctx.Call("putImageData", imgData, 0, 0)

				fmt.Println("Creating Image ...")
				htmlImg := js.Global.Get("Image").New()
				htmlImg.Set("src", canvas.Call("toDataURL"))
				htmlImg.Set("style", "width:450;")

				fmt.Println("Appending Child ...")
				receiveHolder.Call("appendChild", htmlImg)
			}
		}()
	})
}
