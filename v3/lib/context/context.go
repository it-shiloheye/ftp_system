package ftp_context

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/it-shiloheye/ftp_system/v3/lib/base"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

type Context = *ContextStruct

type ContextStruct struct {
	parent_ctx context.Context
	created    time.Time

	cancel_count atomic.Int64
	cancel_c     chan struct{}

	deadline       time.Time
	valid_deadline atomic.Bool

	err error

	*base.MutexedMap[any]

	wg sync.WaitGroup
}

func init() {

	var _ context.Context = (*ContextStruct)(nil)

}

func CreateNewContext() *ContextStruct {
	return &ContextStruct{
		created:    time.Now(),
		MutexedMap: base.NewMutexedMap[any](),
		cancel_c:   make(chan struct{}),
	}
}

func CreateNewContextWithDeadline(t time.Duration) (ctx *ContextStruct) {
	ctx = CreateNewContext()
	ctx.SetDeadline(t)

	return
}

func CreateNewContextWithParent(pctx context.Context) (ctx *ContextStruct) {

	ctx = CreateNewContext()
	ctx.parent_ctx = pctx

	go func() {
		<-pctx.Done()
		ctx.Cancel()
	}()

	return
}

// returns the time context was created
func (ctx *ContextStruct) Created() time.Time {
	return ctx.created
}

// returns whether the current context is before the deadline
func (ctx *ContextStruct) BeforeDeadline() (has_deadline, before_deadline bool) {
	t_d, ok_1 := ctx.Deadline()
	if ok_1 {
		return ok_1, time.Now().Before(t_d)
	}
	return
}

// returns whether the current context is near or past the deadline
func (ctx *ContextStruct) NearDeadline(t time.Duration) (has_deadline, near_deadline bool) {
	t_d, ok_1 := ctx.Deadline()
	if ok_1 {
		return ok_1, t_d.Sub(time.Now()) < t
	}
	return
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
		ok = ctx.valid_deadline.Load()

	}
	return
}

func (ctx *ContextStruct) Value(key any) any {
	return key
}

func (ctx *ContextStruct) Err() error {
	return ctx.err
}

func (ctx *ContextStruct) Cancel() {
	defer func() {
		recover()
	}()
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
		if cptx, ok := ctx.parent_ctx.(Context); ok {
			cptx.Add()
		}

	}
	ctx.wg.Add(1)

	return ctx
}

// waits for all goroutines to return and cancel to be called
func (ctx *ContextStruct) Wait() <-chan struct{} {
	ctx.wg.Wait()

	return ctx.cancel_c // close any process that have not returned
}

// marks goroutine as finished
func (ctx *ContextStruct) Finished() {
	if ctx.parent_ctx != nil {
		if cptx, ok := ctx.parent_ctx.(Context); ok {
			cptx.Finished()
		}

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
	if cptx, _ok := ctx.parent_ctx.(Context); _ok {
		it, ok = cptx.Get(key)
	} else {
		it = ctx.parent_ctx.Value(key)
		ok = it != nil
	}

	return
}

// Sets current context and parent context
func (ctx *ContextStruct) SetParent(key string, val any) error {
	ctx.Set(key, val)

	if ctx.parent_ctx == nil {
		return log_item.NewLogItem("ctx.SetParent", log_item.LogLevelError01).SetMessagef("parent context is nil:\nkey: %s\nval: %v", key, val)
	}
	if cptx, ok := ctx.parent_ctx.(Context); ok {
		cptx.SetParent(key, val)
	}

	return nil
}
