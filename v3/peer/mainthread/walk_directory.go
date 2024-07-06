package mainthread

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

// map[string]*FileItem
type FileMap = map[string]*FileItem
type FileState string

const (
	// under "WalkDir"
	fstate_new         FileState = "new"
	fstate_to_upload   FileState = "to_upload"
	fstate_to_download FileState = "to_download"
	fstate_unchanged   FileState = "unchanged"

	fstate_uploaded     FileState = "uploaded"
	fstate_upload_err   FileState = "upload_error"
	fstate_downloaded   FileState = "downloaded"
	fstate_download_err FileState = "download_error"

	fstate_deleted FileState = "deleted"
	fstate_os_err  FileState = "os_error"

	fstate_changed FileState = "changed" // -> useless, mark as upload or download

)

func (f FileState) String() string {
	return string(f)
}

func (f FileState) IsErr() bool {
	switch f {
	case fstate_download_err:
		return true
	case fstate_upload_err:
		return true
	case fstate_os_err:
		return true
	}

	return false
}

type FileItem struct {
	FileState `json:"file_state"`
	FileHash  string `json:"file_hash"`
	ListIndex int    `json:"db_index"`
	Path      string `json:"file_path"`
	stats     fs.FileInfo
}

func WalkDir(storage_path string, files_map *FileMapType, watcher *fsnotify.Watcher) (err error) {
	loc := log_item.Locf(`func WalkDir(storage_path: "%s", file_map map[string]*FileItem) (err error) `, storage_path)
	file_map := (files_map.FileMap)
	storage_struct.RLock()
	defer storage_struct.RUnlock()

	// files already on the local storage
	err1 := filepath.WalkDir(storage_path, func(file_path string, d fs.DirEntry, err2 error) error {
		loc := log_item.Locf(`err1 := filepath.WalkDir(storage_path:%s, func(dir_path: %s, d fs.DirEntry, err2: %v) error`, storage_path, file_path, err2)

		if err2 != nil {
			Logger.LogErr(loc, err2)
			return nil
		}

		fpath := get_relative_path(storage_path, file_path)

		if is_excluded_dir(file_path) {

			return nil
		}

		if ok, err4 := is_excluded_regex(loc, file_path); ok || err4 != nil {
			if err4 != nil {

				return err4
			}
			return nil
		}
		if d.IsDir() {
			watcher.Add(fpath)
			return nil
		}

		stats, err3 := d.Info()
		if err3 != nil {
			return Logger.LogErr(loc, err3)
		}

		log.Println("found: ", fpath)
		f_stored, ok := file_map.Get(fpath)

		if !ok {
			file_map.Set(fpath, &FileItem{
				FileState: fstate_new,
				Path:      fpath,
				stats:     stats,
			})
			return nil
		}

		if f_stored.SameModtime(stats.ModTime()) && !f_stored.IsErr() {
			f_stored.FileState = fstate_unchanged
			return nil
		}

		if f_stored.IsAfter(stats.ModTime()) || f_stored._upload_err() {
			f_stored.FileState = fstate_to_upload
			f_stored.stats = stats
			return nil
		}

		return nil
	})

	if err1 != nil {
		err = Logger.LogErr(loc, err1)

		return
	}

	return
}

// returns "/" comma delimited path with "dir_path" as root directory
func get_relative_path(dir_path string, file_path string) string {
	if dir_path == file_path {
		return ""
	}

	l_d := len(dir_path)

	tmp01 := file_path[l_d+1:]
	tmp02 := strings.Split(tmp01, string(os.PathSeparator))
	return strings.Join(tmp02, "/")
}

// returns os specific full filepath
func (fi *FileItem) Full(dir string) string {
	tmp01 := dir + "/" + fi.Path
	tmp02 := strings.Split(tmp01, "/")

	return strings.Join(tmp02, string(os.PathSeparator))
}

func (fi *FileItem) _new() bool         { return fi.FileState == "new" }
func (fi *FileItem) _to_upload() bool   { return fi.FileState == "to_upload" }
func (fi *FileItem) _to_download() bool { return fi.FileState == "to_download" }

func (fi *FileItem) _uploaded() bool     { return fi.FileState == "uploaded" }
func (fi *FileItem) _upload_err() bool   { return fi.FileState == "upload_error" }
func (fi *FileItem) _downloaded() bool   { return fi.FileState == "downloaded" }
func (fi *FileItem) _download_err() bool { return fi.FileState == "download_error" }

func (fi *FileItem) _unchanged() bool { return fi.FileState == "unchanged" }
func (fi *FileItem) _deleted() bool   { return fi.FileState == "deleted" }
func (fi *FileItem) _os_err() bool    { return fi.FileState == "os_error" }

func (fi *FileItem) SameModtime(t time.Time) bool {

	return fi.stats != nil && fi.stats.ModTime() == t
}
func (fi *FileItem) IsBefore(t time.Time) bool {
	return fi.stats.ModTime().Before(t)
}

func (fi *FileItem) IsAfter(t time.Time) bool {
	return fi.stats != nil && fi.stats.ModTime().After(t)
}

func is_excluded_regex(loc log_item.Loc, file_path string) (bool, error) {
	for _, excluded := range storage_struct.ExcludeRegex {
		if ok, err3 := regexp.MatchString(excluded, file_path); ok || err3 != nil {
			if err3 != nil {

				return true, Logger.LogErr(loc, err3)
			}
			for _, included := range storage_struct.IncludeRegex {
				if strings.Contains(file_path, included) {
					return false, nil
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func is_excluded_dir(file_path string) bool {
	if file_path == "" {
		return true
	}
	for _, excluded := range storage_struct.ExcludeDirs {
		if strings.Contains(file_path, excluded) {

			for _, included := range storage_struct.IncludeRegex {
				if strings.Contains(file_path, included) {
					return false
				}
			}
			return true
		}
	}
	return false
}

type FiFileInfo struct {
	name    string
	size    int64
	mode    int
	modtime time.Time
	is_dir  bool
}

func (fi *FiFileInfo) Name() string { // base name of the file{
	return fi.name
}
func (fi *FiFileInfo) Size() int64 { // length in bytes for regular files; system-dependent for others{
	return fi.size
}
func (fi *FiFileInfo) Mode() fs.FileMode { // file mode bits{
	return fs.FileMode(fi.mode)
}
func (fi *FiFileInfo) ModTime() time.Time { // modification time{
	return fi.modtime
}
func (fi *FiFileInfo) IsDir() bool { // abbreviation for Mode().IsDir(){
	return fi.is_dir
}
func (fi *FiFileInfo) Sys() any { // underlying data source (can return nil){
	return nil
}
