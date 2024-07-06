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
)

var Logger = logging.Logger
var storage_struct = server_config.Storage

func main() {
	loc := log_item.Loc(`main()`)
	log.Println("this is a test")

	ctx := ftp_context.CreateNewContext()

	logging.InitialiseLogging(".")
	go logging.Logger.Engine(ctx, ".")
	db.ConnectToDB(ctx)

	server_config.LoopReadStorageStruct(3, storage_struct)
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

	go mainthread.PermanentUploadLoop(ctx.Add())

	go mainthread.Loop(ctx.Add())
	ctx.Wait()
}
