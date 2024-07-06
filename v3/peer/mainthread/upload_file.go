package mainthread

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	db "github.com/it-shiloheye/ftp_system/v3/lib/db_access"
	db_access "github.com/it-shiloheye/ftp_system/v3/lib/db_access/generated"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"

	"github.com/fsnotify/fsnotify"
	"github.com/jackc/pgx/v5/pgtype"
)

const after_db_conn = `db_conn := db.DBPool.GetConn()
	files_list, err1 := DB.GetFilesList(ctx, db_conn)
	defer db.DBPool.Return(db_conn)
	db_files_list = files_list`

func UploadFunc(ctx ftp_context.Context, files_map *FileMapType) (err error) {
	loc := log_item.Loc(`func UploadFunc(ctx ftp_context.Context, file_map FileMap) error`)
	after := "setup"
	defer ctx.Finished()
	defer recover_func(&after, loc)
	// tickerf(&loc, 1, "at open of ticker")
	if has_deadline, near_deadline := ctx.NearDeadline(time.Second); has_deadline && near_deadline {
		return nil
	}
	db_conn := db.DBPool.GetConn()
	after = "db_files_list, err1 := DB.GetFilesList(ctx, db_conn)"
	db_files_list, err1 := DB.GetFilesList(ctx, db_conn)
	defer db.DBPool.Return(db_conn)

	if err1 != nil {
		if strings.Contains(err1.Error(), "no rows in result set") {
			return nil
		}
		err = Logger.LogErr(loc, &log_item.LogItem{
			After:     after_db_conn,
			CallStack: []error{err1},
		})
		return
	}

	db_file_map := FileMap{}
	after = "for _, db_fi := range db_files_list"
	for _, db_fi := range db_files_list {

		db_file_map[db_fi.FilePath] = &FileItem{
			Path:      db_fi.FilePath,
			FileState: FileState(db_fi.FileState),
			FileHash:  *db_fi.FileHash,
			stats: &FiFileInfo{
				name:    db_fi.FilePath,
				size:    int64(db_fi.FileSize),
				mode:    int(db_fi.FileMode),
				modtime: db_fi.ModTime.Time,
				is_dir:  false,
			},
		}
	}

	file_map := (files_map.FileMap)

	file_map.Lock()
	defer file_map.Unlock()
	after = "for short_filepath, f_i := range file_map.M"
	for short_filepath, f_i := range file_map.M {
		if has_deadline, near_deadline := ctx.NearDeadline(time.Second); has_deadline && near_deadline {
			return nil
		}
		switch f_i.FileState {
		case fstate_to_upload:
			fallthrough
		case fstate_upload_err:
			f_i.FileState = fstate_to_upload
			err1 := upload_a_file(ctx, f_i)
			if err1 != nil {
				Logger.LogErr(loc, err1)

			}
			continue
		}

		full_filepath := f_i.Full(storage_struct.StorageDirectory)

		after = fmt.Sprintf("db_fi, ok := db_file_map[short_filepath: %s])", full_filepath)
		db_fi, ok := db_file_map[short_filepath]
		if !ok {
			after = "err1 := upload_a_file(ctx, f_i): !ok"
			err1 := upload_a_file(ctx, f_i)
			if err1 != nil {
				Logger.LogErr(loc, err1)

			}
			continue
		}

		if db_fi.FileHash == f_i.FileHash {
			f_i.FileState = fstate_unchanged
			continue
		}

		after = fmt.Sprintf("stats, err2 := os.Stat(full_filepath: %s)", full_filepath)
		stats, err2 := os.Stat(full_filepath)
		if err2 != nil || stats == nil {
			if errors.Is(err2, os.ErrNotExist) || errors.Is(err2, os.ErrInvalid) {
				f_i.FileState = fstate_to_download
				continue
			}

			f_i.FileState = fstate_os_err
			Logger.LogErr(loc, err2)
			continue
		}

		after = "if db_fi.IsBefore(stats.ModTime()) "
		if db_fi.IsBefore(stats.ModTime()) {
			after = "err1 := upload_a_file(ctx, f_i): db_fi.IsBefore"
			err1 := upload_a_file(ctx, f_i)
			if err1 != nil {
				Logger.LogErr(loc, err1)

			}
			continue
		} else if db_fi.IsAfter(stats.ModTime()) {
			f_i.FileState = fstate_to_download
			continue
		}

	}

	return
}

func get_file_type(file_name string) string {

	stp_1 := strings.Split(file_name, ".")
	stp_2 := len(stp_1)
	stp_3 := stp_1[stp_2-1]
	if len(stp_3) > 4 {

		return "unknown"
	}
	return stp_3
}

