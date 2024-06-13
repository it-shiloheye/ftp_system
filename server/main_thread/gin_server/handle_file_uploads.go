package ginserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"io/fs"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/it-shiloheye/ftp_system/server/main_thread/dir_handler"
	"github.com/it-shiloheye/ftp_system/server/main_thread/logging"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

const (
	fs_mode  = fs.FileMode(ftp_base.S_IRWXU | ftp_base.S_IRWXO)
	failure  = "failure"
	success  = "success"
	missing  = "missing"
	uploaded = "uploaded"
	reupload = "reupload"
	saved    = "saved"
)

var Logger = logging.Logger

var filemeta_data = server_dirhandler.FileStorage.FileMetaData
var uploaded_hashes = server_dirhandler.FileStorage.UploadedHashes
var tmp_filedata = server_dirhandler.TmpFileData

var file_hash_chan = make(chan string, 100)

func RegisterRoutes(r *gin.Engine) error {
	loc := "RegisterRoutes(r *gin.Engine) error"
	upload_group := r.Group("/upload", func(ctx *gin.Context) {
		log.Println(ctx.Request.URL.RequestURI())
		ctx.Next()
	})

	tmp_dir := ServerConfig.TmpDirectory
	if len(tmp_dir) < 2 {
		log.Fatalln(&ftp_context.LogItem{
			Location: loc,
			After:    fmt.Sprintf(`tmp_dir := ServerConfig.DirConfig.Path // %s`, tmp_dir),
			Err:      true,

			Message: "no temporary directory given",
		})
	}

	err3 := os.MkdirAll(tmp_dir, fs_mode)
	if err3 != nil && !errors.Is(err3, os.ErrExist) {
		log.Fatalln(&ftp_context.LogItem{
			Location:  loc,
			After:     fmt.Sprintf(`err3 := os.MkdirAll(tmp_dir:"%s", fs.FileMode(ftp_base.S_IRWXU|ftp_base.S_IRWXO))`, tmp_dir),
			CallStack: []error{err3},
		})

	}

	upload_group.GET("/confirm/:file_hash", func(ctx *gin.Context) {
		file_hash := ctx.Param("file_hash")

		state, ok := uploaded_hashes.Get(file_hash)

		if !ok || state != saved {
			ctx.JSON(404, gin.H{
				"state": missing,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"state": uploaded,
		})

	})
	upload_group.POST("/file/:file_hash", func(ctx *gin.Context) {
		file_hash := ctx.Param("file_hash")

		d, err1 := io.ReadAll(ctx.Request.Body)
		if err1 != nil {
			log.Println(err1)
			ctx.AbortWithStatusJSON(422, gin.H{
				"do": "better",
			})
			return
		}

		fh := &filehandler.FileHash{}
		err2 := json.Unmarshal(d, fh)
		if err2 != nil {
			log.Println(err2)
			ctx.AbortWithStatusJSON(500, gin.H{
				"my":    "bad",
				"error": err2.Error(),
			})
			return
		}

		ctx.JSON(201, gin.H{
			"received": file_hash,
		})

		filemeta_data.Set(file_hash, fh)
		log.Println(string(d))
	})

	upload_group.POST("/stream/:file_hash", func(ctx *gin.Context) {
		hash := ctx.Param("file_hash")
		loc := logging.Loc(fmt.Sprintf(`upload_group.POST("/stream/%s", func(ctx *gin.Context)`, hash))

		client_id, err0 := ctx.Cookie("client-id")
		if err0 != nil {
			ctx.JSON(400, gin.H{
				"state": failure,
			})
			err := &ftp_context.LogItem{
				Location: string(loc),
				After:    `client_id, err0 := ctx.Cookie("client-id")`,
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
					"url":     ctx.Request.URL,
				},
				Message:   err0.Error(),
				Err:       true,
				CallStack: []error{err0},
			}
			ctx.JSON(400, gin.H{
				"state": failure,
			})

			Logger.LogErr((loc), err)
			return
		}

		tmp_d := map[string]string{}

		d, err1 := read_request(ctx, &tmp_d)
		if err1 != nil {
			if len(d) > 999 {
				d = d[:999]
			}
			err := &ftp_context.LogItem{
				Location: string(loc),
				After:    `d, err1 := read_request(ctx, &tmp_d)`,
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
					"url":     ctx.Request.URL,
					"data":    string(d),
				},
				Message:   err1.Error(),
				Err:       true,
				CallStack: []error{err1},
			}
			ctx.JSON(400, gin.H{
				"state": failure,
			})

			Logger.LogErr((loc), err)
			return
		}

		file_d, err2 := base64.StdEncoding.DecodeString(tmp_d["data"])
		if err2 != nil {
			err := &ftp_context.LogItem{
				Location: string(loc),
				After:    `file_d, err2 := base64.StdEncoding.DecodeString(tmp_d["data"])`,
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
					"url":     ctx.Request.URL,
				},
				Err:       true,
				Message:   err2.Error(),
				CallStack: []error{err2},
			}

			ctx.JSON(400, gin.H{
				"state": failure,
			})
			Logger.LogErr((loc), err)
			return
		}

		fh, ok := filemeta_data.Get(hash)
		if !ok {
			ctx.JSON(400, gin.H{
				"state": failure,
			})
			return
		}

		bts := filehandler.NewBytesStore()

		_, err3 := bts.ReadFrom(bytes.NewBuffer(file_d))
		if err3 != nil {
			err := &ftp_context.LogItem{
				Location: string(loc),
				After:    `_, err3:= bts.ReadFrom(bytes.NewBuffer(file_d))`,
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
					"url":     ctx.Request.URL,
				},
				Err:       true,
				Message:   "failed to read data",
				CallStack: []error{err3},
			}

			ctx.JSON(400, gin.H{
				"state": failure,
			})
			Logger.LogErr((loc), err)
			return
		}

		file_hash, err4 := bts.Hash()
		if err4 != nil {
			err := &ftp_context.LogItem{
				Location: string(loc),
				After:    `file_hash, err4 := bts.Hash()`,
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
					"url":     ctx.Request.URL,
				},
				Err:       true,
				Message:   "failed to read data",
				CallStack: []error{err4},
			}

			ctx.JSON(400, gin.H{
				"state": failure,
			})
			Logger.LogErr((loc), err)
			return
		}

		if file_hash != hash {

			err := &ftp_context.LogItem{
				Location: string(loc),
				After:    ``,
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
					"url":     ctx.Request.URL,
				},
				Err:     true,
				Message: fmt.Sprintf(`data and hash doesn't match:\n(received data hash:) %s != %s (:expected hash)`, file_hash, hash),
			}

			ctx.JSON(400, gin.H{
				"state": failure,
			})
			Logger.LogErr((loc), err)
			return
		}

		curr_metadata := fmt.Sprintf("%s;%s;%s", &client_id, hash, fmt.Sprint(time.Now()))
		if fh.MetaData == nil {
			// metadata doesn't exist
			fh.MetaData = make(map[string]any)
			fh.MetaData["latest"] = curr_metadata
		} else {
			// metadata exists
			prev, exi := fh.MetaData["latest"]
			if !exi {
				// metadata.latest doesn't exist
				fh.MetaData["latest"] = curr_metadata

			} else {
				existing, ok := fh.MetaData["prev"]
				if !ok {
					// metadata.prev doesn't exist
					fh.MetaData["prev"] = []string{prev.(string)}
				} else {
					// update metadata.prev
					fh.MetaData["prev"] = append(existing.([]string), prev.(string))
				}

				fh.MetaData["latest"] = curr_metadata
			}

		}

		uploaded_hashes.Set(file_hash, uploaded)
		tmp_filedata.Set(file_hash, file_d)

		file_hash_chan <- hash
		ctx.JSON(201, gin.H{
			"state": success,
		})
		return
	})

	return nil
}

