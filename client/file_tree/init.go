package filetree

import (
	"crypto/sha256"
	"log"
)

var bytestore_queue chan *BytesStore = make(chan *BytesStore, 10)

func init() {

	for i := 0; i < 10; i++ {
		bytestore_queue <- &BytesStore{
			bytes: make([]byte, 100_000_00),
			h:     sha256.New(),
		}
	}

	err := FilesList.readFileTreeJSON(file_tree_file)
	if err != nil {
		log.Fatalln(err)
	}
}
