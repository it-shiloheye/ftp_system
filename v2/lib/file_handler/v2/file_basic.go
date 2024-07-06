package filehandler

import (
	"encoding/json"

	"io"
	"io/fs"
	"strings"
	"time"

	"os"

	"github.com/it-shiloheye/ftp_system/v2/lib/base"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging/log_item"
)

type FileType string

type FileBasic struct {
	Path     string   `json:"path"`
	FileType FileType `json:"type"`
	Size     int64    `json:"size"`
	// directly read write to the file
	*os.File
	fs os.FileInfo
	d  fs.DirEntry
}

func init() {
	var _ io.ReadWriteCloser = &FileBasic{}

	var _ fs.DirEntry = &FileBasic{}
}

func (Fo *FileBasic) IsOpen() bool {
	return Fo.File != nil
}

func (Fo *FileBasic) Fs() os.FileInfo {
	return Fo.fs
}

func (Fo *FileBasic) Info() (os.FileInfo, error) {
	return Fo.d.Info()
}

func (Fo *FileBasic) IsDir() bool {
	return Fo.d.IsDir()
}

func (Fo *FileBasic) Type() fs.FileMode {
	return Fo.d.Type()
}

func (Fo *FileBasic) ReadAll() (data []byte, err error) {
	data, err = io.ReadAll(Fo.File)
	if err != nil {
		err = log_item.NewLogItem("FileBasic.ReadAll", log_item.LogLevelError01).
			SetAfter("data, err = io.ReadAll(Fo.File)").
			Set("path", Fo.Path).AppendParentError(err).
			SetMessage(err.Error())

	}
	return
}

func (Fo *FileBasic) ModTime() string {
	return Fo.fs.ModTime().Format(time.RFC822Z)
}

func Open(file_path string) (Fo *FileBasic, err error) {
	loc := log_item.Loc("Fo.File, err = base.OpenFile(Fo.Path, os.O_RDWR|os.O_SYNC)")
	var err1, err2 error
	Fo = &FileBasic{
		Path: file_path,
	}
	Fo.File, err1 = base.OpenFile(file_path, os.O_RDWR|os.O_SYNC)
	if err1 != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("base.OpenFile").
			Set("path", Fo.Path).
			SetMessagef("path: %s \nerror:\n%s", file_path, err1.Error()).
			AppendParentError(err1)

		return

	}

	Fo.fs, err2 = Fo.File.Stat()
	if err2 != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("Fo.fs, err2 = Fo.File.Stat()").
			SetMessagef("path: %s \nerror:\n%s", file_path, err2.Error()).AppendParentError(err2)
		return
	}

	Fo.d = fs.FileInfoToDirEntry(Fo.fs)

	Fo.FileType = Ext(Fo)
	return
}

func Ext(Fo *FileBasic) FileType {
	if Fo == nil {
		panic("Provided nil pointer to filehandler.Ext(Fo *FileBasic) FileType")
	}

	if len(Fo.FileType) > 0 {
		return Fo.FileType
	}

	if Fo.IsDir() {
		Fo.FileType = "dir"
		return Fo.FileType
	}
	stp_1 := strings.Split(Fo.Path, ".")
	stp_2 := len(stp_1)
	stp_3 := stp_1[stp_2-1]
	if len(stp_3) > 4 {
		Fo.FileType = "unknown"
		return Fo.FileType
	}

	Fo.FileType = FileType(stp_3)
	return Fo.FileType
}

func Create(file_path string) (Fo *FileBasic, err error) {
	loc := log_item.Loc("Create(file_path string) (Fo *FileBasic,err error)")
	var err1, err2 error
	Fo = &FileBasic{
		Path: file_path,
	}
	Fo.File, err1 = base.OpenFile(file_path, os.O_RDWR|os.O_SYNC|os.O_CREATE)
	if err1 != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("Fo.File, err1 = base.OpenFile(file_path, os.O_RDWR|os.O_SYNC|os.O_CREATE)").
			SetMessagef("path: %s \nerror:\n%s", file_path, err1.Error()).
			AppendParentError(err1)
		return

	}

	Fo.fs, err2 = Fo.File.Stat()
	if err2 != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("Fo.fs, err2 = Fo.File.Stat()").
			SetMessagef("path: %s \nerror:\n%s", file_path, err2.Error()).AppendParentError(err2)
		return
	}

	Fo.d = fs.FileInfoToDirEntry(Fo.fs)
	Fo.FileType = Ext(Fo)
	return
}

func NewFileBasic(path string) (Fo *FileBasic) {
	Fo = &FileBasic{
		Path: path,
	}
	return
}

func WriteJson(Fo *FileBasic, v any) (err error) {
	loc := log_item.Loc("WriteJson(Fo *FileBasic, v any)(err error)")
	var err1, err2 error
	t, err1 := json.MarshalIndent(v, " ", "\t")
	if err1 != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter(`t, err1 := json.MarshalIndent(v, " ", "\t")`).
			Set("path", Fo.Path).
			SetMessagef("path: %s \nerror:\n%s", Fo.Path, err1.Error()).
			AppendParentError(err1)
		return
	}

	_, err2 = Fo.Write(t)
	if err2 != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetMessagef("_, err2 = Fo.Write(t)").
			SetMessagef("path: %s \nerror:\n%s", Fo.Path, err2.Error()).
			AppendParentError(err2)
		return
	}
	return nil
}