func read_request(ctx *gin.Context, js_obj any) (d []byte, err error) {
	loc := "read_request(ctx *gin.Context, js_obj any) (d []byte, err error)"

	if ctx.Request.Body == nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			Err:      true,
			Message:  "nil request body",
		}
		return
	}

	d, err1 := io.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if err1 != nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			Err:      true,
			After:    `d, err1 := io.ReadAll(ctx.Request.Body)`,
			Body: map[string]any{
				"ip_addr": ctx.RemoteIP(),
			},
			Message:   err1.Error(),
			CallStack: []error{err1},
		}
		return
	}

	err2 := json.Unmarshal(d, js_obj)
	if err2 != nil {
		if len(d) > 999 {
			d = d[:999]
		}
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			Err:      true,
			After:    "err2 := json.Unmarshal(d, js_obj)",
			Body: map[string]any{
				"ip_addr": ctx.RemoteIP(),
				"data":    d,
			},
			Message:   err2.Error(),
			CallStack: []error{err2},
		}

		return
	}

	Logger.Logf(logging.Loc(loc), "received:%s\n%016d bytes", ctx.Request.RequestURI, len(d))

	return
}

func StoreUploadedFiles(ctx ftp_context.Context, storage_path string) error {
	loc := logging.Locf(`StoreUploadedFiles(ctx ftp_context.Context, storage_path: %s) error`, storage_path)
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

					Logger.LogErr(loc, &ftp_context.LogItem{
						Err:       true,
						After:     fmt.Sprintf(`err1 := os.WriteFile(file_name: %s, d_, fs_mode)`, file_name),
						Message:   err1.Error(),
						CallStack: []error{err1},
					})
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

				Logger.LogErr(loc, &ftp_context.LogItem{
					Err:       true,
					After:     fmt.Sprintf(`err1 := os.WriteFile(file_name: %s, d_, fs_mode)`, file_name),
					Message:   err1.Error(),
					CallStack: []error{err1},
				})
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
