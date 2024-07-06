package filehandler

import (
	"errors"
	"fmt"

	"os"
	"time"

	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

type LockFile struct {
	locked bool
	Name   string
}

func Lock(file_path string) (lf *LockFile, err error) {
	err1 := os.MkdirAll(file_path, os.FileMode(os.ModeExclusive))
	if err1 != nil {
		if errors.Is(err1, os.ErrExist) {
			err = err1
			return
		}
		err = &log_item.LogItem{
			Location: log_item.Locf(`Lock("%s" string) (lf *LockFile, err error)`, file_path),
			Time:     time.Now(),
			After:    fmt.Sprintf(`err1 := os.MkdirAll("%s", os.FileMode(os.ModeExclusive))`, file_path),
			Message:  err1.Error(),
		}
		return
	}

	lf = &LockFile{
		Name:   file_path,
		locked: true,
	}

	return
}

func (lf *LockFile) Unlock() error {
	err1 := os.Remove(lf.Name)
	if err1 == nil || errors.Is(err1, os.ErrNotExist) {
		lf.locked = false
		return nil
	}

	err2 := &log_item.LogItem{
		Location: log_item.Loc(`func (lf *LockFile) Unlock() error`),
		Time:     time.Now(),
		After:    fmt.Sprintf(`err1 := os.Remove("%s")`, lf.Name),
		Message:  err1.Error(),
	}
	return err2
}

func (lf *LockFile) Lock() (err error) {
	if lf.locked {
		return fmt.Errorf("already locked")
	}
	err1 := os.MkdirAll(lf.Name, os.FileMode(os.ModeExclusive))
	if err1 != nil {
		if errors.Is(err1, os.ErrExist) {

			err = err1
			return
		}
		err = &log_item.LogItem{
			Location: log_item.Locf(`Lock("%s" string) (lf *LockFile, err error)`, lf.Name),
			Time:     time.Now(),
			After:    fmt.Sprintf(`err1 := os.MkdirAll("%s", os.FileMode(os.ModeExclusive))`, lf.Name),
			Message:  err1.Error(),
		}
		return
	}

	lf.locked = true

	return
}
