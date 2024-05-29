package filetree

import (
	"context"

	// "encoding/base64"

	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"time"

	"github.com/ftp_system_client/logging"
	cfg "github.com/ftp_system_client/main_thread/config"
)

type FileType int

var config *cfg.ConfigStruct = cfg.Config
var FilesList FilesListStruct = newFilesList()
var Logger *logging.LoggerStruct = logging.Logger

const (
	FT_Unknown FileType = iota
	FT_Directory
)

var (
	file_tree_file = config.FileTreeFile
)

func (fl *FilesListStruct) Append(fpath string, name string, ftype FileType) {

	new_fo := &FileObjectStruct{
		Name: name,
		Path: fpath,
		Type: ftype,
	}

	fl.m.Lock()
	fl.Files[fpath] = new_fo
	fl.m.Unlock()
	if ftype != FT_Unknown && ftype != FT_Directory {
		fl.f0c <- fpath
	}

}

func WalkFileTree(ctx context.Context, cancel context.CancelCauseFunc) {
	walkFileTree(ctx, cancel)
}

func walkFileTree(ctx context.Context, cancel context.CancelCauseFunc) {
	defer Logger.Wait()
	c := ctx
	Logger.Add()
	go Logger.Go(c)

	Logger.Add()
	go FilesList.GenHashes(c, cancel)

	Logger.Add()
	go FilesList.GenHashes(c, cancel)

	if len(config.Include) < 1 {
		log.Fatal("need to include at least one directory")
	}
	var err error
	// defer Logger.Close()
	for _, directory := range config.Include {

		err = filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {

			for _, excluded := range config.Exclude {
				if strings.Contains(path, excluded) {
					return err
				}
			}
			name := d.Name()

			FilesList.Append(path, name, file_type(d))
			Logger.LogAppendFile(path, name, d.Type().String())
			return err
		})
		if err != nil {
			Logger.LogErrf("walkFileTree error:\n%s", err.Error())
			cancel(err)
			return
		}
	}

}

var unique_exts []string = []string{"unknown", "directory"}

func file_type(d fs.DirEntry) FileType {
	if d.IsDir() {
		return FT_Directory
	}

	spl_name := strings.Split(d.Name(), ".")
	ext := spl_name[len(spl_name)-1]
	if len(ext) > 5 {
		return FileType(0)
	}

	for i, e_ := range unique_exts {
		if e_ == ext {
			return FileType(i)
		}
	}

	unique_exts = append(unique_exts, ext)
	return file_type(d)
}

func (fl *FilesListStruct) BackgroundWriteFile(ctx context.Context) (err error) {
	defer Logger.Done()

	ok := true

	tm := time.NewTicker(time.Minute)
	logs := []string{}
	it := ""

	for ok {
		select {
		case <-ctx.Done():
			ok = false

		case it, ok = <-fl.logc:

			logs = append(logs, it)
			if ok {
				continue
			}
		case <-tm.C:
		}
		for _, log_ := range logs {
			if len(log_) < 1 {
				continue
			}
		}
		fl.writeFileTreeJSON()

		if len(logs) == len(fl.Files) {
			close(fl.f0c)
			close(fl.logc)
		}
	}

	return
}
