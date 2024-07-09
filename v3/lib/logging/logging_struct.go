package logging

import (
	"errors"
	"fmt"

	"io/fs"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/it-shiloheye/ftp_system/v3/lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	"github.com/it-shiloheye/ftp_system/v3/lib/file_handler/v2"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
)

const fs_mode = fs.FileMode(base.S_IRWXU | base.S_IRWXO)

var Logger = &LoggerStruct{

	log_file:           &filehandler.FileBasic{},
	log_err_file:       &filehandler.FileBasic{},
	log_today_file:     &filehandler.FileBasic{},
	log_today_err_file: &filehandler.FileBasic{},
	comm:               make(chan *log_item.LogItem, 100),
}

var lock = &sync.Mutex{}

type LoggerStruct struct {
	comm               chan *log_item.LogItem
	log_file           *filehandler.FileBasic
	log_err_file       *filehandler.FileBasic
	log_today_file     *filehandler.FileBasic
	log_today_err_file *filehandler.FileBasic
}

func InitialiseLogging(logging_dir string) {
	loc := log_item.Locf("func InitialiseLogging(logging_dir: %s)", logging_dir)
	log.Println("loading logger")
	log.Println(logging_dir, "logging directory")

	log_file_p := logging_dir + "/logs/log_file.txt"
	log_err_file_p := logging_dir + "/logs/log_err_file.txt"
	log_today_file_p := logging_dir + "/logs/sess/" + log_file_name() + ".txt"
	log_today_err_file_p := logging_dir + "/logs/sess/" + log_file_name() + "_err.txt"

	log.Printf("%s\nlog_file_p: %s\nlog_err_file_p: %s\n", loc, log_file_p, log_err_file_p)
	// os.Exit(1)
	var err1, err2, err3, err4, err5 error

	err1 = os.MkdirAll(logging_dir+"/logs/sess", fs.FileMode(base.S_IRWXO|base.S_IRWXU))
	if !errors.Is(err1, os.ErrExist) && err1 != nil {
		a := &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			Message:  err1.Error(),
			Level:    log_item.LogLevelError02, CallStack: []error{err1},
		}
		log.Fatalln(a)
	}

	Logger.log_file, err2 = filehandler.Create(log_file_p)
	if err2 != nil {
		b := &log_item.LogItem{
			Location: loc,
			Time:     time.Now(),
			Message:  err2.Error(),

			CallStack: []error{err2},
		}
		log.Fatalln(b)
	}

	Logger.log_err_file, err3 = filehandler.Create(log_err_file_p)
	if err3 != nil {
		c := &log_item.LogItem{
			Location:  loc,
			Time:      time.Now(),
			Message:   err3.Error(),
			Level:     log_item.LogLevelError02,
			CallStack: []error{err3},
		}
		log.Fatalln(c)
	}

	Logger.log_today_file.File, err4 = base.OpenFile(log_today_file_p, os.O_APPEND|os.O_RDWR|os.O_CREATE)
	if err4 != nil {
		c := &log_item.LogItem{
			Location:  loc,
			Time:      time.Now(),
			Message:   err4.Error(),
			Level:     log_item.LogLevelError02,
			CallStack: []error{err4},
		}
		log.Fatalln(c)
	}

	Logger.log_today_err_file.File, err5 = base.OpenFile(log_today_err_file_p, os.O_APPEND|os.O_RDWR|os.O_CREATE)
	if err5 != nil {
		c := &log_item.LogItem{
			Location:  loc,
			Time:      time.Now(),
			Message:   err5.Error(),
			Level:     log_item.LogLevelError02,
			CallStack: []error{err5},
		}
		log.Fatalln(c)
	}

	log.Println("successfull loaded logger")
}

func (ls *LoggerStruct) Log(li *log_item.LogItem) {

	ls.comm <- li
}

func (ls *LoggerStruct) Logf(loc log_item.Loc, str string, v ...any) {
	ls.comm <- &log_item.LogItem{
		Location: loc,
		Time:     time.Now(),
		Message:  fmt.Sprintf(str, v...),
	}
}

func (ls *LoggerStruct) LogErr(loc log_item.Loc, err error) error {
	e := &log_item.LogItem{
		Location:  loc,
		Time:      time.Now(),
		Level:     log_item.LogLevelError02,
		Message:   err.Error(),
		CallStack: []error{err},
	}
	ls.Log(e)
	return e
}

func (ls *LoggerStruct) Engine(ctx ftp_context.Context, logging_dir string) {
	lock.Lock()
	defer ctx.Finished()
	defer lock.Unlock()

	lock_file := logging_dir + "/log.lock"

	tc := time.NewTicker(time.Second)

	var li *log_item.LogItem

	queue := []*log_item.LogItem{}

	var txt, log_txt, err_txt string
	int_ := 0
	for ok := true; ok; {
		log_txt, err_txt = "", ""
		int_ += 1
		select {
		case _, ok = <-ctx.Done():

		case li = <-ls.comm:
			if li != nil {
				queue = append(queue, li)
			}
			// log.Println(int_, "none: new li")
			continue
		case <-tc.C:
		}

		for _, li := range queue {
			if li == nil {
				continue
			}

			txt = li.String() + "\n"
			if len(txt) < 2 {
				continue
			}
			log_txt = log_txt + txt
			if li.Level > log_item.LogLevelInfo03 {
				err_txt = err_txt + txt
				log.SetOutput(os.Stderr)
				log.Print(txt)
			} else {
				log.SetOutput(os.Stdout)
				log.Printf("\n%s:\n%s", li.Location, li.Message)
			}

		}
		l, err1 := filehandler.Lock(lock_file)
		if err1 != nil && !errors.Is(err1, os.ErrExist) {
			log.Println("err:\n", err1)
			<-time.After(time.Second * 5)
			continue
		}

		if len(log_txt) > 0 {
			log.SetOutput(ls.log_today_file)
			log.Print(log_txt)

			log.SetOutput(ls.log_file)
			log.Print(log_txt)
		}

		if len(err_txt) > 0 {
			log.SetOutput(ls.log_err_file)
			log.Print(err_txt)

			log.SetOutput(ls.log_today_err_file)
			log.Print(err_txt)

		}

		clear(queue)
		l.Unlock()
		log.SetOutput(os.Stdout)
		// log.Println(int_, "none: done engine")
	}

}

func log_file_name() string {

	d := time.Now().Format(time.RFC3339)
	d1 := strings.ReplaceAll(d, " ", "_")
	d2 := strings.ReplaceAll(d1, ":", "")
	d3 := strings.Split(d2, "T")
	return d3[0]
}
