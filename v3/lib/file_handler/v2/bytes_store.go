package filehandler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

type BytesStore struct {
	h hash.Hash
	bytes.Buffer
}

func (bs *BytesStore) Hash() (hash string, err error) {
	loc := log_item.Loc("func (bs *BytesStore) Hash() (hash string, err error)")
	bs.h.Reset()
	_, err1 := bs.WriteTo(bs.h)
	if err1 != nil {
		err = log_item.NewLogItem(loc, log_item.LogLevelError01).
			SetAfter("_, err = bs.CopyTo(bs.h)").
			SetMessage(err1.Error()).
			AppendParentError(err1)
		return
	}
	hash = fmt.Sprintf("%x", bs.h.Sum(nil))

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
