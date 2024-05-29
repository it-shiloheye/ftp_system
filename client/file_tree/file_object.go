package filetree

import (
	"os"
	"strings"
)

type FileObjectStruct struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Type        FileType `json:"type"`
	Hash        string   `json:"hash"`
	Size        int64    `json:"size"`
	ModTime     string   `json:"mod_time"`
	PrevModTime string   `json:"prev_mod_time"`
	file        *os.File
}

func (fo *FileObjectStruct) SearchName(substr string) bool {
	return strings.Contains(fo.Name, substr)
}

func (fo *FileObjectStruct) SearchPath(substr string) bool {
	return strings.Contains(fo.Path, substr)
}
