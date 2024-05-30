package filehandler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"

	ftp_context "github.com/ftp_system_client/main_thread/context"
)

type BytesStore struct {
	h hash.Hash
	bytes.Buffer
}

func (bs *BytesStore) Hash() (hash string, err error) {
	loc := "bs.Hash"
	bs.h.Reset()
	_, err = io.Copy(bs.h, &bs.Buffer)
	if err != nil {
		err = ftp_context.NewLogItem(loc, true).Set("after", "io.Copy").AppendParentError(err)
		return
	}
	hash = fmt.Sprintf("%x", bs.h.Sum(nil))

	return
}

func (bs *BytesStore) CopyFrom(fo io.Reader) (n int64, err error) {
	n, err = io.Copy(&bs.Buffer, fo)
	return
}

func (bs *BytesStore) CopyTo(fo io.Writer) (n int64, err error) {
	n, err = io.Copy(fo, &bs.Buffer)
	return
}

func NewBytesStore() (bs *BytesStore) {
	bs = &BytesStore{
		h:      sha256.New(),
		Buffer: bytes.Buffer{},
	}
	bs.Grow(100_000)
	return
}
