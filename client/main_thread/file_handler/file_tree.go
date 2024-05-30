package filehandler

import (
	"io/fs"
	"os"
	"time"

	"github.com/ftp_system_client/base"
	ftp_context "github.com/ftp_system_client/main_thread/context"
)

type FileType = string

type FileBasic struct {
	Name string             `json:"name"`
	Path string             `json:"path"`
	Err  ftp_context.LogErr `json:"error"`
	Type FileType           `json:"type"`
	Size int64              `json:"size"`
	fo   *os.File
	fs   os.FileInfo
	d    fs.DirEntry
}

type FileTreeStruct struct {
	FileBasic
	Directory string      `json:"directory"`
	FilesList []*FileHash `json:"files_list"`

	LatesttUploadTime time.Time                    `json:"latest_upload_time"`
	LatestHash        time.Time                    `json:"latest_hash"`
	HashQueue         base.MutexedQueue[*FileHash] `json:"hash_queue"`
	UploadQueue       base.MutexedQueue[*FileHash] `json:"upload_queue"`
}

type FileHash struct {
	FileBasic
	Hash        string    `json:"hash"`
	LastModTime time.Time `json:"last_mod_time"`
	PrevModTime time.Time `json:"prev_mod_time"`
}

func init() {

}

func (fts *FileTreeStruct) EnqueueHashing(fh *FileHash) {
	fts.HashQueue.Enqueue(fh)
}

func (fts *FileTreeStruct) DequeueHashing() <-chan int64 {

	return fts.HashQueue.Dequeue()
}

func (fts *FileTreeStruct) MarkDoneHashing(n int64) {
	fts.HashQueue.MarkDone(n)
}

func (fts *FileTreeStruct) EnqueueUpload(fh *FileHash) {
	fts.UploadQueue.Enqueue(fh)
}

func (fts *FileTreeStruct) DequeueUpload() <-chan int64 {

	return fts.UploadQueue.Dequeue()
}

func (fts *FileTreeStruct) MarkDoneUpload(n int64) {
	fts.UploadQueue.MarkDone(n)
}
