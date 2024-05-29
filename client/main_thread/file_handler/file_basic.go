package filehandler

import (
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

func (fo *FileBasic) Open() (ok bool, err error) {
	if fo.fo != nil {
		return true, nil
	}

	fo.fo, err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC)
	if err != nil {
		err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("base.OpenFile %s error:\n%s", fo.Path, err.Error())
	}
	ok = true

	fo.fs, err = fo.fo.Stat()
	if err != nil {
		err = ftp_context.NewLogItem("FileBasic.Open", true).SetMessagef("fo.fo.Stat %s error:\n%s", fo.Path, err.Error())
	}

	return
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
