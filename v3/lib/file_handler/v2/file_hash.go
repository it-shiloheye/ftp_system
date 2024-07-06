package filehandler

import (
	"fmt"
	"time"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

type FileHash struct {
	*FileBasic
	Hash     string         `json:"hash"`
	ModTime  string         `json:"last_mod_time"`
	MetaData map[string]any `json:"meta_data"`
}

func HashFile(Fo *FileBasic, bs *BytesStore) (hash string, err error) {
	loc := log_item.Loc(" NewFileHash(Fo *FileBasic, bs *BytesStore)(hash string, err error)")
	if Fo == nil || Fo.File == nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			Message:  "FileBasic or os.File pointer is nil",
		}
		return
	}
	if bs == nil || bs.h == nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			Message:  "ByteStore pointer provided is invalid",
		}
		return
	}

	var err1, err2 error
	bs.Reset()

	_, err1 = bs.ReadFrom(Fo.File)
	if err1 != nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    `_, err1 = bs.ReadFrom(Fo.File)`,
			Message:  err1.Error(),
			Level:    log_item.LogLevelError02, CallStack: []error{err1},
		}
		return
	}

	hash, err2 = bs.Hash()
	if err2 != nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    `hash, err2 = bs.Hash()`,
			Message:  err2.Error(),
			Level:    log_item.LogLevelError02, CallStack: []error{err2},
		}
		return
	}
	return
}

func NewFileHashOpen(file_path string) (Fh *FileHash, err error) {
	loc := log_item.Loc("NewFileHashOpen(file_path string) (Fh *FileHash, err error)")
	Fh = &FileHash{}
	var err1 error
	Fh.FileBasic, err1 = Open(file_path)
	if err1 != nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    fmt.Sprintf(`Fh.FileBasic, err1 = Open("%s")`, file_path),
			Message:  err1.Error(),
			Level:    log_item.LogLevelError02, CallStack: []error{err1},
		}
		return
	}

	Fh.ModTime = fmt.Sprint(Fh.fs.ModTime())

	return
}

func NewFileHashCreate(file_path string) (Fh *FileHash, err error) {
	loc := log_item.Loc("NewFileHashOpen(file_path string) (Fh *FileHash, err error)")
	Fh = &FileHash{}
	var err1 error
	Fh.FileBasic, err1 = Create(file_path)
	if err1 != nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    fmt.Sprintf(`Fh.FileBasic, err1 = Create("%s")`, file_path),
			Message:  err1.Error(),
			Level:    log_item.LogLevelError02, CallStack: []error{err1},
		}
		return
	}

	return
}

func (fh *FileHash) Get(key string) (it any, ok bool) {
	it, ok = fh.MetaData[key]
	return
}

func (fh *FileHash) Set(key string, val any) {

	fh.MetaData[key] = val
}

func GetFileHash[T any](fh *FileHash, key string) (it T, ok bool) {

	a, ok_ := fh.Get(key)
	if !ok_ {
		return
	}

	it, ok = a.(T)
	return
}
