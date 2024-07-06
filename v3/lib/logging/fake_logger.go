package logging

import (
	"fmt"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type FakeLogger struct {
	C chan *log_item.LogItem

	log_level atomic.Pointer[log_item.LogLevel]

	out  io.Writer
	lock sync.Mutex

	prefix atomic.Pointer[string] // prefix on each line to identify the logger (but see Lmsgprefix)
	flag   atomic.Int32           // properties

}

func (l *FakeLogger) Fatal(v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelFatal,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp

}
func (l *FakeLogger) Fatalf(format string, v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelFatal,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Fatalln(v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelFatal,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Flags() int {
	return int(l.flag.Load())
}
func (l *FakeLogger) Output(calldepth int, s string) error {
	return nil
}
func (l *FakeLogger) Panic(v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelError02,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
	panic(tmp)
}
func (l *FakeLogger) Panicf(format string, v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelError02,
		Message: *l.prefix.Load() + fmt.Sprintf(format, v...),
	}
	l.C <- tmp
	panic(tmp)
}
func (l *FakeLogger) Panicln(v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelError02,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
	panic(tmp)
}
func (l *FakeLogger) Prefix() string {
	return *l.prefix.Load()
}
func (l *FakeLogger) Print(v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelInfo01,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Printf(format string, v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelInfo01,
		Message: *l.prefix.Load() + fmt.Sprintf(format, v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) Println(v ...any) {
	tmp := &log_item.LogItem{
		Time:    time.Now(),
		Level:   log_item.LogLevelInfo01,
		Message: *l.prefix.Load() + fmt.Sprint(v...),
	}
	l.C <- tmp
}
func (l *FakeLogger) SetFlags(flag int) {
	l.flag.Store(int32(flag))
}
func (l *FakeLogger) SetOutput(w io.Writer) {
	l.lock.Lock()
	l.out = w
	l.lock.Unlock()
}
func (l *FakeLogger) SetPrefix(prefix string) {
	l.prefix.Store(&prefix)
}
