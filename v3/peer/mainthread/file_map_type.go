package mainthread

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	ftp_base "github.com/it-shiloheye/ftp_system/v3/lib/base"
	"github.com/it-shiloheye/ftp_system/v3/lib/file_handler/v2"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

const file_map_json_path = "./file_map.json"

type FileMapType struct {
	l  sync.RWMutex
	fl *filehandler.LockFile

	ToUpload   bool                            `json:"to_upload"`
	ToDownload bool                            `json:"to_download"`
	FileMap    *ftp_base.MutexedMap[*FileItem] `json:"file_map"`

	HashMap *ftp_base.MutexedMap[FileState] `json:"hash_map"`
}

func (fmtype *FileMapType) MarshallText() ([]byte, error) {

	return json.MarshalIndent(fmtype, " ", "\t")
}

func NewFileMap() *FileMapType {
	return &FileMapType{
		FileMap: ftp_base.NewMutexedMap[*FileItem](),
		HashMap: ftp_base.NewMutexedMap[FileState](),
		fl:      &filehandler.LockFile{Name: "./locks/file_map.lock"},
	}
}

func (f *FileMapType) Reset() {
	f.ToUpload = false
	f.ToDownload = false
}

func (f *FileMapType) LoadFileMap() error {
	f.l.Lock()
	defer f.l.Unlock()

	d, err1 := os.ReadFile(file_map_json_path)
	if err1 != nil {
		if errors.Is(err1, os.ErrNotExist) {
			return nil
		}
		return Logger.LogErr(log_item.Locf(`d, err1 := os.ReadFile("%s")`, file_map_json_path), err1)

	}

	err2 := json.Unmarshal(d, f)
	if err2 != nil {
		Logger.LogErr(log_item.Locf(`err2 := json.Unmarshal(d:"./file_map.json",f)`), err2)
	}

	return nil
}

func (f *FileMapType) SaveFileMap() error {
	f.l.RLock()
	defer f.l.RUnlock()

	d, err1 := json.MarshalIndent(f, " ", "\t")
	if err1 != nil {
		return Logger.LogErr(log_item.Locf(`d, err1 := json.MarshalIndent(f: FileMapType," ","\t")`), err1)
	}

	err2 := os.WriteFile(file_map_json_path, d, ftp_base.FS_MODE)
	if err2 != nil {
		return Logger.LogErr(log_item.Locf(`err2 := os.WriteFile(file_map_json_path: "%s",d,ftp_base.FS_MODE)`, file_map_json_path), err1)
	}

	return nil
}
