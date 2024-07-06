package mainthread

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	ftp_base "github.com/it-shiloheye/ftp_system/v2/lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system/v2/lib/context"
	db "github.com/it-shiloheye/ftp_system/v2/lib/db_access"
	db_access "github.com/it-shiloheye/ftp_system/v2/lib/db_access/generated"

	"github.com/it-shiloheye/ftp_system/v2/lib/logging"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system/v2/peer/config"
	db_helpers "github.com/it-shiloheye/ftp_system/v2/peer/main_thread/db_access"
	"github.com/jackc/pgx/v5/pgtype"
)

var Logger = logging.Logger
var ServerConfig = server_config.ServerConfig
var DB = db_access.New()

type ToUploadType = map[string]*FilesListItem
type ToUpdateType = map[string]*db_access.GetFilesRow

func Loop(ctx ftp_context.Context, StorageStruct *server_config.StorageStruct) error {
	loc := log_item.Loc(`mainthread.Loop(ctx ftp_context.Context) error`)
	defer ctx.Finished()

	tc := time.NewTicker(time.Minute * 5)

	for ok := true; ok; {

		Logger.Logf(loc, "new loop: %s", fmt.Sprintln(time.Now()))

		tmp_fileslist, err1 := WalkDir(ctx, StorageStruct)
		if err1 != nil {
			Logger.LogErr(loc, err1)
			<-time.After(time.Minute)
			continue
		}

		db_tmp_fileslist, err2 := db_helpers.GetFiles(ctx, StorageStruct)
		if err2 != nil {
			Logger.LogErr(loc, err2)
			<-time.After(time.Minute)
			continue
		}

		log.Println("reaching here")
		/*uniq_to_download*/ _, uniq_to_upload, err3 := process_to_upload_func(tmp_fileslist, db_tmp_fileslist)
		if err3 != nil {
			Logger.LogErr(loc, err3)
			<-time.After(time.Minute)
			continue
		}

		insert_file_rows, _ := upload_func(ctx, uniq_to_upload, StorageStruct)

		for _, f := range insert_file_rows {
			log.Println("uploaded:", f.ID, f.FileHash)
		}

		switch StorageStruct.PeerRole {
		case db_access.PeerRoleTypeStorage:
			err4 := to_download_func(ctx, StorageStruct)
			if err4 != nil {
				Logger.LogErr(loc, err4)
				<-time.After(time.Minute)
				continue
			}

		}

		// download_func(ctx, uniq_to_download, StorageStruct)
		select {
		case <-ctx.Done():
			return Logger.LogErr(loc, ctx.Err())
		case <-tc.C:

			set_up_storagestruct(ctx, StorageStruct)
		}
	}

	return nil
}

type FilesListItem struct {
	Path string
	FD   os.FileInfo
	*os.File
}

func (fsi *FilesListItem) Reopen() error {
	loc := log_item.Locf(`func (fsi *FilesListItem) Reopen(%s) error `, fsi.Path)
	var err5, err6 error
	fsi.File, err5 = os.OpenFile(fsi.Path, os.O_RDONLY, ftp_base.FS_MODE)
	if err5 != nil {
		return Logger.LogErr(loc, err5)
	}

	fsi.FD, err6 = fsi.File.Stat()
	if err6 != nil {
		return Logger.LogErr(loc, err6)
	}

	return nil
}

func upload_func(ctx ftp_context.Context, uniq_to_upload ToUploadType, storage_struct *server_config.StorageStruct) (insert_file_rows []*db_access.InsertFileRow, to_delete []string) {
	loc := log_item.Loc(`upload_func(uniq_to_upload map[string]*FilesListItem) (insert_file_rows []*db_access.InsertFileRow)`)
	conn := db.DBPool.GetConn()
	defer db.DBPool.Return(conn)

	for file_path, FD := range uniq_to_upload {
		d, err4 := os.ReadFile(file_path)
		if err4 != nil {
			Logger.LogErr(loc, err4)
			continue
		}

		name_01 := strings.Split(file_path, string(os.PathSeparator))
		name_0l := len(name_01) - 1
		name := name_01[name_0l]
		insert_file_row, err7 := DB.InsertFile(ctx.Add(), conn, &db_access.InsertFileParams{
			PeerID:           storage_struct.PeerId.Bytes,
			FilePath:         file_path,
			FileType:         filepath.Ext(file_path),
			FileName:         name,
			ModificationDate: pgtype.Timestamp{Time: FD.FD.ModTime(), Valid: true},
			FileState:        db_access.NullFileStatusType{FileStatusType: db_access.FileStatusTypeNew},
			FileData:         d,
		})

		if err7 != nil {
			if strings.Contains(err7.Error(), `duplicate key value violates unique constraint "file_storage_file_hash_key"`) {
				delete(uniq_to_upload, file_path)
				log.Println("exists: ", file_path)
				if storage_struct.OnUpload.DeleteOnUpload {
					if storage_struct.OnUpload.MaxAgeInDaysBeforeDelete < 0 {
						continue
					}

					if time.Since(FD.FD.ModTime()) < time.Hour*24*time.Duration(storage_struct.OnUpload.MaxAgeInDaysBeforeDelete) {
						continue
					}

				}

				continue
			}
			Logger.LogErr(loc, err7)
			continue
		}

		insert_file_rows = append(insert_file_rows, insert_file_row)
	}

	return
}

