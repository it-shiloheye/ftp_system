package dir_handler

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"path/filepath"
	"strings"

	"time"

	initialiseclient "github.com/it-shiloheye/ftp_system/client/init_client"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
	// filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var ClientConfig = initialiseclient.ClientConfig

type ReadDirResult struct {
	FilesList []*filehandler.FileBasic
	ToRehash  []string
	ToUpload  []string
}

func ticker(loc string, i int) {

	// log.Println(loc, i)
}

func ReadDir(ctx ftp_context.Context, dir_data initialiseclient.DirConfig) (rd ReadDirResult, err ftp_context.LogErr) {
	loc := "ReadDir(ctx ftp_context.Context, dir_data initialiseclient.DirConfig) (err ftp_context.LogErr)"
	defer ctx.Finished()
	ticker(loc, 1)
	rd = ReadDirResult{
		ToRehash: []string{},
		ToUpload: []string{},
	}
	excluded_dirs_ := FlatMap[string](
		append(ClientConfig.ExcludeDirs, ".git"),
		ClientConfig.ExcluedFile,
		dir_data.ExcludeDirs,
		dir_data.ExcluedFile,
	)
	ticker(loc, 2)
	dirs_excluded_dirs_list := []string{}
	dir_uniq := map[string]bool{}
	for _, d := range excluded_dirs_ {
		if len(d) < 1 || dir_uniq[d] {
			continue
		}
		a := strings.Join(strings.Split(d, string(os.PathSeparator)), "/")
		b := strings.Join(strings.Split(d, string(os.PathSeparator)), "\\")
		dirs_excluded_dirs_list = append(dirs_excluded_dirs_list, a, b)
		dir_uniq[a] = true
		dir_uniq[b] = true
	}
	ticker(loc, 3)
	// log.Println("dir data path", dir_data.Path, "\n", strings.Join(dirs_list, "\n"))
	var err1 error
	rd.FilesList, err1 = list_file_tree(dir_data.Path, dirs_excluded_dirs_list)
	if err1 != nil {
		return rd, &ftp_context.LogItem{Location: loc, Time: time.Now(),
			After: `a, err1 := filehandler.ReadDir(ctx, dir_data.Path, FlatMap(
				append(ClientConfig.ExcludeDirs, ".git"),
				ClientConfig.ExcluedFile,
				dir_data.ExcludeDirs,
				dir_data.ExcluedFile,
			))`,
			Message:   "didn't successfully read dir",
			CallStack: []error{err1},
		}
	}
	ticker(loc, 4)
	for _, file := range rd.FilesList {
		f_path := file.Path
		log.Println(f_path)
		if file.IsDir() {
			continue
		}
		FileTree.FileMap.Set(f_path, &filehandler.FileHash{
			FileBasic: file,
		})

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
	log.Println("successfully read dir at ", time.Now().Format(time.RFC822))
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

		if NilError(err2) {

			return &ftp_context.LogItem{
				Location:  loc,
				Time:      time.Now(),
				After:     after,
				Message:   err2.Error(),
				CallStack: []error{err2},
			}
		}

		if fs_d.IsDir() {
			return nil
		}

		for _, excluded := range exclude_paths {
			if strings.Contains(path, excluded) {
				log.Println("excluded:", excluded, path)
				return nil
			}

		}

		tmp, err3 := filehandler.Open(path)
		if NilError(err3) {

			return &ftp_context.LogItem{
				Location:  loc,
				Time:      time.Now(),
				After:     fmt.Sprintf(`tmp, err3 := filehandler.Open("%s")`, path),
				Message:   err3.Error(),
				CallStack: []error{err3},
			}
		}
		log.Println("appended:", path)
		out = append(out, tmp)

		return nil
	})

	if NilError(err1) {
		return out, err1
	}

	return
}
