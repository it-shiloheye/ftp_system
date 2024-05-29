package filehandler

import (
	"bytes"
	"io"
	"io/fs"
	"strings"
	"time"

	"os"

	"github.com/ftp_system_client/base"
	ftp_context "github.com/ftp_system_client/main_thread/context"
)

type FileBasic struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Err  error  `json:"error"`
	fo   *os.File
	fs   os.FileInfo
	d    fs.DirEntry
}

func (fo *FileBasic) Read(buf []byte) (n int, err error) {
	n, err = io.ReadFull(fo.fo, buf)
	if err != nil {
		err = ftp_context.NewLogItem("FileBasic.Read", true).SetMessagef("io.ReadFull %s error:\n%s", fo.Path, err.Error())
	}
	return
}

func (fo *FileBasic) Write(buf []byte) (n int, err error) {
	_n, err := io.CopyN(fo.fo, bytes.NewReader(buf), int64(len(buf)))
	n = int(_n)
	return

}

func (fo *FileBasic) IsOpen() bool {
	return fo.fo != nil
}

func (fo *FileBasic) ReadAll() (data []byte, err error) {
	data, err = io.ReadAll(fo.fo)
	if err != nil {
		err = ftp_context.NewLogItem("FileBasic.ReadAll", true).SetMessagef("io.ReadAll %s error:\n%s", fo.Path, err.Error())
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

	fo.fo, fo.Err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC)
	if fo.Err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("base.OpenFile %s error:\n%s", fo.Path, fo.Err.Error())
		return fo
	}

	fo.fs, fo.Err = fo.fo.Stat()
	if fo.Err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("fo.fo.Stat %s error:\n%s", fo.Path, fo.Err.Error())
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
	if fo.fo != nil {
		return fo
	}

	fo.fo, fo.Err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC|os.O_CREATE)
	if fo.Err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("base.OpenFile %s error:\n%s", fo.Path, fo.Err.Error())

		return fo
	}
	fo.fs, fo.Err = fo.fo.Stat()
	if fo.Err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("fo.fo.Stat %s error:\n%s", fo.Path, fo.Err.Error())
		return fo
	}

	return fo
}

func (fo *FileBasic) CreateWithDir() *FileBasic {
	if fo.fo != nil {
		return fo
	}

	_dir := strings.Split(fo.Path, "\\")
	l := len(_dir) - 2
	if l > 1 {
		dir := _dir[:l]

		fo.Err = os.MkdirAll(strings.Join(dir, "\\"), fs.FileMode(base.S_IRWXO|base.S_IRWXU))
		if fo.Err != nil {
			return fo
		}
	}

	fo.fo, fo.Err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC)
	if fo.Err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("base.OpenFile %s error:\n%s", fo.Path, fo.Err.Error())
		return fo
	}

	fo.fs, fo.Err = fo.fo.Stat()
	if fo.Err != nil {
		fo.Err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("fo.fo.Stat %s error:\n%s", fo.Path, fo.Err.Error())
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
