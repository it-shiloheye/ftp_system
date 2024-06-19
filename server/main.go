package main

import (
	"strings"
	"time"

	initialiseserver "github.com/it-shiloheye/ftp_system/server/initialise_server"
	ginserver "github.com/it-shiloheye/ftp_system/server/main_thread/gin_server"
	"github.com/it-shiloheye/ftp_system_client/main_thread/dir_handler"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	"github.com/it-shiloheye/ftp_system_lib/logging"
	"github.com/it-shiloheye/ftp_system_lib/logging/log_item"
)

var ServerConfig = initialiseserver.ServerConfig
var Logger = logging.Logger

func main() {
	loc := log_item.Loc("main")
	data_dir := ServerConfig.DirConfig.Path
	if len(strings.Trim(data_dir, " ")) < 1 {
		ServerConfig.DirConfig.Path = "./data"
		data_dir = ServerConfig.DirConfig.Path
	}

	Logger.Logf(loc, "server started")
	ctx := ftp_context.CreateNewContext()
	logging.InitialiseLogging(data_dir)

	go dir_handler.UpdateFileTree(ctx, data_dir+"/file-tree.lock", data_dir+"/file-tree.json")
	go WriteConfig(ctx.Add())
	go Logger.Engine(ctx.Add(), data_dir)
	go ginserver.StoreUploadedFiles(ctx.Add(), data_dir)
	go ginserver.UpdateClientFileTree(ctx.Add(), data_dir)
	defer ctx.Wait()
	ginserver.NewServer(ctx)
}

func WriteConfig(ctx ftp_context.Context) {
	loc := log_item.Loc("WriteConfig(ctx ftp_context.Context)")
	defer ctx.Finished()
	tc := time.NewTicker(time.Minute * 5)

	for ok := true; ok; {
		select {
		case _, ok = <-ctx.Done():
		case <-tc.C:
		}

		if err1 := initialiseserver.WriteConfigToFile(); err1 != nil {

			err := log_item.NewLogItem(loc, log_item.LogLevelError01).
				SetAfter("err1 := write_config();").SetMessage(err1.Error())
			Logger.LogErr(loc, err)
		}
	}
}
