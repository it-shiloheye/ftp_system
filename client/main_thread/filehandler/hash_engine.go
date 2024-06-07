package dir_handler

import (
	"sync"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

type HashPiston struct {
	sync.WaitGroup
	f_path chan string
	done   chan string
	err    chan error
}

func NewHashingPiston() *HashPiston {
	return &HashPiston{
		f_path: make(chan string, 10),
		done:   make(chan string, 10),
		err:    make(chan error, 10),
	}
}

func (hp *HashPiston) Send(f string) {
	hp.Add(1)
	hp.f_path <- f
}

func (hp *HashPiston) Get() <-chan string {
	return hp.done
}

func (hp *HashPiston) HandleError(ctx ftp_context.Context) {
	defer ctx.Finished()
	for ok := true; ok; {
		select {
		case _, ok = <-ctx.Done():
			break
		case err := <-hp.err:
			Logger.LogErr("HandleError", err)

		}
	}
}

func (hp *HashPiston) Piston(ctx ftp_context.Context) {
	hp.Add(1)
	hashing_piston(ctx, hp.f_path, hp.done, hp.err)
	hp.Done()
}

func (hp *HashPiston) Close() {
	close(hp.f_path)
}

func (hp *HashPiston) Wait() {
	hp.WaitGroup.Wait()
	close(hp.done)
	close(hp.err)
}

func (hp *HashPiston) Reset() {
	hp.f_path = make(chan string, 10)
	hp.done = make(chan string, 10)
	hp.err = make(chan error, 10)
}

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

			FileTree.FileMap.Set(file_path, tmp_fh)
		}

		tmp_fh.Hash, err2 = filehandler.HashFile(tmp_fh.FileBasic, bts)
		if err2 != nil {
			err <- ftp_context.NewLogItem(loc, true).
				SetMessage(err2.Error()).
				SetAfterf(`tmp_fh.Hash, err2 =	filehandler.HashFile(tmp_fh.FileBasic,bts)`).
				AppendParentError(err2)

			continue
		}

		FileTree.HashQueue.Set(file_path, tmp_fh)

	}

}
