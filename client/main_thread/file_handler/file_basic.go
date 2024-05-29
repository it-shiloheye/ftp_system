package filehandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"strings"
	"time"

	"os"

	"github.com/ftp_system_client/base"
	ftp_context "github.com/ftp_system_client/main_thread/context"
)

type FileBasic struct {
	Name string             `json:"name"`
	Path string             `json:"path"`
	Err  ftp_context.LogErr `json:"error"`
	fo   *os.File
	fs   os.FileInfo
	d    fs.DirEntry
}

func init() {
	var _ io.ReadWriteCloser = &FileBasic{}
}

func (fo *FileBasic) Close() error {
	return fo.fo.Close()
}

func (fo *FileBasic) Read(buf []byte) (n int, err error) {
	n, err = io.ReadFull(fo.fo, buf)
	if err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Read", true).
			Set("after", "io.ReadFull").
			Set("path", fo.Path).AppendParentError(err)
		return n, fo.Err
	}
	return
}

func (fo *FileBasic) Write(buf []byte) (n int, err error) {
	_n, err := io.CopyN(fo.fo, bytes.NewReader(buf), int64(len(buf)))
	n = int(_n)
	if err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Write", true).
			Set("after", "io.CopyN").
			Set("path", fo.Path).
			AppendParentError(err)
		return n, fo.Err
	}
	return

}

func (fo *FileBasic) IsOpen() bool {
	return fo.fo != nil
}

func (fo *FileBasic) ReadAll() (data []byte, err error) {
	data, err = io.ReadAll(fo.fo)
	if err != nil {
		err = ftp_context.NewLogItem("FileBasic.ReadAll", true).
			Set("after", "io.ReadAll").
			Set("path", fo.Path).AppendParentError(err).
			SetMessage(err.Error())

	}
	return
}

func (fo *FileBasic) ModTime() string {
	return fo.fs.ModTime().Format(time.RFC822Z)
}

func (fo *FileBasic) Open() *FileBasic {
	if fo.fo != nil {
		return fo
	}
	var err error

	fo.fo, err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC)
	if err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).
			SetMessagef("base.OpenFile %s error:\n%s", fo.Path, err).AppendParentError(err)
		return fo
	}

	fo.fs, err = fo.fo.Stat()
	if err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).
			SetMessagef("fo.fo.Stat %s error:\n%s", fo.Path, err).AppendParentError(err)
		return fo
	}

	return fo
}

func (fo *FileBasic) Ext() string {
	stp_1 := strings.Split(fo.Name, ".")
	stp_2 := len(stp_1)
	stp_3 := stp_1[stp_2-1]
	if len(stp_3) > 4 {
		return "unknown"
	}

	return stp_3
}

func (fo *FileBasic) Create() *FileBasic {
	loc := "FileBasic.Create"
	if fo.fo != nil {
		return fo
	}

	var err error
	fo.fo, err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC|os.O_CREATE)
	if err != nil {
		fo.Err = ftp_context.NewLogItem(loc, true).
			SetMessagef("base.OpenFile %s error:\n%s", fo.Path, err).AppendParentError(err)

		return fo
	}
	fo.fs, err = fo.fo.Stat()
	if err != nil {
		fo.Err = ftp_context.NewLogItem(loc, true).
			SetMessagef("fo.fo.Stat %s error:\n%s", fo.Path, err).AppendParentError(err)
		return fo
	}

	return fo
}

func (fo *FileBasic) CreateWithDir() *FileBasic {
	if fo.fo != nil {
		return fo
	}

	var err error
	_dir := strings.Split(fo.Path, "\\")
	l := len(_dir) - 2
	if l > 1 {
		dir := _dir[:l]

		err = os.MkdirAll(strings.Join(dir, "\\"), fs.FileMode(base.S_IRWXO|base.S_IRWXU))
		if err != nil && !errors.Is(err, os.ErrExist) {
			fo.Err = ftp_context.NewLogItem("FileBasic.CreateWithDir", true).
				SetMessagef("os.MkdirAll %s error:\n%s", dir, err).AppendParentError(err)
			return fo
		}
	}

	fo.fo, err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC)
	if err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).
			SetMessagef("base.OpenFile %s error:\n%s", fo.Path, err).AppendParentError(err)
		return fo
	}

	fo.fs, err = fo.fo.Stat()
	if err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).
			SetMessagef("fo.fo.Stat %s error:\n%s", fo.Path, err).AppendParentError(err)
		return fo
	}

	return fo
}

func NewFileBasic(path string) (fo *FileBasic) {
	fo = &FileBasic{
		Path: path,
	}
	return
}

func (fo *FileBasic) WriteJson(v any) ftp_context.LogErr {
	loc := "FileBasic.WriteJson"
	var err error

	t, err := json.MarshalIndent(v, " ", "\t")
	if err != nil {
		fo.Err = ftp_context.NewLogItem(loc, true).
			Set("after", "json.MarshalIndent").
			Set("path", fo.Path).
			SetMessage(err.Error()).AppendParentError(err, fo.Err)
		return fo.Err
	}

	_, err = fo.Write(t)
	if err != nil {
		fo.Err = ftp_context.NewLogItem(loc, true).
			SetMessagef("json.MarshalIndent").
			Set("path", fo.Path).
			SetMessage(err.Error()).AppendParentError(err, fo.Err)
		return fo.Err
	}
	return nil
}
func (fo *FileBasic) ResetError() {
	fo.Err = nil
}
