package ginserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/it-shiloheye/ftp_system/server/initialise_server"
	ftp_dirhandler "github.com/it-shiloheye/ftp_system_client/main_thread/dir_handler"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
	"github.com/it-shiloheye/ftp_system_lib/logging/log_item"
)

var client_filetrees = ftp_base.NewMutexedMap[*ftp_dirhandler.FileTreeJson]()

// receives and saves entire file tree from client
func ConfirmFileTree(ctx *gin.Context) {
	loc := log_item.Loc("func ConfirmFileTree(ctx *gin.Context) ")
	client_id, err1 := ctx.Cookie("client-id")

	if err1 != nil {
		ctx.JSON(403, gin.H{
			"state": failure,
		})
		if err1 != nil {
			Logger.LogErr(loc, err1)
		}
		return
	}

	client_storage_path, _ := client_directory_path(client_id)

	ServerConfig.Clients.Set(client_id, &initialiseserver.ClientIDStruct{
		ClientId:    client_id,
		CommonName:  "",
		StoragePath: client_storage_path,
	})

	l, err3 := filehandler.Lock(client_storage_path + "-confirm-filetree.lock")
	if err3 != nil {
		ctx.JSON(500, gin.H{
			"state": failure,
		})
		Logger.LogErr(loc, err3)
		return
	}

	defer l.Unlock()

	f_tree, ok := client_filetrees.Get(client_id)
	if !ok {
		f_tree = ftp_dirhandler.NewFileTreeJson()
	}
	defer client_filetrees.Set(client_id, f_tree)

	err4 := ctx.BindJSON(f_tree)
	if err4 != nil {
		ctx.JSON(500, gin.H{
			"state": failure,
		})
		Logger.LogErr(loc, err4)
		return
	}

	log.Printf("filestate tree\n%v\n\n", f_tree)
	server_response := map[string]string{}

	hash_list := f_tree.FileState.Keys()
	for _, file_ := range hash_list {
		_, file_p := get_dir_and_name(client_storage_path, file_)
		_, err1 := os.ReadFile(file_p)
		if err1 != nil {
			server_response[file_] = missing
		}
	}

	log.Printf("response\n%v\n\n%v\n\n", &server_response, hash_list)
	server_response["state"] = success
	ctx.JSON(200, &server_response)

	err6 := WriteClientFileTree(client_id)
	if err6 != nil {
		Logger.LogErr(loc, err6)
		return
	}

}

func get_dir_and_name(data_dir string, file_hash string) (dir string, file_path string) {
	dir1, dir2 := file_hash[:2], file_hash[2:4]
	dir_p := strings.Join([]string{data_dir, dir1, dir2}, string(os.PathSeparator))

	return dir_p, dir_p + string(os.PathSeparator) + file_hash
}

type FileTreeUnit struct {
	client_id string
	dir_id    string
}

func client_directory_path(client_id string) (dir_p string, filetree_json string) {
	dir_p = ServerConfig.DirConfig.Path + string(os.PathSeparator) + client_id
	filetree_json = dir_p + string(os.PathSeparator) + fmt.Sprintf("file-tree.json")
	return
}

func WriteClientFileTree(client_id string) error {
	loc := log_item.Locf(`WriteClientFileTree(client_id: %s) error `, client_id)
	ftree, ok := client_filetrees.Get(client_id)
	if !ok {
		return Logger.LogErr(loc, fmt.Errorf("missing filetree for:\nclient-id:\t%s", client_id))
	}

	client_filetrees.Lock()
	defer client_filetrees.Unlock()
	dir_p, file_p := client_directory_path(client_id)
	os.MkdirAll(dir_p, fs_mode)

	file_, err1 := os.OpenFile(file_p, os.O_CREATE|os.O_TRUNC, fs_mode)
	if err1 != nil {
		return err1
	}

	JS := json.NewEncoder(file_)
	JS.SetIndent("", "  ")
	err2 := JS.Encode(ftree)
	if err2 != nil {
		return err2
	}
	return nil
}

var ftree_channel = make(chan FileTreeUnit, 100)

func UpdateClientFileTree(ctx ftp_context.Context, data_dir string) {
	loc := log_item.Loc("UpdateFileTree(ctx ftp_context.Context)")

	defer ctx.Finished()
	tc := time.NewTicker(time.Minute)
	var l *filehandler.LockFile

	tmp := map[string]FileTreeUnit{}
	for ok := true; ok; {
		select {
		case <-tc.C:
		case _, ok = <-ctx.Done():

		}
		l, _ = filehandler.Lock(data_dir + string(os.PathSeparator) + "client-filetree.lock")

		for _, ftree_unit := range client_filetrees.Keys() {
			err := WriteClientFileTree(ftree_unit)
			if err != nil {
				Logger.LogErr(loc, err)
			}
			Logger.Logf(loc, "updated filetree: %s successfully", ftree_unit)
		}
		l.Unlock()
		clear(tmp)
	}
}

func UploadBulkFiles(ctx *gin.Context) {
	loc := log_item.Loc(`UploadBulkFiles(ctx gin.Context)`)

	client_id, _ := ctx.Cookie("client-id")

	// dir_id, _ := ctx.Cookie("dir-id")

	req_body := map[string]string{}

	err0 := ctx.BindJSON(&req_body)
	if err0 != nil {
		ctx.JSON(400, gin.H{
			"state": failure,
		})
		Logger.LogErr(loc, &log_item.LogItem{
			Message:   "not able to bind_json",
			CallStack: []error{err0},
		})
		return
	}

	_, filetree_p := client_directory_path(client_id)

	tmp_filetree := ftp_dirhandler.NewFileTreeJson()
	filetree_file, err1 := filehandler.Open(filetree_p)
	if err1 != nil {
		ctx.JSON(500, gin.H{
			"state": failure,
		})
		Logger.LogErr(loc, &log_item.LogItem{
			Message:   fmt.Sprintf("not able to open filetree_p: %s", filetree_p),
			CallStack: []error{err1},
		})
		return
	}
	defer filetree_file.Close()

	err2 := json.NewDecoder(filetree_file).Decode(tmp_filetree)
	if err2 != nil {
		ctx.JSON(500, gin.H{
			"state": failure,
		})
		Logger.LogErr(loc, &log_item.LogItem{
			Message:   fmt.Sprintf("not able to load json from filetree_p: %s", filetree_p),
			CallStack: []error{err2},
		})
		return
	}

	found_file := 0
	dir_p, _ := client_directory_path(client_id)

	decode := base64.StdEncoding.DecodeString

	c := len(req_body)

filemap_loop:
	for _, fh := range tmp_filetree.FileMap.M {

		for k, v := range req_body {
			if fh.Hash != k {
				continue
			}

			found_file += 1

			file_dir, file_p := get_dir_and_name(dir_p, k)
			os.MkdirAll(file_dir, fs_mode)
			data, err3 := decode(v)
			if err3 != nil {
				ctx.JSON(400, gin.H{
					"state": failure,
				})
				Logger.LogErr(loc, &log_item.LogItem{
					Message:   "not able to decode data stream:",
					CallStack: []error{err3},
				})
				return
			}
			err4 := os.WriteFile(file_p, data, fs_mode)
			if err4 != nil {
				ctx.JSON(400, gin.H{
					"state": failure,
				})
				Logger.LogErr(loc, &log_item.LogItem{
					Message:   "not able to decode data stream:",
					CallStack: []error{err3},
				})
				return
			}
			if found_file == c {
				break filemap_loop
			}
			break
		}

	}

	ctx.JSON(201, gin.H{
		"state": success,
	})

}
