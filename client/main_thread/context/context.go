package ftp_context

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ftp_system_client/base"
)

type Context = *ContextStruct

type ContextStruct struct {
	parent_ctx *ContextStruct

	created time.Time

	cancel_count atomic.Int64
	cancel_c     chan struct{}

	deadline       time.Time
	valid_deadline atomic.Bool

	err error

	base.MutexedMap[any]

	wg sync.WaitGroup
}

func init() {

	var _ context.Context = (*ContextStruct)(nil)

}

func CreateNewContext() *ContextStruct {
	return &ContextStruct{
		created:    time.Now(),
		MutexedMap: base.NewMutexedMap[any](),
	}
}

func CreateNewContextWithDeadline(t time.Duration) (ctx *ContextStruct) {
	ctx = CreateNewContext()
	ctx.SetDeadline(t)

	return
}

func CreateNewContextWithParent(pctx *ContextStruct) (ctx *ContextStruct) {

	ctx = CreateNewContext()
	ctx.parent_ctx = pctx

	go func() {
		<-pctx.Done()
		ctx.Cancel()
	}()

	return
}

func (ctx *ContextStruct) Created() time.Time {
	return ctx.created
}

func (ctx *ContextStruct) Done() <-chan struct{} {
	if ctx.cancel_c == nil {
		ctx.cancel_c = make(chan struct{})
	}

	ctx.cancel_count.Add(1)

	return ctx.cancel_c
}

func (ctx *ContextStruct) Deadline() (t time.Time, ok bool) {
	if ctx.valid_deadline.Load() {
		t = ctx.deadline
		ok = true

	}
	return
}

func (ctx *ContextStruct) Value(key any) any {
	panic("don't use this")
}

func (ctx *ContextStruct) Err() error {
	return ctx.err
}

func (ctx *ContextStruct) Cancel() {
	close(ctx.cancel_c)
}

func (ctx *ContextStruct) SetDeadline(t time.Duration) (deadline time.Time) {

	ctx.Lock()
	deadline = (time.Now().Add(t))
	ctx.deadline = deadline
	ctx.valid_deadline.Store(true)

	ctx.Unlock()
	k := time.After(t)

	go func() {
		<-k
		ctx.Lock()
		if ctx.valid_deadline.Load() {
			ctx.Cancel()
		}
		ctx.Unlock()
	}()
	return
}

func (ctx *ContextStruct) CancelDeadline() (ok bool) {

	if time.Now().Before(ctx.deadline) {
		ctx.valid_deadline.Store(false)
		return false
	}

	return
}

// adds new goroutine to waitgroup
func (ctx *ContextStruct) Add() Context {
	if ctx.parent_ctx != nil {
		ctx.parent_ctx.Add()

	}
	ctx.wg.Add(1)

	return ctx
}

// does not unblock on cancel
func (ctx *ContextStruct) Wait() {
	ctx.wg.Wait()
	ctx.Cancel() // close any process that have not returned
}

// marks goroutine as finished
func (ctx *ContextStruct) Finished() {
	if ctx.parent_ctx != nil {
		ctx.parent_ctx.Finished()
	}
	ctx.wg.Done()
}

func (ctx *ContextStruct) NewChild() Context {

	return CreateNewContextWithParent(ctx)
}

func (ctx *ContextStruct) Get(key string) (it any, ok bool) {

	it, ok = ctx.MutexedMap.Get(key)
	if ok {
		return
	}
	it, ok = ctx.parent_ctx.Get(key)
	return
}

// Sets current context and parent context
func (ctx *ContextStruct) SetParent(key string, val any) error {
	ctx.Set(key, val)

	if ctx.parent_ctx == nil {
		return NewLogItem("ctx.SetParent", true).SetMessagef("parent context is nil:\nkey: %s\nval: %v", key, val)
	}
	ctx.parent_ctx.SetParent(key, val)
	return nil
}
