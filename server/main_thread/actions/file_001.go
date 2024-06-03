package actions

import (
	"fmt"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler"
)

func Write_directory_files_list(dir_path string, files []filehandler.FileBasic) (err *ftp_context.LogItem) {
	loc := "Write_directory_files_list(dir_path string, files []filehandler.FileBasic) (err *ftp_context.LogItem)"
	name := func() string {
		a := time.Now()
		b := fmt.Sprintf("files/%d/%02d_%02d.json", a.Year(), a.Month(), a.Day())
		return b
	}()

	txt_file := filehandler.NewFileBasic(dir_path + "\\" + name)

	err = txt_file.Create().
		WriteJson(files)

	if err != nil {
		return ftp_context.NewLogItem(loc, true).SetAfter("txt_file.Create().WriteJson(files)").AppendParentError(err)
	}

	return
}
