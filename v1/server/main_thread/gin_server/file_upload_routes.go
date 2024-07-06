package ginserver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"path/filepath"

	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/it-shiloheye/ftp_system/server/initialise_server"
	ftp_dirhandler "github.com/it-shiloheye/ftp_system_client/main_thread/dir_handler"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
	"github.com/it-shiloheye/ftp_system_lib/logging/log_item"
)

var client_filetrees = ftp_base.NewMutexedMap[*ftp_dirhandler.FileTreeJson]()

func init() {

	for _, client_id := range ServerConfig.Clients.Keys() {
		client_dir_p, ftree_json := client_directory_path(client_id)

		ftree := ftp_dirhandler.NewFileTreeJson()

		fo, err1 := os.Open(ftree_json)
		if err1 != nil {
			if !errors.Is(err1, os.ErrNotExist) {
				log.Fatalln(err1)
			}
			os.MkdirAll(client_dir_p, fs_mode)

			d, err1 := json.MarshalIndent(ftree, " ", "\t")
			if err1 != nil {
				log.Fatalln(err1)
			}

			err2 := os.WriteFile(ftree_json, d, fs_mode)
			if err1 != nil {
				log.Fatalln(err2)
			}

		} else {
			err2 := json.NewDecoder(fo).Decode(ftree)
			if err2 != nil {
				log.Fatalln(err2)
			}
			fo.Close()
		}

		client_filetrees.Set(client_id, ftree)
	}
}

func PrintToJson(to_json any, name string, kill ...bool) {
	d, err1 := json.MarshalIndent(to_json, " ", "\t")
	if err1 != nil {
		log.Fatalln(err1)
	}
	err2 := os.WriteFile("./data/"+name+".json", d, fs.FileMode(ftp_base.S_IRWXU|ftp_base.S_IRWXO))
	if err2 != nil {
		log.Fatalln(err2)
	}

	if len(kill) > 0 {
		if kill[0] {
			log.Println("exiting PrintToJson: ", name)
			os.Exit(1)
		}
	} else {
		log.Println("exiting PrintToJson: ", name)
		os.Exit(1)
	}
}

func confirm_uploaded_files(client_id string, hashed_files []string) (to_hash map[string]bool, err error) {
	loc := log_item.Locf(`confirm_uploaded_files(client_id: %s, f_tree *ftp_dirhandler.FileTreeJson)(tmp map[string]string,err error)`, client_id)

	to_hash = map[string]bool{}
	client_storage_path, _ := client_directory_path(client_id)
	client_ftree, ok := client_filetrees.Get(client_id)
	if !ok {

		err = Logger.LogErr(loc, &log_item.LogItem{CallStack: []error{fmt.Errorf("missing client from store")}, After: fmt.Sprintf(`client_ftree, %v := client_filetrees.Get(%s)`, ok, client_id)})
		return
	}

	l, err3 := filehandler.Lock(client_storage_path + "-confirm-filetree.lock")
	if err3 != nil {
		err = Logger.LogErr(loc, &log_item.LogItem{
			After:     fmt.Sprintf(`l, err3 := filehandler.Lock("%s")`, client_storage_path+"-confirm-filetree.lock"),
			CallStack: []error{err3},
		})
		return
	}
	defer l.Unlock()
	Logger.Logf(loc, "client_storage_path:\t%s", client_storage_path)

	client_ftree.RLock()
	for _, hash := range hashed_files {

		_, ok := client_ftree.FileMap.M[hash]
		to_hash[hash] = ok

	}

	err1 := filepath.WalkDir(client_storage_path, func(path string, d fs.DirEntry, err1 error) error {
		loc := log_item.Locf(`fs.WalkDir(os.DirFS(client_storage_path: %s), client_storage_path, func(path string, d fs.DirEntry, err1: %v) error `, client_storage_path, err1)
		if err1 != nil {
			err = Logger.LogErr(loc, err1)
			log.Fatalln(err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		file_name_hash := d.Name()
		if strings.Contains(file_name_hash, ".json") {
			return nil
		}

		if to_hash[file_name_hash] {
			delete(to_hash, file_name_hash)
		}

		client_ftree.FileState.M[file_name_hash] = ftp_dirhandler.FileStateUploaded
		return nil

	})

	if err1 != nil {
		err = Logger.LogErr(loc, err1)
		log.Println(loc, "\nerror after walkdir")
		return
	}

	PrintToJson(client_ftree, "test-ftree")
	client_ftree.RUnlock()
	client_filetrees.Set(client_id, client_ftree)
	log.Println(`client_filetrees.Set(client_id:`, client_id, `, client_ftree)\n`, client_ftree)
	err2 := WriteClientFileTree(client_id)
	if err2 != nil {
		err = Logger.LogErr(loc, err2)
		return
	}

	return
}

type ConfirmFileTreeResponse struct {
	State    string   `json:"state"`
	ToUpload []string `json:"to_upload"`
}

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

	server_res := ConfirmFileTreeResponse{
		ToUpload: []string{},
	}
	defer log.Println(&server_res)
	_, ok := client_filetrees.Get(client_id)
	if !ok {

		log.Println("client_filetree missing")
		client_storage_path, _ := client_directory_path(client_id)
		ServerConfig.Clients.Set(client_id, &initialiseserver.ClientIDStruct{
			ClientId:    client_id,
			CommonName:  "",
			StoragePath: client_storage_path,
		})

		server_res.State = "missing"

		ctx.JSON(201, &server_res)

		return
	}

	err4 := json.NewDecoder(ctx.Request.Body).Decode(&server_res.ToUpload)
	if err4 != nil {

		ctx.JSON(500, gin.H{
			"state": failure,
		})
		Logger.LogErr(loc, err4)
		return
	}

	to_hash, err1 := confirm_uploaded_files(client_id, server_res.ToUpload)

	if err1 != nil {
		Logger.LogErr(loc, err1)
		server_res.State = failure
		ctx.JSON(400, &server_res)
		return
	}

	server_res.State = success
	for k, uploaded_ := range to_hash {
		server_res.ToUpload = append(server_res.ToUpload, k)
		if !uploaded_ && server_res.State != success {
			server_res.State = missing
		}
	}
	ctx.JSON(200, &server_res)

}

func get_dir_and_name(data_dir string, file_hash string) (dir string, file_path string) {
	dir1, dir2 := file_hash[:2], file_hash[2:4]
	dir_p := strings.Join([]string{data_dir, dir1, dir2}, string(os.PathSeparator))

	return dir_p, dir_p + string(os.PathSeparator) + file_hash
}

func client_directory_path(client_id string) (dir_p string, filetree_json string) {
	dir_p = ServerConfig.DirConfig.Path + string(os.PathSeparator) + client_id
	filetree_json = dir_p + string(os.PathSeparator) + "file-tree.json"
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
	f_tree, ok := client_filetrees.Get(client_id)
	if !ok {
		f_tree = ftp_dirhandler.NewFileTreeJson()
	}
	defer WriteClientFileTree(client_id)
	defer client_filetrees.Set(client_id, f_tree)

	f_tree.RLock()
	defer f_tree.RUnlock()

	found_file := 0
	dir_p, _ := client_directory_path(client_id)

	decode := base64.StdEncoding.DecodeString

	c := len(req_body)

filemap_loop:
	for _, fh := range f_tree.FileMap.M {

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

			Logger.Logf(loc, "uploaded: %s", fh.Hash)
			f_tree.FileState.Set(fh.Hash, ftp_dirhandler.FileStateUploaded)
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
