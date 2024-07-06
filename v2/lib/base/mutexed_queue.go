package base

import (
	"fmt"
	"sync"
	"time"
)

type MutexedQueue[T any] struct {
	sync.RWMutex
	Queue   MutexedMap[T] `json:"queue"`
	Done    JsonAtomicInt `json:"done"`
	Pending JsonAtomicInt `json:"pending"`
	queue_c chan int64
}

func NewMutexedQueue[T any]() (Mq *MutexedQueue[T]) {
	Mq = &MutexedQueue[T]{
		Queue:   NewMutexedMap[T](),
		queue_c: make(chan int64),
		Done:    JsonAtomicInt{},
		Pending: JsonAtomicInt{},
	}

	Mq.Done.Store(0)
	Mq.Pending.Store(0)
	return
}

func str_conv(n int64) string {
	return fmt.Sprintf("%d", n)
}

func (mq *MutexedQueue[T]) Enqueue(item T) (retry bool) {
	n := mq.Pending.Load()
	if !mq.Pending.CompareAndSwap(n, n+1) {
		<-time.After(time.Microsecond * 60)

		return mq.Enqueue(item)
	}
	mq.Queue.Set(str_conv(n), item)
	mq.queue_c <- n
	return false
}

func (mq *MutexedQueue[T]) Dequeue() <-chan int64 {
	return mq.queue_c
}

func (mq *MutexedQueue[T]) Len() (n int64) {

	return mq.Queue.count.Load()
}

func (mq *MutexedQueue[T]) Clear() {
	mq.Lock()
	mq.Queue.Clear()
	mq.queue_c = make(chan int64)
	mq.Done.Store(0)
	mq.Pending.Store(0)
	mq.Unlock()

}

func (mq *MutexedQueue[T]) Get(n int64) (it T, ok bool) {
	if n >= mq.Len() || n < 0 {
		return
	}
	return mq.Queue.Get(str_conv(n))
}

func (mq *MutexedQueue[T]) MarkDone(n int64) {
	mq.Done.Add(1)
	mq.Queue.Delete(str_conv(n))
}
