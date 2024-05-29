package base

import (
	"sync"
	"sync/atomic"
	"time"
)

type MutexedQueue[T any] struct {
	sync.RWMutex
	Queue []T `json:"queue"`
	count atomic.Int64
}

func NewMutexedQueue[T any]() (Mq *MutexedQueue[T]) {
	Mq = &MutexedQueue[T]{
		Queue: make([]T, 10),
		count: atomic.Int64{},
	}

	Mq.count.Store(0)
	return
}

func (mq *MutexedQueue[T]) Enqueue(item T) {
	mq.Lock()
	mq.Queue = append(mq.Queue, item)
	mq.Unlock()

}

func (mq *MutexedQueue[T]) Dequeue() <-chan T {
	c := make(chan T, 1)
	n := mq.count.Load()
	for !mq.count.CompareAndSwap(n, n+1) {
		<-time.After(time.Microsecond * 10)
		n = mq.count.Load()

	}

	mq.RLock()
	t := mq.Queue[n]
	mq.RUnlock()
	c <- t
	close(c)
	return c
}

func (mq *MutexedQueue[T]) Len() (n int64) {
	mq.RLock()
	n = int64(len(mq.Queue))
	mq.RUnlock()

	return
}

func (mq *MutexedQueue[T]) Clear() {
	mq.Lock()
	clear(mq.Queue)
	mq.count.Store(0)
	mq.Unlock()

}

func (mq *MutexedQueue[T]) Pos() (n int64) {
	return mq.count.Load()
}

func (mq *MutexedQueue[T]) Get(n int) (it T, ok bool) {
	if n >= int(mq.Len()) {
		return
	}
	mq.RLock()
	it = mq.Queue[n]
	mq.RUnlock()

	return
}
