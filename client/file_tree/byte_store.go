package filetree

import (
	"errors"
	"hash"
	"io"
	"log"
	"time"

	"fmt"
	"os"

	"github.com/ftp_system_client/logging"
	"github.com/it-shiloheye/ftp_system_lib/base"
)

type BytesStore struct {
	h     hash.Hash
	bytes []byte
}

func (bs *BytesStore) Hash() (hash string) {
	// log.Println("hashFile called")

	bs.h.Write(bs.bytes)
	hash = fmt.Sprintf("%x", bs.h.Sum(nil))
	bs.h.Reset()

	return
}

func (bs *BytesStore) ReadFile(file_path string) (f_ *os.File, err error) {
	err_loc := logging.ErrorLocation("bs.ReadFile")
	f_, err = base.OpenFile(file_path, os.O_RDWR|os.O_SYNC)
	if err != nil {
		Logger.LogErr(err_loc, "base.OpenFile", err)
		return
	}

	f_s, _ := f_.Stat()

	diff := f_s.Size() - int64(len(bs.bytes))
	if diff > 0 {
		bs.bytes = make([]byte, len(bs.bytes)+int(2*diff))
	}

	_, err = f_.Read(bs.bytes)
	if err != nil && !errors.Is(err, io.EOF) {
		Logger.LogErr(err_loc, "f_.Read(bs.bytes)", err.Error())
		return
	}

	return
}

func (bs *BytesStore) Reset() {
	clear(bs.bytes)
}

const TimeFormat = time.RFC822Z

func hash_fileobject(fo *FileObjectStruct, bs *BytesStore) (err error) {

	if bs == nil {
		log.Fatalln("hash_fileobject(fo *FileObjectStruct, bs *BytesStore) error: nil bs pointer")
	}

	if fo == nil {
		log.Fatalln("hash_fileobject(fo *FileObjectStruct, bs *BytesStore) error: nil fo pointer")
	}
	defer func() {
		bs.Reset()
		bytestore_queue <- bs
	}()

	fo.file, err = bs.ReadFile(fo.Path)
	if err != nil && !errors.Is(err, io.EOF) {
		err = fmt.Errorf("hash_fileobject(fo *FileObjectStruct):\nbs.readFile(fo.Path):\n%v", err)
		return
	}
	defer fo.file.Close()

	f_stats, err := fo.file.Stat()
	if err != nil {
		err = fmt.Errorf("hash_fileobject(fo *FileObjectStruct):\nfile.State() :\n%v", err)
		return
	}

	f_mod := f_stats.ModTime()

	if len(fo.ModTime) < 1 {
		fo.ModTime = f_mod.Format(TimeFormat)

		fo.Hash = bs.Hash()

		return
	}

	fo_mod, err := time.Parse(TimeFormat, fo.ModTime)
	if err != nil {
		err = fmt.Errorf("hash_fileobject(fo *FileObjectStruct):\ntime.Parse(time.RFC822Z,fo.ModTime) :\n%v", err)
		return
	}

	if fo_mod == f_mod {
		return
	}
	return
}
