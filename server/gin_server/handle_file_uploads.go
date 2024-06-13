package ginserver

import (
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
	"github.com/it-shiloheye/ftp_system/server/main_thread/logging"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

const (
	fs_mode = fs.FileMode(ftp_base.S_IRWXU | ftp_base.S_IRWXO)
	failure = "failure"
	success = "success"
)

var Logger = logging.Logger

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

		tmp_file_hold.Set(file_hash, fh)
		log.Println(string(d))
	})

	upload_group.POST("/stream/:file_hash", func(ctx *gin.Context) {
		hash := ctx.Param("file_hash")
		loc := logging.Loc(fmt.Sprintf(`upload_group.POST("/stream/%s", func(ctx *gin.Context)`, hash))

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
					"data":    string(d),
				},
				Message:   err1.Error(),
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
				},
				Message:   err2.Error(),
				CallStack: []error{err2},
			}

			ctx.JSON(400, gin.H{
				"state": failure,
			})
			Logger.LogErr((loc), err)
			return
		}

		fh, ok := tmp_file_hold.Get(hash)
		if !ok {
			ctx.JSON(400, gin.H{
				"state": failure,
			})
			return
		}
		if fh.Hash != hash {

			err := &ftp_context.LogItem{
				Location: string(loc),
				After:    fmt.Sprintf(`fh, ok := tmp_file_hold.Get(hash:"%s")`, hash),
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
				},
				Message:   "data and hash doesn't match",
				CallStack: []error{err2},
			}

			ctx.JSON(400, gin.H{
				"state": failure,
			})
			Logger.LogErr((loc), err)
			return
		}

		if fh.MetaData == nil {
			fh.MetaData = make(map[string]any)
		}

		if fh.Hash == hash {

			fh.MetaData["contents"] = file_d
		}
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
	defer ctx.Finished()
	tc := time.NewTicker(time.Minute * 5)
	hash_list := []string{}

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

		for _, file_hash := range hash_list {
			log.Println(`saving: `, file_hash)
			<-time.After(time.Millisecond * 1000)
		}
	}

	return nil
}
