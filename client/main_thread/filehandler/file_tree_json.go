package dir_handler

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"sync"
	"time"

	"os"

	"github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"

	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
	// "golang.org/x/sync/syncmap"
)

var FileTree = NewFileTreeJson()

type FileTreeJson struct {
	lock    sync.RWMutex
	FileMap ftp_base.MutexedMap[*filehandler.FileHash] `json:"files"`

	HashQueue ftp_base.MutexedMap[*filehandler.FileHash] `json:"hash_queue"`
	Uploaded  ftp_base.MutexedMap[*filehandler.FileHash] `json:"uploaded"`
}

func init() {
	log.Println("loading filetree")

	FileTree.Lock()
	defer FileTree.Unlock()
	file_tree_path := ClientConfig.DataDir + "/file-tree.json"
	log.Println(file_tree_path)
	b, err1 := os.ReadFile(file_tree_path)
	if err1 != nil {
		if errors.Is(err1, os.ErrNotExist) {

			tmp, err2 := json.MarshalIndent(FileTree, " ", "\t")
			if err2 != nil {
				log.Fatalln(err2)
			}
			err3 := os.WriteFile(file_tree_path, tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
			if err2 != nil {
				log.Fatalln(err3)
			}

			log.Println("successfully loaded filetree")
			return
		}
		log.Fatalln(err1)
	}

	err3 := json.Unmarshal(b, FileTree)
	if err3 != nil {
		log.Fatalln(err3)
	}

	log.Println("successfully loaded filetree")
}

func NewFileTreeJson() *FileTreeJson {
	return &FileTreeJson{
		FileMap: ftp_base.NewMutexedMap[*filehandler.FileHash](),

		HashQueue: ftp_base.NewMutexedMap[*filehandler.FileHash](),
		Uploaded:  ftp_base.NewMutexedMap[*filehandler.FileHash](),
	}
}

func WriteFileTree(ctx ftp_context.Context) (err ftp_context.LogErr) {
	loc := "WriteFileTree() (err ftp_context.LogErr)"
	FileTree.RLock()
	defer FileTree.RUnlock()

	file_tree_path := ClientConfig.DataDir + "/file-tree.json"
	tmp, err1 := json.MarshalIndent(FileTree, " ", "\t")
	if err1 != nil {
		return &ftp_context.LogItem{Location: loc, Time: time.Now(),
			Err:       true,
			After:     `tmp, err1 := json.MarshalIndent(FileTree, " ", "\t")`,
			Message:   err1.Error(),
			CallStack: []error{err1},
		}
	}
	err2 := os.WriteFile(file_tree_path, tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
	if err2 != nil {
		return &ftp_context.LogItem{Location: loc, Time: time.Now(),
			Err:       true,
			After:     `err2 := os.WriteFile(file_tree_path, tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))`,
			Message:   err2.Error(),
			CallStack: []error{err2},
		}
	}

	return
}

var tc = time.NewTimer(time.Millisecond)

func (ft *FileTreeJson) Lock() {
	ft.lock.Lock()
	ft.HashQueue.Lock()
	ft.Uploaded.Lock()
	ft.FileMap.Lock()

}
func (ft *FileTreeJson) Unlock() {
	ft.HashQueue.Unlock()
	ft.Uploaded.Unlock()
	ft.FileMap.Unlock()
	ft.lock.Unlock()
}

func (ft *FileTreeJson) RLock() {
	ft.lock.RLock()
	ft.HashQueue.RLock()
	ft.Uploaded.RLock()
	ft.FileMap.RLock()

}
func (ft *FileTreeJson) RUnlock() {
	ft.HashQueue.RUnlock()
	ft.Uploaded.RUnlock()
	ft.FileMap.RUnlock()
	ft.lock.RUnlock()
}
