package ftp_context

import (
	"encoding/json"
	"fmt"
	"time"
)

type LogErr = *LogItem

type LogItem struct {
	Location  string         `json:"location"`
	Time      time.Time      `json:"time"`
	Body      map[string]any `json:"body"`
	Message   string         `json:"message"`
	Err       bool           `json:"is_error"`
	CallStack error          `json:"call_stack"`
}

func (li *LogItem) Error() string {
	if !li.Err {
		return ""
	}
	return li.to_string()
}

func (li *LogItem) to_string() string {
	stp_1 := func() string {
		b, err := json.MarshalIndent(li, "\t", " ")
		if err != nil {
			panic(fmt.Sprint("LogItem.Error json.MarshalIndent ", err))
		}
		return string(b)
	}()

	stp_2 := fmt.Sprintf("%s:\n%s:\n%s", li.Time.Format(time.RFC822Z), li.Location, stp_1)

	stp_3 := func() string {
		if li.Err {
			return fmt.Sprint("[ERR] ", stp_2)
		}
		return fmt.Sprint("[LOG] ", stp_2)
	}()

	return stp_3
}

func (li *LogItem) Set(key string, value any) *LogItem {
	if li.Body == nil {
		li.Body = make(map[string]any)
	}
	li.Body[key] = value

	return li
}

func (li *LogItem) Get(key string) (it any, ok bool) {
	if li.Body == nil {
		li.Body = make(map[string]any)
		return
	}
	it, ok = li.Body[key]
	return
}

func (li *LogItem) ParentError(err error) *LogItem {
	li.Err = true
	li.CallStack = err
	return li
}

func NewLogItem(loc string, err bool) (lt *LogItem) {
	lt = &LogItem{
		Err:      err,
		Location: loc,
		Time:     time.Now(),
	}
	return
}

func (lt *LogItem) SetMessage(v ...any) *LogItem {
	lt.Message = fmt.Sprint(v...)
	return lt
}

func (lt *LogItem) SetMessagef(str string, v ...any) *LogItem {
	lt.Message = fmt.Sprintf(str, v...)
	return lt
}
