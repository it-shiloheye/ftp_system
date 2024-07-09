package mainthread

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"

	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	db "github.com/it-shiloheye/ftp_system/v3/lib/db_access"
	db_access "github.com/it-shiloheye/ftp_system/v3/lib/db_access/generated"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system/v3/peer/config"
)

var storage_struct = server_config.Storage

var Logger = logging.Logger

var DB = db.DB

func tickerf(loc *log_item.Loc, n int, v ...any) {
	Logger.Logf(*loc, "\n%03d\t%s", n, fmt.Sprint(v...))
}

func Loop(ctx ftp_context.Context) error {
	loc := log_item.Locf(`func Loop(ctx ftp_context.Context)=>error`)
	after := "setup "
	defer ctx.Finished()
	defer recover_func(&after, loc)

	files_map := NewFileMap()
	files_map.LoadFileMap()
	file_watcher := SetUpCurrentWatcher(files_map)
	var ev fsnotify.Event
	for ok := true; ok; {
		files_map.Reset()
		log.Println("new loop in mainthread")
		child_ctx := ctx.NewChild()
		t := time.Minute * time.Duration(storage_struct.PollIntervalMinutes)
		child_ctx.SetDeadline(t)

		switch storage_struct.PeerRole {
		case db_access.PeerRoleTypeClient:
			ctx.Add()
			func() {
				defer recover_func(&after, loc)

				files_map.FileMap.Clear()
				after = "WalkDir"
				log.Println(after)
				if err01 := WalkDir(storage_struct.StorageDirectory, files_map, file_watcher); err01 != nil {
					Logger.LogErr(loc, err01)
					<-time.After(time.Minute)

				}

				after = "UploadFunc"
				log.Println(after)
				if err02 := UploadFunc(child_ctx.Add(), files_map); err02 != nil {
					Logger.LogErr(loc, err02)
					<-time.After(time.Minute)

				}

				after = "ClientDownloadFunc"
				log.Println(after)
				if err03 := ClientDownloadFunc(child_ctx.Add(), files_map); err03 != nil {
					Logger.LogErr(loc, err03)
					<-time.After(time.Minute)

				}
				log.Println("client loop successful")

			}()
		case db_access.PeerRoleTypeStorage:
			after = "StorageDownloadFunc"
			StorageDownloadFunc(ctx, files_map)
		}

		select {
		case <-ctx.Done():
			after = "ctx.Done()"
			return ctx.Err()
		case <-time.After(t):
			after = "time.After(t)"
		case ev, ok = <-file_watcher.Events:
			after = "<-file_watcher.Events"
			log.Println(ev)
		}

		after = "server_config.LoopReadStorageStruct"
		server_config.LoopReadStorageStruct(3, storage_struct)
		after = "files_map.SaveFileMap"
		// files_map.SaveFileMap()
	}

	return nil
}

func StorageDownloadFunc(ctx *ftp_context.ContextStruct, files_map *FileMapType) (err error) {
	loc := log_item.Loc("StorageDownloadFunc(ctx *ftp_context.ContextStruct, file_map map[string]*FileItem, db_files_list []*db_access.GetFilesListRow) (err error)")

	files_map.FileMap.Lock()
	defer files_map.FileMap.Unlock()

	// file_map := files_map.FileMap

	db_conn := db.DBPool.GetConn()
	defer db.DBPool.Return(db_conn)

	_, err1 := DB.DownloadStoreBulk(ctx, db_conn, storage_struct.PeerId.Bytes)
	if err1 != nil {

		return Logger.LogErr(loc, err1)
	}

	return nil
}

func SetUpCurrentWatcher(files_map *FileMapType) (wc *fsnotify.Watcher) {
	loc := log_item.Loc(`func SetUpWatcher(files_map *FileMapType)`)

	file_map := files_map.FileMap
	log.Println("watcher created")
	watcher, err1 := fsnotify.NewWatcher()

	log.Println("watcher created 2")
	if err1 != nil {
		Logger.LogErr(loc, err1)
		return
	}
	wc = watcher

	filepath.WalkDir(storage_struct.StorageDirectory, func(fpath string, d fs.DirEntry, err error) error {
		if ok, err1 := is_excluded_regex(loc, fpath); is_excluded_dir(fpath) || ok || err1 != nil {
			if err1 != nil {
				return err1
			}
			return nil
		}
		if d.IsDir() {
			watcher.Add(fpath)
			return nil
		}
		stats, err3 := d.Info()
		if err3 != nil {
			return Logger.LogErr(loc, err3)
		}

		f_stored, ok := file_map.Get(fpath)
		if ok {

			f_stored.stats = stats
			return nil
		}

		if !ok {
			file_map.Set(fpath, &FileItem{
				FileState: fstate_new,
				Path:      fpath,
				stats:     stats,
			})
			return nil
		}

		return nil
	})

	log.Println("watcher set up finished")

	return
}

func recover_func(after *string, loc log_item.Loc) {
	logging.RecoverFunc(after, loc)
}
