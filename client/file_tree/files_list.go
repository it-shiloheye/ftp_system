package filetree

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/it-shiloheye/ftp_system/client/logging"
	"github.com/it-shiloheye/ftp_system_lib/base"
)

type FilesListStruct struct {
	m          sync.RWMutex
	Included   []string                     `json:"included"`
	Excluded   []string                     `json:"excluded"`
	Extensions []string                     `json:"extensions"`
	Files      map[string]*FileObjectStruct `json:"files"`

	f0c  chan string
	logc chan string
}

func newFilesList() (fl FilesListStruct) {
	fl = FilesListStruct{
		Included: config.Include,
		Excluded: config.Exclude,
		Files:    map[string]*FileObjectStruct{},
		f0c:      make(chan string),
		logc:     make(chan string),
	}

	fl.readFileTreeJSON(file_tree_file)

	return
}

func (fl *FilesListStruct) writeFileTreeJSON() {

	fl.m.Lock()
	fl.Extensions = unique_exts
	file_tree_json, _ := json.MarshalIndent(fl, " ", "\t")
	fl.m.Unlock()
	err := base.WriteFile(file_tree_file, file_tree_json)
	if err != nil {
		log.Fatalln(err)
	}

}

func (fl *FilesListStruct) readFileTreeJSON(f_path string) (err error) {
	fl.m.Lock()
	defer fl.m.Unlock()

	f, err := base.OpenFile(f_path, os.O_RDONLY)
	if err != nil && !errors.Is(err, io.EOF) {
		err = fmt.Errorf(fmt.Sprint(logging.ErrorLocation("readFileTreeJSON"), " OpenFile ", f_path, " ", err))
		return
	}
	fs, _ := f.Stat()
	data := make([]byte, fs.Size()+1)
	n, err := f.Read(data)
	if err != nil && !errors.Is(err, io.EOF) {

		err = fmt.Errorf(fmt.Sprint(logging.ErrorLocation("readFileTreeJSON"), " Read ", f_path, " ", err))
		return
	}

	if n < 1 {
		err = fmt.Errorf(fmt.Sprint(logging.ErrorLocation("readFileTreeJSON"), " Read ", f_path, " ", "file read is less than 1 byte"))
		return
	}
	err = json.Unmarshal(data[:n], fl)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			err = fmt.Errorf("%v\nsyntax error at byte offset %d", err, e.Offset)
		}

		err = fmt.Errorf(fmt.Sprint(logging.ErrorLocation("readFileTreeJSON"), " Unmarshall ", f_path, "\n", err, "\n", string(data[:n])))
		return
	}
	return
}

func (fl *FilesListStruct) GenHashes(ctx context.Context, cancel context.CancelCauseFunc) (err error) {
	defer Logger.Done()

	err_loc := logging.ErrorLocation("fl.GenHashes")
	tm := time.NewTicker(time.Minute)
	ok := true
	files := []string{}
	idx := ""

	defer func() {
		if err != nil {
			cancel(err)
		}
	}()

	count := struct {
		Total  atomic.Int32
		Hashed atomic.Int32
	}{
		Total:  atomic.Int32{},
		Hashed: atomic.Int32{},
	}
	count.Total.Store(0)
	count.Hashed.Store(0)

	var bs *BytesStore
	i := 0
	// MainLooop:
	for ok {

		select {
		case <-ctx.Done():
			ok = false

		case idx, ok = <-fl.f0c:
			// Logger.Logf("Adding: %s to hash queue", idx)
			files = append(files, idx)
			count.Total.Add(1)
			if ok {

				continue
			}
		case <-tm.C:

		case bs = <-bytestore_queue:
			if len(files) == i {
				<-time.After(time.Millisecond * 100)
				bytestore_queue <- bs

			} else {
				// log.Fatalln("TRYING TO HASH")

				file_ := files[i]

				fl.m.RLock()
				f, ok := fl.Files[file_]
				fl.m.RUnlock()
				Logger.Add()
				if ok {
					go func(fo *FileObjectStruct, bs *BytesStore) {
						defer Logger.Done()
						err := hash_fileobject(f, bs)
						if err != nil {
							Logger.LogErr(err_loc, "hash_fileobject", err)
							return
						}
						count.Hashed.Add(1)
						fl.logc <- Logger.Log("hashed file: ", fo.Name, "\nhash:\t", fo.Hash)

					}(f, bs)

					i += 1

				} else {
					<-time.After(time.Millisecond * 100)
				}

			}
		}

		if (i + 1) == len(files) {
			clear(files)
			i = 0
		}

	}

	return
}
