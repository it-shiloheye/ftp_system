package base

import "bytes"

type BufferPool struct {
	bs chan *bytes.Buffer
}

func (bp *BufferPool) Get() *bytes.Buffer {

	b := <-bp.bs
	return b
}

func (bp *BufferPool) Return(b *bytes.Buffer) {
	b.Reset()
	bp.bs <- b
}

func NewBufferPool(n int) *BufferPool {

	ab := &BufferPool{
		bs: make(chan *bytes.Buffer, n+1),
	}

	for i := 0; i < (n + 1); i++ {
		ab.bs <- &bytes.Buffer{}
	}

	return ab
}

type BytesPool struct {
	bs chan *[]byte
}

func (bp *BytesPool) Get() *[]byte {

	b := <-bp.bs
	return b
}

func (bp *BytesPool) Return(b *[]byte) {
	clear(*b)

	bp.bs <- b
}

func NewBytesPool(n int) *BytesPool {

	ab := &BytesPool{
		bs: make(chan *[]byte, n+1),
	}

	for i := 0; i < (n + 1); i++ {
		ab.bs <- &[]byte{}
	}

	return ab
}
