package server_config

import (
	"encoding/json"
	"log"
	"os"

	"sync"
	"time"

	ftp_base "github.com/it-shiloheye/ftp_system/v3/lib/base"

	db_access "github.com/it-shiloheye/ftp_system/v3/lib/db_access/generated"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

var Logger = logging.Logger

var Storage = NewStorageStruct()

type StorageStruct struct {
	sync.RWMutex

	PeerId           pgtype.UUID            `json:"peer_id"`
	PeerRole         db_access.PeerRoleType `json:"peer_role"`
	PeerName         string                 `json:"peer_name"`
	StorageDirectory string                 `json:"storage_directory"`

	ExcludeDirs         []string       `json:"exclude_dirs"`
	ExcludeRegex        []string       `json:"exclude_regex"` // exclude directory within storage
	IncludeRegex        []string       `json:"include_regex"` // in excluded but still should be included
	UploadDirs          OnUploadStruct `json:"on_upload"`
	PollIntervalMinutes int            `json:"poll_interval_minutes"`
}

type OnUploadStruct struct {
	DeleteOnUpload           bool `json:"delete_on_upload"`
	MaxAgeInDaysBeforeDelete int  `json:"max_age_in_days_before_delete"`
	// upload only, no download
	UploadDirs []string `json:"upload_directories_fullpath"`
}

const config_filepath = "./config.json"

func NewStorageStruct() (sts *StorageStruct) {
	sts = &StorageStruct{

		PeerRole:    db_access.PeerRoleTypeClient,
		ExcludeDirs: []string{".git", "tmp", "~", "node_modules", "$"},

		ExcludeRegex:        []string{},
		IncludeRegex:        []string{},
		PollIntervalMinutes: 5,
		UploadDirs: OnUploadStruct{
			DeleteOnUpload:           false,
			MaxAgeInDaysBeforeDelete: -1,
			UploadDirs:               []string{},
		},
	}

	return
}

func WriteToDisk(storage_struct *StorageStruct) error {
	file_path := config_filepath
	loc := log_item.Locf(`WriteToDisk(file_path: %s, str_struct *StorageStruct) error`, file_path)
	storage_struct.Lock()
	defer storage_struct.Unlock()
	d, err1 := json.MarshalIndent(storage_struct, " ", "\t")
	if err1 != nil {
		return Logger.LogErr(loc, err1)
	}

	err2 := os.WriteFile(file_path, d, ftp_base.FS_MODE)
	if err2 != nil {
		return Logger.LogErr(loc, err2)
	}
	return nil
}

func ReadFromDisk(storage_struct *StorageStruct) (err error) {
	file_path := config_filepath
	loc := log_item.Locf(`ReadFromDisk(file_path: %s) (sts *StorageStruct, err error)`, file_path)
	storage_struct.Lock()
	defer storage_struct.Unlock()
	f, err1 := os.Open(file_path)
	if err1 != nil {
		err = Logger.LogErr(loc, err1)
		return
	}
	defer f.Close()

	err2 := json.NewDecoder(f).Decode(storage_struct)
	if err2 != nil {
		err = Logger.LogErr(loc, err2)
		return
	}

	return
}

func LoopReadStorageStruct(tries int, sts *StorageStruct) (err error) {
	file_path := config_filepath
	loc := log_item.Locf(`func LoopRead(file_path: "%s", tries: %03d) (sts *StorageStruct, err error)`, file_path, tries)

mainloop:
	for {
		if tries > 0 {
			tries -= 1
		} else {
			break
		}

		err1 := ReadFromDisk(sts)
		if err1 != nil {
			err = logging.Logger.LogErr(loc, err1)
			log.Println("error reading config.json, rewriting config.json")
			err2 := WriteToDisk(sts)
			if err2 != nil {
				log.Fatalln(Logger.LogErr(loc, err2))
			}
			<-time.After(time.Minute)
			continue
		}

		has_storage_dir := len(sts.StorageDirectory) > 5
		has_upload_dirs := len(sts.UploadDirs.UploadDirs) > 0
		if has_storage_dir {

			_, err2 := os.ReadDir(sts.StorageDirectory)
			if err2 != nil {
				log.Println("error reading storage directory: ", sts.StorageDirectory)
				err = Logger.LogErr(loc, err2)
				<-time.After(time.Minute)
				log.Println("plead correct issue in config:\n", err2)
				continue
			}
		}

		if has_upload_dirs {
			for _, dir_path := range sts.UploadDirs.UploadDirs {
				if len(dir_path) < 5 {
					continue
				}
				_, err3 := os.ReadDir(dir_path)
				if err3 != nil {
					log.Println("error reading directory: ", dir_path)
					err = Logger.LogErr(loc, err3)
					<-time.After(time.Minute)
					log.Println("plead correct issue in config \"UploadDirs\":\n", err3)
					continue mainloop
				}
			}
		}

		if !has_storage_dir && !has_upload_dirs {
			log.Println("please add storage directory in config, or upload director(y|ies)")
			<-time.After(time.Minute)
			continue
		}

		return
	}

	return
}
