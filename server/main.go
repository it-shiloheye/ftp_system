package main

import (
	"time"

	ginserver "github.com/it-shiloheye/ftp_system/server/main_thread/gin_server"
	initialiseserver "github.com/it-shiloheye/ftp_system/server/initialise_server"
	"github.com/it-shiloheye/ftp_system_client/main_thread/dir_handler"

	"github.com/it-shiloheye/ftp_system/server/main_thread/logging"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

var ServerConfig = initialiseserver.ServerConfig
var Logger = logging.Logger

func main() {
	loc := logging.Loc("main")
	Logger.Logf(loc, "server started")
	ctx := ftp_context.CreateNewContext()

	go dir_handler.UpdateFileTree(ctx, ServerConfig.DirConfig.Path+"/file-tree.lock", ServerConfig.DirConfig.Path+"/file-tree.json")
	go WriteConfig(ctx.Add())
	go Logger.Engine(ctx.Add())
	go ginserver.StoreUploadedFiles(ctx.Add(), ServerConfig.Path)
	defer ctx.Wait()
	ginserver.NewServer(ctx)
}

func WriteConfig(ctx ftp_context.Context) {
	loc := logging.Loc("WriteConfig(ctx ftp_context.Context)")
	defer ctx.Finished()
	tc := time.NewTicker(time.Minute * 5)

	for ok := true; ok; {
		select {
		case _, ok = <-ctx.Done():
			break
		case <-tc.C:
		}

		if err1 := initialiseserver.WriteConfigToFile(); err1 != nil {

			err := ftp_context.NewLogItem(string(loc), true).
				SetAfter("err1 := write_config();").SetMessage(err1.Error())
			Logger.LogErr(loc, err)
		}
	}
}
