package api

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	ftp_base "github.com/it-shiloheye/ftp_system/v3/lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

func RegisterApi(ctx ftp_context.Context, r *gin.Engine) {
	loc := log_item.Loc(`func RegisterApi(r *gin.Engine)`)
	// cors_ := cors.NewCors()

	_ = loc

	api_route := r.Group("api")

	limit_requests := map[string]chan struct{}{}

	DBPool.PopulateConns(ctx, 5)
	api_route.Use(func(ctx *gin.Context) {
		// cors_.RenderCors(ctx)
		ti := time.Now()
		remote_ip, req_method, req_url, req_status := ctx.RemoteIP(), ctx.Request.Method, ctx.Request.RequestURI, ""
		limiter, ok := limit_requests[remote_ip]
		if !ok {
			limit_requests[remote_ip] = make(chan struct{}, 1)
			limiter, _ = limit_requests[remote_ip]
			limiter <- struct{}{}
		}
		<-limiter
		defer func() {
			limiter <- struct{}{}
		}()
		ctx.Next()

		if ctx.Request.Response != nil {
			req_status = ctx.Request.Response.Status
		}

		Logger.Log(&log_item.LogItem{
			Message: fmt.Sprintf("%s | %s | %s | %s | %s", remote_ip, req_method, req_url, req_status, time.Since(ti).String()),
		})
	})

	api_route.GET("/files_list", get_files_list_route)

	print_route(r)
}

func print_route(r *gin.Engine) {
	loc := log_item.Loc(`func print_route(r *gin.Engine)`)
	r_info := r.Routes()

	paths_list := map[string][]*HandlerRoutesJSON{}
	for _, path_i := range r_info {

		t, ok := paths_list[path_i.Path]
		if !ok {
			paths_list[path_i.Path] = []*HandlerRoutesJSON{
				{
					Path:   path_i.Path,
					Method: path_i.Method,
				},
			}
		} else {
			paths_list[path_i.Path] = append(t, &HandlerRoutesJSON{
				Path:   path_i.Path,
				Method: path_i.Method,
			})
		}

	}

	r_info_json, err1 := json.MarshalIndent(&paths_list, " ", "\t")
	if err1 != nil {
		Logger.LogErr(loc, err1)
		return
	}

	err2 := os.WriteFile("./api_routes.json", r_info_json, ftp_base.FS_MODE)
	if err2 != nil {
		Logger.LogErr(loc, err2)
		return
	}
}
