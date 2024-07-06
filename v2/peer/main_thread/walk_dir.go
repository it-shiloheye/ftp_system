package mainthread

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ftp_base "github.com/it-shiloheye/ftp_system/v2/lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system/v2/lib/context"

	"github.com/it-shiloheye/ftp_system/v2/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system/v2/peer/config"
)

func WalkDir(ctx ftp_context.Context, storage_struct *server_config.StorageStruct) (files_list []*FilesListItem, err error) {
	loc := log_item.Loc(`func WalkDir(ctx ftp_context.Context, storage_struct *server_config.StorageStruct)(files_list []FilesListItem, err  error)`)
	files_list = []*FilesListItem{}

	uniq := map[string]bool{}
	for _, dir_path := range storage_struct.IncludeDirs {
		err1 := filepath.WalkDir(dir_path, func(path string, d fs.DirEntry, err2 error) error {
			loc := log_item.Locf(`err1 := filepath.WalkDir(dir_path: %s,func(path %s, d fs.DirEntry, err2: %s) error `, dir_path, path, fmt.Sprint(err2))
			if err2 != nil {
				return Logger.LogErr(loc, err2)
			}

			if _, ok := uniq[path]; ok || d.IsDir() {
				return nil
			}
			uniq[path] = true

			for _, invalid_dir := range storage_struct.ExcludeDirs {
				if strings.Contains(path, invalid_dir) {
					return nil
				}
			}
			for _, invalid_regex := range storage_struct.ExcludeRegex {
				if not_ok, err3 := regexp.MatchString(invalid_regex, path); not_ok || err3 != nil {
					if err3 != nil {
						Logger.LogErr(loc, &log_item.LogItem{
							After:     fmt.Sprintf(`not_ok, err3 := regexp.MatchString(invalid_regex: %s ,path: %s);not_ok || err3 != nil `, invalid_regex, path),
							Level:     log_item.LogLevelError02,
							CallStack: []error{err3},
						})
					}
					return nil
				}
			}

			add := false

			if len(storage_struct.IncludeRegex) > 0 {
				for _, fpath := range storage_struct.IncludeRegex {
					ok, err4 := regexp.MatchString(fpath, path)
					if ok {
						add = true
						break
					}
					if err4 != nil {
						Logger.LogErr(loc, err4)
					}
				}

			} else {
				add = true
			}
			if add {

				tmp := &FilesListItem{
					Path: path,
				}

				var err5, err6 error
				tmp.File, err5 = os.OpenFile(path, os.O_RDONLY, ftp_base.FS_MODE)
				if err5 != nil {
					return Logger.LogErr(loc, err5)
				}

				tmp.FD, err6 = tmp.File.Stat()
				if err6 != nil {
					return Logger.LogErr(loc, err6)
				}
				tmp.Close()
				files_list = append(files_list, tmp)
			}
			return nil
		})

		if err1 != nil {
			err = Logger.LogErr(loc, err1)

			return
		}
	}

	for _, fpath := range storage_struct.IncludeFiles {
		if _, ok := uniq[fpath]; ok {
			continue
		}
		uniq[fpath] = true

		tmp := &FilesListItem{
			Path: fpath,
		}
		err1 := tmp.Reopen()
		if err1 != nil {
			Logger.LogErr(loc, err1)
			continue
		}
		tmp.Close()
		files_list = append(files_list, tmp)

	}

	return
}
