package dir_handler

import (
	"sync"

	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var FileStorage = NewFileStorageStruct()

var lock_fs = sync.Mutex{}

type UploadedFileStruct struct {
	ClientId  string   `json:"client_id"`
	FilesList []string `json:"files_list"`
	LocalPath string   `json:"local_path"`
}

type FileStorageStruct struct {
	// store file_hash data
	FileMetaData ftp_base.MutexedMap[*filehandler.FileHash] `json:"files_metadata"`
	// store client id data
	ClientData ftp_base.MutexedMap[*UploadedFileStruct] `json:"client_data"`
}

func NewFileStorageStruct() *FileStorageStruct {
	lock_fs.Lock()
	return &FileStorageStruct{
		FileMetaData: ftp_base.NewMutexedMap[*filehandler.FileHash](),
		ClientData:   ftp_base.NewMutexedMap[*UploadedFileStruct](),
	}
}
