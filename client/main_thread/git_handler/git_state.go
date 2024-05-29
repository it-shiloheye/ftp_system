package githandler

import (
	"time"

	filehandler "github.com/ftp_system_client/main_thread/file_handler"
)

type GitStateStruct struct {
	Directory        string    `json:"directory"`
	LastCommit       string    `json:"last_commit"`
	MostRecentChange time.Time `json:"most_recent_change"`

	CommitMessage filehandler.FileBasic   `json:"commit_message"`
	PendingFiles  []*filehandler.FileHash `json:"PendingFile"`
	UploadedFiles map[string]*filehandler.FileHash
}