func process_to_upload_func(
	tmp_fileslist []*FilesListItem,
	db_tmp_fileslist []*db_access.GetFilesRow) (
	unique_to_update ToUpdateType,
	uniq_to_upload ToUploadType,
	err error) {
	loc := log_item.Loc(` download_func(ctx ftp_context.Context, uniq_to_download ToDownloadType, storage_struct *server_config.StorageStruct) error`)

	half := len(tmp_fileslist) / 2
	uniq_to_upload = make(map[string]*FilesListItem, half)
	unique_to_update = make(map[string]*db_access.GetFilesRow, half)

	// tmp_fileslist from WalkDir gives a list of all the files on this machine
	for _, file_item := range tmp_fileslist {
		uniq_to_upload[file_item.Path] = file_item
	}

	// "get files from db" returns a list of all files for this peer
	// use a map to rerun
	for _, file_uploaded := range db_tmp_fileslist {
		fd_tmp, ok := uniq_to_upload[file_uploaded.FilePath]
		if !ok {
			continue
		}

		stats, err1 := fd_tmp.Stat()
		if err1 != nil {
			Logger.LogErr(loc, err1)
			<-time.After(time.Microsecond * 20)
			continue
		}

		if stats.ModTime() == file_uploaded.ModificationDate.Time ||
			stats.ModTime().Before(file_uploaded.ModificationDate.Time) {
			// if there have been no changes to the file since, do nothing
			delete(uniq_to_upload, file_uploaded.FilePath)
			continue
		} else {
			// upload the file, it will generate a new hash
			// update the file, mark the new file prev_hash as the old file
			unique_to_update[file_uploaded.FilePath] = file_uploaded
		}

	}

	return
}

func delete_func(ctx ftp_context.Context, to_delete *[]string, storage_struct *server_config.StorageStruct) error {

	return nil
}

func to_download_func(
	ctx ftp_context.Context,
	storage_struct *server_config.StorageStruct) error {
	loc := log_item.Locf(`func to_download_func(ctx ftp_context.Context, uniq_to_update ToUpdateType, storage_struct *server_config.StorageStruct) error `)
	storage_struct.RLock()
	defer storage_struct.RUnlock()

	conn := db.DBPool.GetConn()
	defer db.DBPool.Return(conn)
	path_sep := func() string {
		return string(os.PathSeparator)
	}

	download_files := map[string]*db_access.GetFilesRow{}
	uniq_file := map[string]os.FileInfo{}
	create_path := func(str ...string) string {
		return strings.Join(str, path_sep())
	}

	for _, peer_sub := range storage_struct.SubscribedPeers {

		files, err1 := DB.GetFiles(ctx, conn, peer_sub.PeerUuid.Bytes)
		if err1 != nil {
			Logger.LogErr(loc, err1)
			continue
		}

		if len(files) == 0 {
			continue
		}

		filepath.Walk(peer_sub.LocalName, func(fpath string, info fs.FileInfo, err error) error {
			loc := log_item.Locf(`filepath.Walk(peer_sub.LocalName: %s,func(path: %s, info fs.FileInfo, err error) error`, peer_sub.LocalName, fpath)
			if err != nil {
				Logger.LogErr(loc, err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			uniq_file[fpath] = info

			return nil
		})

		for _, v := range files {
			tmp_tm := v.ModificationDate.Time
			str_tm := strings.Split(fmt.Sprint(tmp_tm), " ")
			dir_01 := strings.Split(str_tm[0], "-")
			dir_02 := create_path(storage_struct.StorageDirectory, "storage", peer_sub.LocalName, strings.Join(dir_01, path_sep()))

			file_hash := *v.FileHash
			file_01 := create_path(dir_02, file_hash[:6]+"_"+v.FileName)
			_, ok := uniq_file[file_01]

			if ok {
				continue
			}

			os.MkdirAll(dir_02, ftp_base.FS_MODE)

			download_files[file_01] = v

		}
	}

	for fpath, file_row := range download_files {

		f_ds, err1 := DB.GetFileData(ctx, conn, file_row.FileHash)
		if err1 != nil {
			Logger.LogErr(loc, err1)
			continue
		}
		if len(f_ds) < 1 {
			Logger.LogErr(loc, fmt.Errorf("file: %s file_hash: %s missing from Database", file_row.FilePath, *file_row.FileHash))
			continue
		}

		file_d := f_ds[0]

		err2 := os.WriteFile(fpath, file_d.FileData, ftp_base.FS_MODE)
		if err2 != nil {
			Logger.LogErr(loc, err2)
			continue
		}

	}

	return nil
}
