package backend

import (
	"path/filepath"
	"sort"

	"github.com/schollz/croc/v9/src/croc"
)

type ReceiveFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Receive struct {
	ID    uint64        `json:"id"`
	Files []ReceiveFile `json:"files"`
}

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

	Server.receiveDataMtx.Lock()
	recv.ID = Server.CurrentReceiveID
	Server.ReceiveData[recv.ID] = recv
	Server.CurrentReceiveID++
	Server.receiveDataMtx.Unlock()

	return
}