func upload_a_file(ctx ftp_context.Context, fi *FileItem) error {
	loc := log_item.Loc(`upload_a_file(ctx ftp_context.Context, fi *FileItem) error`)

	db_conn := db.DBPool.GetConn()
	defer db.DBPool.Return(db_conn)
	// tickerf(&loc, 1, "before readfile", fi.Path)
	file_path := fi.Full(storage_struct.StorageDirectory)
	d, err1 := os.ReadFile(file_path)
	if err1 != nil {
		if errors.Is(err1, os.ErrNotExist) {
			fi.FileState = fstate_to_download
			return nil
		}
		fi.FileState = fstate_os_err
		return Logger.LogErr(loc, &log_item.LogItem{
			After:     fmt.Sprintf(`d, err1 := os.ReadFile(file_path: %s)`, file_path),
			CallStack: []error{err1},
		})
	}

	// log.Fatalln(d)
	// tickerf(&loc, 2, "before uploading file_data")
	mod_time := pgtype.Timestamptz{Time: fi.stats.ModTime(), Valid: true}
	added_data, err2 := DB.UploadFilesStepOneUploadData(ctx, db_conn,
		&db_access.
			UploadFilesStepOneUploadDataParams{
			ModTime:   mod_time,
			FileSize:  int32(fi.stats.Size()),
			FileDataB: d,
		})
	if err2 != nil {

		fi.FileState = fstate_upload_err
		return Logger.LogErr(loc, &log_item.LogItem{
			After:     `added_data, err2 := DB.UploadFilesStepOneUploadData(ctx, db_conn,`,
			CallStack: []error{err2},
		})
	}

	if len(fi.FileHash) > 0 {
		if fi.FileHash == *added_data.FileHash {
			fi.FileState = fstate_unchanged
			return nil
		}
	}
	// tickerf(&loc, 3, "before updating file_metadata")
	metadata_id, err3 := DB.UploadFilesStepTwoUploadMetadata(ctx, db_conn,
		&db_access.
			UploadFilesStepTwoUploadMetadataParams{
			FilePath:   fi.Path,
			FileType:   get_file_type(fi.Path),
			FileDataID: added_data.ID,
			FileMode:   int32(fi.stats.Mode()),
			ModTime:    mod_time,
		})

	if err3 != nil {
		fi.FileState = fstate_upload_err
		return Logger.LogErr(loc, &log_item.LogItem{
			After:     `metadata_id, err3 := DB.UploadFilesStepTwoUploadMetadata(ctx, db_conn,`,
			CallStack: []error{err3},
		})
	}

	file_state := fi.FileState
	if file_state != fstate_new {
		file_state = fstate_changed
	}
	// tickerf(&loc, 4, "before updating file_tracker")
	err4 := DB.UpdateFileTracker(ctx, db_conn, &db_access.
		UpdateFileTrackerParams{
		PeerID:        storage_struct.PeerId.Bytes,
		FileMetaID:    metadata_id,
		CurrentHashID: added_data.ID,
		FileState:     file_state.String(),
	})

	if err4 != nil {
		fi.FileState = fstate_upload_err
		return Logger.LogErr(loc, &log_item.LogItem{
			After:     `err4 := DB.UpdateFileTracker(ctx, db_conn, &db_access.`,
			CallStack: []error{err4},
		})
	}

	fi.FileState = fstate_uploaded
	// tickerf(&loc, 5, "successfully updated: a file", fi.Path, added_data.FileHash)
	return nil
}

func PermanentUploadLoop(ctx ftp_context.Context) error {
	defer ctx.Finished()
	loc := log_item.Loc(`PermanentUploadLoop(ctx ftp_context.Context) error`)

	files_map := NewFileMap()
	max_age := storage_struct.UploadDirs.MaxAgeInDaysBeforeDelete *
		int(time.Hour) *
		24
	db_conn := db.DBPool.GetConn()
	defer db.DBPool.Return(db_conn)
	for _, permanent_upload_dir := range storage_struct.UploadDirs.UploadDirs {

		err1 := WalkDir(permanent_upload_dir, files_map, &fsnotify.Watcher{})
		if err1 != nil {
			Logger.LogErr(loc, err1)
			<-time.After(time.Minute)
			continue
		}

		dir_id, err2 := DB.UploadStoreStepOnePeerDir(ctx, db_conn, &db_access.UploadStoreStepOnePeerDirParams{
			PeerID:  storage_struct.PeerId,
			DirPath: permanent_upload_dir,
		})
		if err2 != nil {
			Logger.LogErr(loc, err2)
			<-time.After(time.Minute)
			continue
		}

		files_map.l.Lock()
		for _, file_i := range files_map.FileMap.M {
			d, err3 := os.ReadFile(file_i.Full(permanent_upload_dir))
			if err3 != nil {
				Logger.LogErr(loc, err3)
				<-time.After(time.Second)
				continue
			}

			file_modtime := pgtype.Timestamptz{Time: file_i.stats.ModTime(), Valid: true}
			file_upload, err4 := DB.UploadStoreStepTwoUploadFile(ctx, db_conn, &db_access.
				UploadStoreStepTwoUploadFileParams{
				ModTime:   file_modtime,
				FileSize:  int32(file_i.stats.Size()),
				FileDataB: d,
			})
			if err4 != nil {
				Logger.LogErr(loc, err4)
				<-time.After(time.Second)
				continue
			}

			file_metadata_id, err5 := DB.UploadStoreStepThreeUpdateMetadata(ctx, db_conn, &db_access.
				UploadStoreStepThreeUpdateMetadataParams{
				FilePath:   file_i.Path,
				FileType:   get_file_type(file_i.Path),
				FileDataID: file_upload.ID,
				FileMode:   int32(file_i.stats.Mode()),
				ModTime:    file_modtime,
				DirID:      &dir_id,
			})

			if err5 != nil {
				Logger.LogErr(loc, err5)
				<-time.After(time.Second)
				continue
			}

			DB.UpdateFileTracker(ctx, db_conn, &db_access.
				UpdateFileTrackerParams{
				PeerID:        storage_struct.PeerId.Bytes,
				FileMetaID:    file_metadata_id,
				CurrentHashID: file_upload.ID,
				FileState:     "uploaded",
			})

			if !storage_struct.UploadDirs.DeleteOnUpload || max_age < 0 {
				continue
			}

			tds := time.Since(file_i.stats.ModTime())
			to_delete := tds >= time.Duration(tds)
			if !to_delete {
				continue
			}

			err6 := os.Remove(file_i.Full(permanent_upload_dir))
			if err6 != nil {
				Logger.LogErr(loc, err6)
				<-time.After(time.Millisecond)

			}

		}
		files_map.l.Unlock()
	}

	return nil
}
