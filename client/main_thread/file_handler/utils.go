package filehandler

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	ftp_context "github.com/ftp_system_client/main_thread/context"
)

func ReadDir(ctx ftp_context.Context, dir_path string, exclude []string) (files []FileBasic, err error) {
	ctx.Add()
	defer ctx.Finished()
	if len(dir_path) < 1 {
		log.Fatal("need to give directory path")
	}

	err = filepath.WalkDir(dir_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, excluded := range exclude {
			if strings.Contains(path, excluded) {
				return nil
			}
		}

		fo := FileBasic{
			Name: d.Name(),
			Path: path,
			Type: "",
			d:    d,
		}

		fo.Ext()
		
		files = append(files, fo)

		return nil
	})

	return
}
