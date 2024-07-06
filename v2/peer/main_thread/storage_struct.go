package mainthread

import (
	"encoding/json"
	"errors"

	"log"
	"os"
	"time"

	ftp_base "github.com/it-shiloheye/ftp_system/v2/lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system/v2/lib/context"

	"github.com/it-shiloheye/ftp_system/v2/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system/v2/peer/config"
	db_helpers "github.com/it-shiloheye/ftp_system/v2/peer/main_thread/db_access"
)

func NewStorageStruct(ctx ftp_context.Context) (StorageStruct *server_config.StorageStruct) {
	StorageStruct = server_config.NewStorageStruct()

	server_config.WriteToDisk("./config.example.json", StorageStruct)

	set_up_storagestruct(ctx, StorageStruct)

	return
}

func set_up_storagestruct(ctx ftp_context.Context, StorageStruct *server_config.StorageStruct) {
	config_filepath := StorageStruct.StorageDirectory + "/config.json"
	loc := log_item.Locf(`set_up_storagestruct := func(StorageStruct :%s) :`, config_filepath)
	f, err1 := os.OpenFile(config_filepath, os.O_RDWR|os.O_EXCL, ftp_base.FS_MODE)
	prev_role := StorageStruct.PeerRole
	if err1 == nil {
		defer f.Close()
		err2 := json.NewDecoder(f).Decode(StorageStruct)
		if err2 != nil {
			log.Fatalln(Logger.LogErr(loc, err2))
		}

		if len(StorageStruct.IncludeDirs) < 1 {
			log.Printf("please add at least one directory to included_dirs in: %s\n", config_filepath)
			<-time.After(time.Minute)
			set_up_storagestruct(ctx, StorageStruct)
			return
		}

		f.Close()

		err3 := server_config.WriteToDisk(config_filepath, StorageStruct)
		if err3 != nil {
			log.Fatalln(Logger.LogErr(loc, err3))
		}
		err4 := db_helpers.UpdatePeerRole(ctx, StorageStruct, &prev_role)
		if err4 != nil {
			log.Fatalln(Logger.LogErr(loc, err4))
		}

		return
	}

	if !errors.Is(err1, os.ErrNotExist) {
		log.Fatalln(Logger.LogErr(loc, err1))
	}

	err3 := server_config.WriteToDisk(config_filepath, StorageStruct)
	if err3 != nil {
		log.Fatalln(Logger.LogErr(loc, err3))
	}

	log.Printf("please fill: %s appropriately\n", config_filepath)
	<-time.After(time.Minute)
	set_up_storagestruct(ctx, StorageStruct)
}
