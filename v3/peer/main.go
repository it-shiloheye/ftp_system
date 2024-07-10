package main

import (
	"log"
	"strings"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	db "github.com/it-shiloheye/ftp_system/v3/lib/db_access"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system/v3/peer/config"
	"github.com/it-shiloheye/ftp_system/v3/peer/mainthread"
	"github.com/it-shiloheye/ftp_system/v3/peer/mainthread/db_helpers"
	// networkpeer "github.com/it-shiloheye/ftp_system/v3/peer/network-peer"
)

var Logger = logging.Logger
var storage_struct = server_config.Storage

func main() {
	loc := log_item.Loc(`main()`)
	log.Println("this is a test")

	_ = loc

	ctx := ftp_context.CreateNewContext()

	logging.InitialiseLogging(".")
	go logging.Logger.Engine(ctx, ".")
	db.ConnectToDB(ctx)

	server_config.LoopReadStorageStruct(3, storage_struct)
	connect_db_client(ctx, storage_struct)
	go mainthread.PermanentUploadLoop(ctx.Add())

	go mainthread.Loop(ctx.Add())

	for ok := true; ok; {
		err_c := make(chan error) //networkpeer.CreateBrowserServer(ctx)

		select {
		case _, ok = <-ctx.Done():
		case err := <-err_c:
			Logger.LogErr(loc, err)

		}
	}
	ctx.Wait()
}

func connect_db_client(ctx ftp_context.Context, storage_struct *server_config.StorageStruct) {
	loc := log_item.Loc(`func connect_client(ctx ftp_context.Context, storage_struct *server_config.StorageStruct)`)
	Logger.Logf(loc, "attempting to connect to db")
	for {
		err1 := db_helpers.ConnectClient(ctx, storage_struct)
		if err1 != nil {
			if strings.Contains(err1.Error(), "A socket operation was attempted to an unreachable network") {
				Logger.LogErr(loc, err1)
				<-time.After(time.Minute)
				continue
			}

			log.Fatalln(Logger.LogErr(loc, err1))
		} else {
			break
		}
	}
}
