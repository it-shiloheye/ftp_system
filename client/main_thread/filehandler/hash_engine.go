package dir_handler

import (
	
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

func hashing_piston(ctx ftp_context.Context, file_path_chan <-chan string, done chan string, err chan error) {
	loc := "hashing_piston(ctx ftp_context.Context,file_path_chan <-chan string, err chan error)"
	defer ctx.Finished()

	bts := filehandler.NewBytesStore()
	var file_path string

	var tmp_fh *filehandler.FileHash
	exists := false
	var err1, err2 error
	for ok := true; ok; {
		select {
		case _, ok = <-ctx.Done():
			break
		case file_path, ok = <-file_path_chan:
		}

		bts.Reset()
		tmp_fh, exists = FileTree.FileMap.Get(file_path)
		if !exists {
			tmp_fh, err1 = filehandler.NewFileHashOpen(file_path)
			if err1 != nil {
				err <- ftp_context.NewLogItem(loc, true).
					SetMessage(err1.Error()).
					SetAfterf(`tmp_fh, err1 := filehandler.NewFileHashOpen(%s)`, file_path).
					AppendParentError(err1)
				continue
			}

			FileTree.FileMap.Set(file_path,tmp_fh)
		}

		tmp_fh.Hash, err2 = filehandler.HashFile(tmp_fh.FileBasic, bts)
		if err2 != nil {
			err <- ftp_context.NewLogItem(loc, true).
				SetMessage(err2.Error()).
				SetAfterf(`tmp_fh.Hash, err2 =	filehandler.HashFile(tmp_fh.FileBasic,bts)`).
				AppendParentError(err2)
			continue
		}

		FileTree.HashQueue.Set(file_path,tmp_fh)
	}

}
