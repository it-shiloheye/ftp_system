package githandler

import (
	ftp_context "github.com/ftp_system_client/main_thread/context"
	filehandler "github.com/ftp_system_client/main_thread/file_handler"
)

type GitStateStruct struct {
	F        filehandler.FileHash       `json:"git_state"`
	FileTree filehandler.FileTreeStruct `json:"file_tree"`

	buffer_store chan *filehandler.BytesStore

	CommitMessage string `json:"commit_message"`
}

func (g *GitStateStruct) populate_buffer_store(n int) {
	if n < 1 {
		return
	}
	if g.buffer_store == nil {
		g.buffer_store = make(chan *filehandler.BytesStore)
	}

	for n > 0 {
		g.buffer_store <- filehandler.NewBytesStore()
		n -= 1
	}
}

func (g *GitStateStruct) HashEngine(ctx ftp_context.Context, n int) <-chan ftp_context.LogErr {
	err_ := make(chan ftp_context.LogErr, 1)

	for n > 0 {
		e_ := g.hash_engine(ctx)
		n -= 1
		ctx.Add()
		go func(ec <-chan ftp_context.LogErr) {
			defer ctx.Finished()
			var e ftp_context.LogErr
			for ok := true; ok; {
				select {
				case <-ctx.Done():
					return
				case e, ok = <-ec:
					err_ <- e
				}
			}
		}(e_)
	}

	return err_
}

func (g *GitStateStruct) hash_engine(ctx ftp_context.Context) (err_c <-chan ftp_context.LogErr) {
	loc := "hash_engine"
	ctx.Add()
	ctx.Add()
	defer ctx.Finished()
	err_ := make(chan ftp_context.LogErr, 1)

	g.populate_buffer_store(10)
	go func() {
		defer ctx.Finished()
	main_loop:
		for {
			select {
			case <-ctx.Done():
				break main_loop

			case n := <-g.FileTree.DequeueHashing():
				f, ok := g.FileTree.HashQueue.Get(n)
				if !ok {
					g.FileTree.EnqueueHashing(f)
					break
				}
				ctx.Add()
				go func(f *filehandler.FileHash) {

					defer ctx.Finished()

					b := <-g.buffer_store
					b.Reset()
					_, err := b.ReadFrom(f)
					if err != nil {
						err_ <- ftp_context.NewLogItem(loc, true).Set("after", "BytesStore.ReadFrom(FileHash)").AppendParentError(err)

						g.FileTree.EnqueueHashing(f)
						g.buffer_store <- b
						return
					}
					f.Hash, err = b.Hash()
					if err != nil {
						err_ <- ftp_context.NewLogItem(loc, true).Set("after", "BytesStore.Hash").AppendParentError(err)

						g.FileTree.EnqueueHashing(f)
						g.buffer_store <- b
						return
					}

					g.FileTree.MarkDoneHashing(n)
					g.FileTree.EnqueueUpload(f)
				}(f)
			}
		}

	}()

	return err_
}
