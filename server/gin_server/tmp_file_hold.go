package ginserver

import (
	"github.com/it-shiloheye/ftp_system_client/main_thread/dir_handler"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var tmp_file_hold = ftp_base.NewMutexedMap[*filehandler.FileHash]()
var contents = ftp_base.NewMutexedMap[[]byte]()

var FileTree = dir_handler.FileTree

func init() {
	file_tree_path := ServerConfig.DirConfig.Path + "/file-tree.json"

	dir_handler.InitialiseFileTree(file_tree_path)
}
