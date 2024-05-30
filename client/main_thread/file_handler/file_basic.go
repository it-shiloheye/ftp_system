package filehandler

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"strings"
	"time"

	"os"

	"github.com/ftp_system_client/base"
	ftp_context "github.com/ftp_system_client/main_thread/context"
)

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
	d := string(buf)
	if len(d) > 1000 {
		d = d[:999]
	}
	if len(buf) < 1 {
		fo.Err = ftp_context.NewLogItem("FileBasic.Write", true).
			SetMessage("No data to write").
			Set("path", fo.Path).
			Set("buf", d).
			AppendParentError(err)
		return n, fo.Err
	}

	if fo.fo == nil {
		err_tmp := fo.Create().Err

		if err_tmp != nil {

			fo.Err = ftp_context.NewLogItem("FileBasic.Write", true).
				Set("after", "fo.Open()").
				SetMessage("unable to open file").
				Set("path", fo.Path).
				Set("buf", d).
				AppendParentError(err_tmp, err)

			return n, fo.Err
		}

		fo.Err = nil
		err = nil

	}

	_n, err := fo.fo.Write(buf)
	n = int(_n)
	if err != nil {
		log.Println(err)
		fo.Err = ftp_context.NewLogItem("FileBasic.Write", true).
			Set("after", "fo.fo.Write").
			Set("copied", n).
			Set("path", fo.Path).
			Set("buf", d).
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
			Set("after", "base.OpenFile").
			Set("path", fo.Path).
			SetMessage("unable to open file").
			AppendParentError(err)
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
	if fo.d.IsDir() {
		fo.Type = "dir"
		return ""
	}
	stp_1 := strings.Split(fo.Name, ".")
	stp_2 := len(stp_1)
	stp_3 := stp_1[stp_2-1]
	if len(stp_3) > 4 {
		return "unknown"
	}

	fo.Type = stp_3

	return stp_3
}

func (fo *FileBasic) Create() *FileBasic {
	loc := "FileBasic.Create"
	if fo.fo != nil {
		return fo
	}
	var err error

	fo.fo, err = base.OpenFile(fo.Path, os.O_RDWR|os.O_SYNC|os.O_CREATE)
	if err_pre := err; err != nil {
		if errors.Is(err, os.ErrNotExist) {

			d := strings.Split(fo.Path, "/")
			dir := strings.Join(d[:len(d)-1], "\\")
			err = os.MkdirAll(dir, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
			if err != nil && !errors.Is(err, os.ErrExist) {
				fo.Err = ftp_context.NewLogItem(loc, true).
					Set("after", "os.MkdirAll").
					Set("path", fo.Path).
					Set("Mkdir", dir).
					AppendParentError(err, err_pre)

				return fo
			}
			return fo.Create()
		}
		fo.Err = ftp_context.NewLogItem(loc, true).
			Set("after", "base.OpenFile").
			Set("path", fo.Path).
			SetMessagef("failed to open file").AppendParentError(err)

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
			SetMessagef("fo.Write").
			Set("path", fo.Path).
			SetMessage(err.Error()).AppendParentError(err, fo.Err)
		return fo.Err
	}
	return nil
}
func (fo *FileBasic) ResetError() {
	fo.Err = nil
}
