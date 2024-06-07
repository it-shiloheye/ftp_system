package dir_handler

import (
	"fmt"
	"io/fs"

	"os"

	"path/filepath"
	"strings"

	"time"

	initialiseclient "github.com/it-shiloheye/ftp_system/client/init_client"
	"github.com/it-shiloheye/ftp_system/client/main_thread/logging"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
	// filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var ClientConfig = initialiseclient.ClientConfig
var Logger = logging.Logger

type ReadDirResult struct {
	FilesList []*filehandler.FileBasic
	ToRehash  []string
	ToUpload  []string
}

func ticker(loc string, i int) {

	// Logger.Logf(loc, "%d", i)
}

func ReadDir(ctx ftp_context.Context, dir_data initialiseclient.DirConfig) (rd ReadDirResult, err ftp_context.LogErr) {
	loc := "ReadDir(ctx ftp_context.Context, dir_data initialiseclient.DirConfig) (err ftp_context.LogErr)"
	defer ctx.Finished()
	ticker(loc, 1)
	rd = ReadDirResult{
		ToRehash: []string{},
		ToUpload: []string{},
	}

	ticker(loc, 2)
	dirs_excluded_dirs_list := func() (tmp []string) {
		excluded_dirs_ := FlatMap[string](
			append(ClientConfig.ExcludeDirs, ".git"),
			ClientConfig.ExcluedFile,
			dir_data.ExcludeDirs,
			dir_data.ExcluedFile,
		)
		dir_uniq := map[string]bool{}
		for _, d := range excluded_dirs_ {
			if len(d) < 1 || dir_uniq[d] {
				continue
			}
			if !strings.Contains(d, "\\") && !strings.Contains(d, "/") {
				tmp = append(tmp, d)
				dir_uniq[d] = true
				continue
			}
			a := strings.Join(strings.Split(d, string(os.PathSeparator)), "/")
			b := strings.Join(strings.Split(d, string(os.PathSeparator)), "\\")
			tmp = append(tmp, a, b)
			dir_uniq[a] = true
			dir_uniq[b] = true
		}

		return
	}()

	ticker(loc, 3)

	var err1 error
	rd.FilesList, err1 = list_file_tree(dir_data.Path, dirs_excluded_dirs_list)
	if err1 != nil {

		tmp_err := &ftp_context.LogItem{Location: loc, Time: time.Now(),
			After: `a, err1 := filehandler.ReadDir(ctx, dir_data.Path, FlatMap(
				append(ClientConfig.ExcludeDirs, ".git"),
				ClientConfig.ExcluedFile,
				dir_data.ExcludeDirs,
				dir_data.ExcluedFile,
			))`,

			Message: "didn't successfully read dir",
			Err:     true, CallStack: []error{err1},
		}
		Logger.LogErr(loc, tmp_err)
		return rd, tmp_err
	}
	ticker(loc, 4)
	for _, file := range rd.FilesList {
		f_path := file.Path
		Logger.Logf(loc, f_path)
		if file.IsDir() {
			continue
		}
		_, exists_filemap := FileTree.FileMap.Get(f_path)
		if !exists_filemap {
			FileTree.FileMap.Set(f_path, &filehandler.FileHash{
				FileBasic: file,
				ModTime:   file.Fs().ModTime(),
			})
		}

		hashed, ok := FileTree.HashQueue.Get(f_path)
		if !ok {

			rd.ToRehash = append(rd.ToRehash, f_path)
			continue
		}
		time_passed := file.Fs().ModTime().Equal(hashed.ModTime)
		if time_passed {
			rd.ToRehash = append(rd.ToRehash, f_path)
			continue
		}

		uploaded, ok := FileTree.Uploaded.Get(f_path)
		if !ok {
			rd.ToUpload = append(rd.ToUpload, f_path)
			continue
		}

		current_version := uploaded.Hash == hashed.Hash
		if !current_version {
			rd.ToUpload = append(rd.ToUpload, f_path)
			continue
		}
	}

	ticker(loc, 6)
	Logger.Logf(loc, "successfully read dir at %s", time.Now().Format(time.RFC822))
	return
}

func FlatMap[T any](lists ...[]T) (res []T) {
	l := 0
	for _, listlet := range lists {
		l += len(listlet)
	}
	res = make([]T, l)
	for _, listlet := range lists {
		res = append(res, listlet...)
	}

	return
}

func NilError(err error) bool {
	if err != nil {
		if len(err.Error()) > 0 {
			return true
		}
	}

	return false
}

func list_file_tree(dir_path string, exclude_paths []string) (out []*filehandler.FileBasic, err error) {
	loc := "list_file_tree(dir_path string, exclude_paths []string) (out []*filehandler.FileBasic, err error) "
	err1 := filepath.WalkDir(dir_path, func(path string, fs_d fs.DirEntry, err2 error) error {
		after := fmt.Sprintf(`filepath.WalkDir("%s", func("%s", _ fs.DirEntry, err2 error) error `, dir_path, path)
		// ticker(loc, 1)
		if err2 != nil {

			return &ftp_context.LogItem{
				Location: loc,
				Time:     time.Now(),
				After:    after,

				Message: err2.Error(),
				Err:     true, CallStack: []error{err2},
			}
		}
		// ticker(loc, 2)

		if fs_d.IsDir() {
			return nil
		}

		for _, excluded := range exclude_paths {
			if strings.Contains(path, excluded) {
				// Logger.Logf(loc, "excluded: %s %s", excluded, path)

				return nil
			}

		}
		// ticker(loc, 3)
		tmp, err3 := filehandler.Open(path)
		if NilError(err3) {

			tmp_Err := &ftp_context.LogItem{
				Location: loc,
				Time:     time.Now(),
				After:    fmt.Sprintf(`tmp, err3 := filehandler.Open("%s")`, path),
				Message:  err3.Error(),

				Err: true, CallStack: []error{err3},
			}
			Logger.LogErr(loc, tmp_Err)
			return tmp_Err
		}

		// Logger.Logf(loc, "appended: %s", path)
		out = append(out, tmp)
		ticker(loc, 4)
		return nil
	})

	ticker(loc, 15)
	if NilError(err1) {
		return out, err1
	}

	return
}
