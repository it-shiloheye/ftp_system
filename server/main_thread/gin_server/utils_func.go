package ginserver

import (
	"encoding/json"
	"io"
	"io/fs"

	"time"

	"github.com/gin-gonic/gin"
	server_dirhandler "github.com/it-shiloheye/ftp_system/server/main_thread/dir_handler"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	"github.com/it-shiloheye/ftp_system_lib/logging"
	"github.com/it-shiloheye/ftp_system_lib/logging/log_item"
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

var filemeta_data = &server_dirhandler.FileStorage.FileMetaData
var uploaded_hashes = &server_dirhandler.FileStorage.UploadedHashes
var tmp_filedata = &server_dirhandler.TmpFileData

var file_hash_chan = make(chan string, 100)

func RouteGuard(route string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		loc := log_item.Locf(`r.Group("/%s", func(ctx *gin.Context)`, route)

		// log.Printf("\n\ncookies:\n %v\n\n\n", ctx.Request.Cookies())
		client_id, err1 := ctx.Cookie("client-id")
		if err1 != nil {
			ctx.JSON(400, gin.H{
				"state": failure,
			})
			Logger.LogErr(loc, &log_item.LogItem{
				Location: loc,
				After:    `client_id, err1 := ctx.Cookie("client-id")`,
				Body: map[string]any{
					"ip_addr": ctx.RemoteIP(),
					"url":     ctx.Request.URL,
				},
				Level:     log_item.LogLevelError02,
				Message:   "failed to read data",
				CallStack: []error{err1},
			})
			return
		}
		ctx.Next()
		Logger.Logf(loc, "request from: %s\t%s", client_id, ctx.RemoteIP())
	}
}

func read_request(ctx *gin.Context, js_obj any) (d []byte, err error) {
	loc := log_item.Loc("read_request(ctx *gin.Context, js_obj any) (d []byte, err error)")

	if ctx.Request.Body == nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			Level:    log_item.LogLevelError02,
			Message:  "nil request body",
		}
		return
	}

	d, err1 := io.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if err1 != nil {
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			Level:    log_item.LogLevelError02,
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
		err = &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			Level:    log_item.LogLevelError02,
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

	Logger.Logf(log_item.Loc(loc), "received:%s\n%016d bytes", ctx.Request.RequestURI, len(d))

	return
}
