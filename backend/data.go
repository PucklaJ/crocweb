package backend

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	mp4 "github.com/abema/go-mp4"
	"github.com/schollz/croc/v9/src/croc"
)

type ReceiveFile struct {
	Name   string                 `json:"name"`
	Size   int64                  `json:"size"`
	Custom map[string]interface{} `json:"custom"`
}

type Receive struct {
	ID    uint64        `json:"id"`
	Files []ReceiveFile `json:"files"`

	time time.Duration
}

const ReceiveTime = time.Minute * 5

func FromCroc(receiver *croc.Client) (recv Receive) {
	for _, fi := range receiver.FilesToTransfer {
		recv.Files = append(recv.Files, ReceiveFile{
			Name: filepath.Join(fi.FolderRemote, fi.Name),
			Size: fi.Size,
		})
	}

	sort.Slice(recv.Files, func(i, j int) bool {
		return recv.Files[i].Name < recv.Files[j].Name
	})

	recv.time = ReceiveTime

	Server.receiveDataMtx.Lock()
	recv.ID = Server.CurrentReceiveID
	recv.GenerateCustom()
	Server.ReceiveData[recv.ID] = recv
	Server.CurrentReceiveID++
	Server.receiveDataMtx.Unlock()

	return
}

func (recv *Receive) GenerateCustom() {
	waitChan := make(chan bool, runtime.GOMAXPROCS(-1))
	var wg sync.WaitGroup
	for _i := range recv.Files {
		waitChan <- true

		wg.Add(1)
		go func(i int) {
			defer func() {
				<-waitChan
				wg.Done()
			}()

			ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(recv.Files[i].Name), "."))

			switch ext {
			case "png", "jpg", "jpeg", "gif":
				recv.Files[i].SetCustomImage()
			case "mp4", "mov":
				recv.Files[i].SetCustomMP4()
			}
		}(_i)
	}

	wg.Wait()
}

func (rf *ReceiveFile) SetCustomImage() {
	file, err := os.Open(rf.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open \"%s\" to create custom: %s\n", rf.Name, err)
		return
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load image config of \"%s\" to create custom: %s\n", rf.Name, err)
		return
	}

	rf.Custom = make(map[string]interface{})
	rf.Custom["width"] = cfg.Width
	rf.Custom["height"] = cfg.Height
}

func (rf *ReceiveFile) SetCustomMP4() {
	file, err := os.Open(rf.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open \"%s\" to create custom: %s\n", rf.Name, err)
		return
	}
	defer file.Close()

	info, err := mp4.Probe(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to probe \"%s\" video to create custom: %s\n", rf.Name, err)
		return
	}

	width := info.Tracks[0].AVC.Width
	height := info.Tracks[0].AVC.Height

	rf.Custom = make(map[string]interface{})
	rf.Custom["width"] = width
	rf.Custom["height"] = height
}
