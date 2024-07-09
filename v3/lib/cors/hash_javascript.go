package cors

import (
	"encoding/base64"
	"fmt"
	"hash"
	"os"
	"strings"

	"golang.org/x/crypto/blake2b"
)

type HashAlgorithm int

const (
	HashAlgo256 HashAlgorithm = 256
	HashAlgo384 HashAlgorithm = 384
	HashAlgo512 HashAlgorithm = 512
)

type HashedFile struct {
	file_name string
	raw       string
	nonce     string
	algo      HashAlgorithm
}

func NewHashedFile(path string, algo HashAlgorithm) (h *HashedFile, err error) {

	f_, err := os.ReadFile(path)
	if err != nil {
		return
	}

	h = &HashedFile{
		file_name: path,
		raw:       shrinkFile(string(f_)), // string() casts to utf-8 string
		algo:      algo,
	}

	h.nonce = hashData(h.raw, algo)
	return

}

func shrinkFile(str string) string {

	if len(str) < 1 {
		panic(fmt.Errorf("no file"))
	}
	return str
}

func hashData(data string, algo HashAlgorithm) string {
	d_ := []byte(data)
	var hash hash.Hash
	var err error
	// Hash the data
	switch algo {
	case HashAlgo256:
		hash, err = blake2b.New256(nil)
	case HashAlgo384:
		hash, err = blake2b.New384(nil)
	case HashAlgo512:
		hash, err = blake2b.New512(nil)
	default:
		panic(fmt.Errorf("invalid hash algorithm chosen: %d", algo))
	}
	if err != nil {
		panic(err)
	}
	hash.Write(d_)
	digest := hash.Sum(nil)
	return fmt.Sprintf("sha%d-%s", int(algo), base64.RawStdEncoding.EncodeToString(digest))
}

func (h *HashedFile) Replace(mstr map[string]string) *HashedFile {

	for key, value := range mstr {
		h.raw = strings.Replace(h.raw, key, value, -1)
	}

	h.raw = shrinkFile(h.raw)
	h.nonce = hashData(h.raw, h.algo)

	return h
}

func (h *HashedFile) Str() string {
	return h.raw
}

func (h *HashedFile) Nonce() string {
	return h.nonce
}

func (h *HashedFile) Reload() (*HashedFile, error) {

	h, err := NewHashedFile(h.file_name, h.algo)

	if err != nil {
		return nil, err
	}

	return h, nil
}
