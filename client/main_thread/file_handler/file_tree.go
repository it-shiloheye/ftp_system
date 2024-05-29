package filehandler

import (
	"time"

	"github.com/ftp_system_client/base"
)

type FileTreeStruct struct {
	Directory string     `json:"directory"`
	FilesList []FileHash `json:"files_list"`

	LastUpload time.Time               `json:"last_upload"`
	LastHash   time.Time               `json:"last_hash"`
	HashQueue  base.MutexedMap[string] `json:"hash_queue"`
}

type FileHash struct {
	FileBasic
	LastModTime time.Time `json:"last_mod_time"`
	PrevModTime time.Time `json:"prev_mod_time"`
}

func init() {

}
