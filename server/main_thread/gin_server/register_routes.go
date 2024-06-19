package ginserver

import (
	"errors"
	"fmt"
	"time"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/it-shiloheye/ftp_system_lib/logging/log_item"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

func RegisterRoutes(r *gin.Engine) error {
	loc := log_item.Loc("RegisterRoutes(r *gin.Engine) error")

	client_group := r.Group("/client", RouteGuard("client"))

	upload_group := r.Group("/upload", RouteGuard("upload"))

	tmp_dir := ServerConfig.TmpDirectory
	if len(tmp_dir) < 2 {
		log.Fatalln(&log_item.LogItem{
			Location: loc,
			After:    fmt.Sprintf(`tmp_dir := ServerConfig.DirConfig.Path // %s`, tmp_dir),
			Level:    log_item.LogLevelError02,

			Message: "no temporary directory given",
		})
	}

	err3 := os.MkdirAll(tmp_dir, fs_mode)
	if err3 != nil && !errors.Is(err3, os.ErrExist) {
		log.Fatalln(&log_item.LogItem{
			Location:  loc,
			After:     fmt.Sprintf(`err3 := os.MkdirAll(tmp_dir:"%s", fs.FileMode(ftp_base.S_IRWXU|ftp_base.S_IRWXO))`, tmp_dir),
			CallStack: []error{err3},
		})

	}

	client_group.POST("/register", func(ctx *gin.Context) {

		ctx.JSON(400, gin.H{
			"state": failure,
		})
	})

	// sends entire filetree from client to server
	upload_group.POST("/confirm", ConfirmFileTree)

	upload_group.POST("/bulk", UploadBulkFiles)

	return nil
}

func StoreUploadedFiles(ctx ftp_context.Context, storage_path string) error {
	loc := log_item.Locf(`StoreUploadedFiles(ctx ftp_context.Context, storage_path: %s) error`, storage_path)
	defer ctx.Finished()
	tc := time.NewTicker(time.Minute * 5)
	hash_list := []string{}

	generate_dir := func(name string) (dir string, file_name string) {
		dir_a, dir_b, f_name := name[:2], name[2:5], name
		dir = dir_a + string(os.PathSeparator) + dir_b
		file_name = dir + string(os.PathSeparator) + f_name
		return
	}

	var s string
	for ok := true; ok; {
		select {
		case _, ok = <-ctx.Done():
			break
		case s = <-file_hash_chan:
			if len(s) > 0 {
				hash_list = append(hash_list, s)
			}
			continue
		case <-tc.C:
			break
		}

	hashlist_loop:
		for _, file_hash := range hash_list {
			<-time.After(time.Millisecond * 1000)
			if state, ok := uploaded_hashes.Get(file_hash); ok && state == saved {
				continue
			}

			log.Println(`saving: `, file_hash)
			dir, file_name := generate_dir(file_hash)
			os.MkdirAll(dir, fs_mode)

			dir_entry, _ := os.ReadDir(dir)

			if len(dir_entry) < 1 {
				d_, ok := tmp_filedata.Get(file_hash)
				if !ok {
					uploaded_hashes.Set(file_hash, missing)
					file_hash_chan <- file_hash
					continue hashlist_loop
				}
				err1 := os.WriteFile(file_name, d_, fs_mode)
				if err1 != nil {

					Logger.LogErr(loc, err1)
					file_hash_chan <- file_hash
					continue hashlist_loop
				}
				uploaded_hashes.Set(file_hash, saved)
				tmp_filedata.Delete(file_hash)
				continue hashlist_loop
			}

			for _, f_d := range dir_entry {
				if f_d.IsDir() {
					continue
				}
				if f_d.Name() == file_name {

					uploaded_hashes.Set(file_hash, saved)
					continue hashlist_loop
				}

			}

			d_, ok := tmp_filedata.Get(file_hash)
			if !ok {
				uploaded_hashes.Set(file_hash, missing)
				file_hash_chan <- file_hash
				continue hashlist_loop
			}
			err1 := os.WriteFile(file_name, d_, fs_mode)
			if err1 != nil {

				Logger.LogErr(loc, err1)
				file_hash_chan <- file_hash
				continue hashlist_loop
			}
			uploaded_hashes.Set(file_hash, saved)
			tmp_filedata.Delete(file_hash)
			continue hashlist_loop

		}

		clear(hash_list)
	}

	return nil
}
