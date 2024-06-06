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

	initialiseclient "github.com/it-shiloheye/ftp_system/client/init_client"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	"github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var Logger = &LoggerStruct{
	lock:  sync.Mutex{},
	comm:  make(chan *ftp_context.LogItem, 100),
	err_c: make(chan error, 100),
}
var ClientConfig = &initialiseclient.ClientConfigStruct{}

type LoggerStruct struct {
	lock  sync.Mutex
	comm  chan *ftp_context.LogItem
	err_c chan error
}

var log_file = &filehandler.FileBasic{}
var log_err_file = &filehandler.FileBasic{}

func init() {
	log.Println("loading logger")

	loc := "ftp_system/client/main_thread/logging/logging_struct.go"
	log_file_p := ClientConfig.DataDir + "/log/log_file.txt"
	log_err_file_p := ClientConfig.DataDir + "/log/log_err_file.txt"
	log.Printf("%s\n%s/n%s\n", loc, log_file_p, log_err_file_p)
	var err1, err2, err3 error

	err1 = os.MkdirAll(ClientConfig.DataDir+"/log", fs.FileMode(ftp_base.S_IRWXO|ftp_base.S_IRWXU))
	if !errors.Is(err1, os.ErrExist) && err1 != nil {
		a := &ftp_context.LogItem{
			Location:  loc,
			Time:      time.Now(),
			Message:   err1.Error(),
			CallStack: []error{err1},
		}
		log.Fatalln(a)
	}

	log_file.File, err2 = ftp_base.OpenFile(log_file_p, os.O_APPEND|os.O_RDWR|os.O_CREATE)
	if err2 != nil {
		b := &ftp_context.LogItem{
			Location:  loc,
			Time:      time.Now(),
			Message:   err2.Error(),
			Err:       true,
			CallStack: []error{err2},
		}
		log.Fatalln(b)
	}

	log_err_file.File, err3 = ftp_base.OpenFile(log_err_file_p, os.O_APPEND|os.O_RDWR|os.O_CREATE)
	if err3 != nil {
		c := &ftp_context.LogItem{
			Location:  loc,
			Time:      time.Now(),
			Message:   err2.Error(),
			Err:       true,
			CallStack: []error{err3},
		}
		log.Fatalln(c)
	}

	log.Println("successfull loaded logger")
}

func (ls *LoggerStruct) Log(li *ftp_context.LogItem) {
	if li.Err {
		ls.err_c <- li
	}
	ls.comm <- li
}

func (ls *LoggerStruct) Logf(loc, str string, v ...any) {
	ls.comm <- &ftp_context.LogItem{
		Location: loc,
		Time:     time.Now(),
		Message:  fmt.Sprintf(str, v...),
	}
}

func (ls *LoggerStruct) LogErr(loc string, err error) {
	e := &ftp_context.LogItem{
		Location:  loc,
		Time:      time.Now(),
		Err:       true,
		Message:   err.Error(),
		CallStack: []error{err},
	}
	ls.Log(e)
}

func (ls *LoggerStruct) Engine(ctx ftp_context.Context) {
	ls.lock.Lock()
	defer ctx.Finished()
	defer ls.lock.Unlock()

	tc := time.NewTicker(time.Second)
	var li *ftp_context.LogItem
	var lerr error
	queue := []*ftp_context.LogItem{}
	err_queue := []error{}

	var log_txt, err_txt string
	for ok := true; ok; {
		log_txt, err_txt = "", ""
		select {
		case _, ok = <-ctx.Done():
			break
		case li = <-ls.comm:
			queue = append(queue, li)
			continue

		case lerr = <-ls.err_c:
			err_queue = append(err_queue, lerr)
			continue
		case <-tc.C:
		}

		for _, li := range queue {
			log_txt += li.String() + "\n"
		}

		log.SetOutput(log_file)
		log.Print(log_txt)

		log.SetOutput(os.Stdout)
		log.Print(log_txt)

		for _, li := range err_queue {
			err_txt += li.Error() + "\n"
		}

		log.SetOutput(log_err_file)
		log.Print(err_txt)

		log.SetOutput(os.Stderr)
		log.Print(err_txt)

		clear(queue)
		clear(err_queue)
	}

}

func log_file_name() string {

	d := time.Now().Format(time.RFC3339)
	d1 := strings.ReplaceAll(d, " ", "_")
	d2 := strings.ReplaceAll(d1, ":", "")
	d3 := strings.Split(d2, "T")
	return d3[0]
}
