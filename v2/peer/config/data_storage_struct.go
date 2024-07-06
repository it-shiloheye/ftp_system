package server_config

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	ftp_base "github.com/it-shiloheye/ftp_system/v2/lib/base"
	db_access "github.com/it-shiloheye/ftp_system/v2/lib/db_access/generated"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/it-shiloheye/ftp_system/v2/lib/logging"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging/log_item"
)

var Logger = logging.Logger

type StorageStruct struct {
	sync.RWMutex

	PeerId           pgtype.UUID            `json:"peer_id"`
	PeerRole         db_access.PeerRoleType `json:"peer_role"`
	StorageDirectory string                 `json:"storage_directory"`
	IncludeDirs      []string               `json:"include_dirs"`
	ExcludeDirs      []string               `json:"exclude_dirs"`
	IncludeFiles     []string               `json:"include_files"`
	ExcludeRegex     []string               `json:"exclude_regex"`
	IncludeRegex     []string               `json:"include_regex"`
	OnUpload         OnUploadStruct         `json:"on_upload"`
	PollInterval     time.Duration          `json:"poll_interval"`
	SubscribedPeers  []*SubscribePeers      `json:"subscribed_peers"`
}

type OnUploadStruct struct {
	DeleteOnUpload           bool `json:"delete_on_upload"`
	MaxAgeInDaysBeforeDelete int  `json:"max_age_in_days_before_delete"`
}
type SubscribePeers struct {
	PeerUuid               pgtype.UUID `json:"peer_uuid"`
	CopyDirectoryStructure bool        `json:"copy_directory_structure"`
	LocalName              string      `json:"local_name"`
	PublishOnEdits         string      `json:"publish_on_edits"`
	CascadeDeletes         string      `json:"cascade_deletes"`
	DownloadOnChanges      string      `json:"download_on_changes"`
}

func NewStorageStruct() (sts *StorageStruct) {
	sts = &StorageStruct{
		IncludeDirs:      []string{},
		PeerRole:         db_access.PeerRoleTypeClient,
		StorageDirectory: "./data",
		ExcludeDirs:      []string{".git", "tmp", "~"},
		IncludeFiles:     []string{},
		ExcludeRegex:     []string{},
		IncludeRegex:     []string{},
		PollInterval:     time.Minute * 5,
		SubscribedPeers: []*SubscribePeers{{
			CopyDirectoryStructure: true,
			LocalName:              "test",
		}},
		OnUpload: OnUploadStruct{
			DeleteOnUpload:           false,
			MaxAgeInDaysBeforeDelete: -1,
		},
	}

	return
}

func WriteToDisk(file_path string, storage_struct *StorageStruct) error {
	loc := log_item.Locf(`WriteToDisk(file_path: %s, str_struct *StorageStruct) error`, file_path)

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

func ReadFromDisk(file_path string, sts *StorageStruct) (err error) {
	loc := log_item.Locf(`ReadFromDisk(file_path: %s) (sts *StorageStruct, err error)`, file_path)

	f, err1 := os.Open(file_path)
	if err1 != nil {
		err = Logger.LogErr(loc, err1)
		return
	}

	err2 := json.NewDecoder(f).Decode(sts)
	if err2 != nil {
		err = Logger.LogErr(loc, err2)
		return
	}

	return
}
