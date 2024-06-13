package server_dirhandler

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
	// store files uploaded by client
	ClientData ftp_base.MutexedMap[*UploadedFileStruct] `json:"client_data"`
	// uploadedhash
	UploadedHashes ftp_base.MutexedMap[string] `json:"uploaded_hashes"`
}

func NewFileStorageStruct() *FileStorageStruct {
	lock_fs.Lock()
	return &FileStorageStruct{
		FileMetaData:   ftp_base.NewMutexedMap[*filehandler.FileHash](),
		ClientData:     ftp_base.NewMutexedMap[*UploadedFileStruct](),
		UploadedHashes: ftp_base.NewMutexedMap[string](),
	}
}

var TmpFileData = ftp_base.NewMutexedMap[[]byte]()
