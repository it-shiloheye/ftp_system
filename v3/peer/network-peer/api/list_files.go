package api

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	db "github.com/it-shiloheye/ftp_system/v3/lib/db_access"
	db_access "github.com/it-shiloheye/ftp_system/v3/lib/db_access/generated"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system/v3/peer/config"
	"github.com/it-shiloheye/ftp_system/v3/peer/mainthread/db_helpers"
)

var Logger = logging.Logger

var DBPool = &db.DBPoolStruct{}
var DB = db.DB
var storage_struct = server_config.Storage

const (
	success string = "success"
	failure string = "failure"
)

type FilesListResponsType struct {
	FilePath string    `json:"path"`
	FileType string    `json:"type"`
	FileHash *string   `json:"hash"`
	ModTime  time.Time `json:"mod_time"`
}

type HandlerRoutesJSON struct {
	Path   string `json:"uri"`
	Method string `json:"method"`
}

var files_list_cache = &FilesListCache{}

type SearchCache struct {
	plain []*FilesListResponsType
	list  map[string][]*FilesListResponsType
}

type FilesListCache struct {
	sync.RWMutex
	search_cache SearchCache
	List         []*db_access.GetFilesListRow
	last_fetch   time.Time
}

func (fls *FilesListCache) cached_fetch(ctx context.Context) ([]*db_access.GetFilesListRow, error) {

	storage_struct.RLock()
	defer storage_struct.RUnlock()
	if fls.List == nil || time.Since(fls.last_fetch) > (time.Duration(server_config.Storage.PollIntervalMinutes/2)*time.Minute) {
		return fls.fetch(ctx)
	}
	fls.RLock()
	defer fls.RUnlock()

	return fls.List, nil
}

func (fls *FilesListCache) fetch(ctx context.Context) ([]*db_access.GetFilesListRow, error) {
	loc := log_item.Loc(`func (fls *FilesListCache) fetch(ctx context.Context) ([]*db_access.GetFilesListRow, error)`)
	db_conn := DBPool.GetConn()
	defer DBPool.Return(db_conn)
	fls.Lock()
	defer fls.Unlock()
	files_list, err1 := DB.GetFilesList(ctx, db_conn)
	fls.last_fetch = time.Now()

	fls.List = files_list
	if err1 != nil {
		if db_helpers.CheckNoRowsInResultSet(err1) {

			return fls.List, nil
		}

		Logger.LogErr(loc, err1)

		return fls.List, err1
	}

	for _, file_i := range files_list {
		fls.search_cache.plain = append(fls.search_cache.plain, &FilesListResponsType{
			FilePath: file_i.FilePath,
			FileType: file_i.FileType,
			FileHash: file_i.FileHash,
			ModTime:  file_i.ModTime.Time,
		})
	}

	fls.search_cache.list = map[string][]*FilesListResponsType{}

	return fls.List, nil

}

func (fls *FilesListCache) search(ctx *gin.Context) (response_list []*FilesListResponsType) {
	response_list = []*FilesListResponsType{}
	name_search, name_search_ok := ctx.Request.URL.Query()["name"]
	found := !name_search_ok
	fls.RLock()
	defer fls.RUnlock()
	if found {
		return fls.search_cache.plain
	}

	lower_cased_params := []string{}
	for _, name_search_param := range name_search {
		if len(name_search_param) < 1 {
			continue
		}
		unescaped_param, err := url.QueryUnescape(name_search_param)
		if err != nil {
			lower_cased_params = append(lower_cased_params, strings.ToLower(name_search_param))
		}
		lower_cased_params = append(lower_cased_params, strings.ToLower(unescaped_param))
	}

	cache_key := strings.Join(lower_cased_params, "|")
	search_key, found_list := fls.search_cache.list[cache_key]
	if found_list {
		return search_key
	}

	for _, files_list_item := range fls.search_cache.plain {
		file_p_lowercase := strings.ToLower(files_list_item.FilePath)

		if name_search_ok && name_search != nil {
			found = false
			for _, name_search_param := range lower_cased_params {
				if strings.Contains(file_p_lowercase, name_search_param) {
					found = true
					break
				}
			}
		}

		if found {
			response_list = append(response_list, files_list_item)
		}
	}

	fls.search_cache.list[cache_key] = response_list

	return
}

func recover_func(after *string, loc log_item.Loc) {
	logging.RecoverFunc(after, loc)
}

func get_files_list_route(ctx *gin.Context) {
	loc := log_item.Locf(`api_route.GET("/files_list", func(ctx *gin.Context)`)
	after := "setup"
	defer recover_func(&after, loc)

	response_list := []*FilesListResponsType{}
	files_list, err1 := files_list_cache.cached_fetch(ctx)
	if err1 != nil {
		after = `if db_helpers.CheckNoRowsInResultSet(err1)`
		if db_helpers.CheckNoRowsInResultSet(err1) {
			ctx.JSON(200, gin.H{
				success: true,
				"data":  response_list,
			})
			return
		}

		Logger.LogErr(loc, err1)
		after = `ctx.JSON(500, gin.H`
		ctx.JSON(500, gin.H{
			success: false,
		})
		return
	}

	after = `if files_list == nil`
	if files_list == nil {
		ctx.JSON(200, gin.H{
			success: true,
			"data":  response_list,
		})
		return
	}

	after = `name_search, name_search_ok := ctx.Request.URL.Query()["name"]`
	response_list = files_list_cache.search(ctx)

	ctx.JSON(200, gin.H{
		success: true,
		"data":  response_list,
	})
}
