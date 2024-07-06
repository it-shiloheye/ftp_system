package mainthread

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	ftp_base "github.com/it-shiloheye/ftp_system/v3/lib/base"

	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	db "github.com/it-shiloheye/ftp_system/v3/lib/db_access"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

const file_hash_err = `error during:
SELECT
    id, file_hash, mod_time, file_size, file_data_b, creation_time
FROM
    file_data
WHERE
    file_data.file_hash = %s;\n`

func log_error_file_hash(file_i *FileItem, loc log_item.Loc, err error) {
	if err == nil || file_i == nil {
		return
	}
	Logger.LogErr(loc, err)
	log.Printf(file_hash_err, file_i.FileHash)
}

const file_path_err = `error during:
SELECT * from file_data 
JOIN file_metadata ON file_metadata.file_data_id = file_data.id
where file_metadata.file_path = %s;\n`

func log_error_file_path(file_i *FileItem, loc log_item.Loc, err error) {
	if err == nil || file_i == nil {
		return
	}
	Logger.LogErr(loc, err)
	log.Printf(file_path_err, file_i.Path)
}

func ClientDownloadFunc(ctx ftp_context.Context, files_map *FileMapType) (err error) {
	loc := log_item.Loc(`ClientDownloadFunc(ctx ftp_context.Context, files_map *FileMapType) (err error)`)
	log.Println(`ClientDownloadFunc(ctx ftp_context.Context, files_map *FileMapType) (err error)`)
	after := "setup"
	defer recover_func(&after, loc)
	if has_deadline, near_deadline := ctx.NearDeadline(time.Second); has_deadline && near_deadline {
		return nil
	}
	file_map := (files_map.FileMap)
	db_conn := db.DBPool.GetConn()
	after = "db_files_list, err1 := DB.GetFilesList(ctx, db_conn)"
	db_files_list, err1 := DB.GetFilesList(ctx, db_conn)
	defer db.DBPool.Return(db_conn)

	if err1 != nil {
		if strings.Contains(err1.Error(), "no rows in result set") {
			log.Println("no rows in result set")
			return nil
		}
		err = Logger.LogErr(loc, &log_item.LogItem{
			After:     after_db_conn,
			CallStack: []error{err1},
		})
		return
	}

	after = "for _, db_fi := range db_files_list"
	for _, db_fi := range db_files_list {
		if has_deadline, near_deadline := ctx.NearDeadline(time.Second); has_deadline && near_deadline {
			return nil
		}
		log.Println("db_fi: ", db_fi.FilePath)
		tmp := &FileItem{
			Path:     db_fi.FilePath,
			FileHash: *db_fi.FileHash,
			stats: &FiFileInfo{
				name:    db_fi.FilePath,
				size:    int64(db_fi.FileSize),
				mode:    int(db_fi.FileMode),
				modtime: db_fi.ModTime.Time,
				is_dir:  false,
			},
		}

		fi, ok := file_map.Get(db_fi.FilePath)
		if !ok {
			after = fmt.Sprintf(`"fi, ok := file_map.Get(db_fi.FilePath: %s) && !ok`, db_fi.FilePath)
			if err2 := download_file(ctx, tmp); err2 != nil {
				Logger.LogErr(loc, err2)
				continue
			}

			tmp.FileState = fstate_downloaded
			file_map.Set(tmp.Path, tmp)
			continue
		}

		if fi.IsBefore(tmp.stats.ModTime()) {
			after = fmt.Sprintf(` fi.IsBefore(tmp.stats.ModTime()); FilePath: %s`, fi.Path)
			if err2 := download_file(ctx, tmp); err2 != nil {
				Logger.LogErr(loc, err2)
				continue
			}

			tmp.FileState = fstate_downloaded
			file_map.Set(tmp.Path, tmp)
			continue
		}

		if fi._download_err() || fi._to_download() || fi._os_err() {
			after = fmt.Sprintf(`if fi._download_err() || fi._to_download() || fi._os_err(); FilePath: %s`, fi.Path)
			if err2 := download_file(ctx, tmp); err2 != nil {
				Logger.LogErr(loc, err2)
				continue
			}

			tmp.FileState = fstate_downloaded
			file_map.Set(tmp.Path, tmp)
			continue
		}
	}

	return nil
}

func download_file(ctx ftp_context.Context, db_fi *FileItem) error {
	loc := log_item.Locf(`download_file(ctx ftp_context.Context,db_fi *FileItem: %s) error`, db_fi.FileHash)

	file_path := db_fi.Full(storage_struct.StorageDirectory)
	log.Println("downloading: ", file_path)
	db_conn := db.DBPool.GetConn()
	defer db.DBPool.Return(db_conn)
	db_data_res, err1 := DB.DownloadFileStepOneGetLatestData(ctx, db_conn, &db_fi.FileHash)
	if err1 != nil {
		db_fi.FileState = fstate_download_err
		return Logger.LogErr(loc, err1)
	}

	for {

		err2 := os.WriteFile(file_path, db_data_res.FileDataB, ftp_base.FS_MODE)
		if err2 != nil {
			if errors.Is(err2, os.ErrNotExist) || errors.Is(err2, os.ErrInvalid) {
				err3 := MkDir(file_path)
				if err3 != nil {

					db_fi.FileState = fstate_os_err
					return Logger.LogErr(loc, err3)
				}
				continue
			}
			db_fi.FileState = fstate_download_err
			return Logger.LogErr(loc, err2)
		}

		break
	}

	err3 := os.Chtimes(file_path, time.Time{}, db_data_res.ModTime.Time)
	if err3 != nil {
		db_fi.FileState = fstate_download_err
		return Logger.LogErr(loc, err3)
	}

	log.Println("successfully downloaded: ", file_path)
	return nil
}

func MkDir(file_path string) error {
	loc := log_item.Locf(`MkDir(file_path: %s) error`, file_path)
	dir_path_s := strings.Split(file_path, string(os.PathSeparator))
	l_d := len(dir_path_s) - 1
	dir_p := strings.Join(dir_path_s[:l_d], string(os.PathSeparator))
	err2 := os.MkdirAll(dir_p, ftp_base.FS_MODE)
	if err2 == nil || errors.Is(err2, os.ErrExist) {
		return nil
	}
	return Logger.LogErr(loc, err2)
}

func MkDirStorage(peer_id uuid.UUID, file_i *FileItem) (file_path string) {

	f_p_s := strings.Split(file_i.Path, "/")
	last_p_s := strings.Split(f_p_s[len(f_p_s)-1], ".")
	name := last_p_s[0]

	a, b, c := file_i.FileHash[:2], file_i.FileHash[2:4], name+"-"+file_i.FileHash

	dir_p := strings.Join([]string{storage_struct.StorageDirectory, "store", peer_id.String(), a, b}, "/")
	os.MkdirAll(dir_p, ftp_base.FS_MODE)

	return dir_p + "/" + c
}
