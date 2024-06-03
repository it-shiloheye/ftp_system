package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/it-shiloheye/ftp_system_lib/base"
)

const (
	log_folder = "./logs"
)

var Logger = &LoggerStruct{
	wg: &sync.WaitGroup{},
	c:  make(chan string),
}

type LoggerStruct struct {
	wg *sync.WaitGroup
	c  chan string
}

func (d *LoggerStruct) Close() {
	close(d.c)
}

func (d *LoggerStruct) Add() {
	d.wg.Add(1)
}

func (d *LoggerStruct) LogAppendFile(fpath string, name string, ftype string) {
	d.c <- fmt.Sprintf("path: %s\nname: %s\ntype: %s", fpath, name, ftype)
}

func (d *LoggerStruct) Log(v ...any) string {
	s := fmt.Sprint(v...)
	d.c <- s
	return s
}
func (d *LoggerStruct) Logf(str string, v ...any) string {
	s := fmt.Sprintf(str, v...)
	d.c <- s
	return s
}

type ErrorLocation string

func (d *LoggerStruct) LogErr(loc ErrorLocation, v ...any) string {
	s := fmt.Sprintf("%s error:\n%s", loc, fmt.Sprint(v...))
	d.c <- s
	return s
}
func (d *LoggerStruct) LogErrf(loc ErrorLocation, str string, v ...any) string {
	s := fmt.Sprintf("%s error:\n%s", loc, fmt.Sprintf(str, v...))
	d.c <- s
	return s
}

func (d *LoggerStruct) Done() {
	d.wg.Done()
}

func (d *LoggerStruct) Go(ctx context.Context) (err error) {

	defer d.wg.Done()
	ok := true
	tmr := time.NewTicker(time.Second)
	log_list := []string{}
	it := ""

	ex_, _ := exists(log_folder)
	if !ex_ {
		err = os.MkdirAll(log_folder, os.ModeDir)
		if err != nil {
			log.Printf("LoggerStruct.Go error: %s", err.Error())
			return
		}
	}
	var lg_file *os.File
	lg_file, err = base.OpenFile(log_folder+"/"+log_file_name(), os.O_CREATE|os.O_APPEND|os.O_RDWR)
	if err != nil {
		log.Printf("base.OpenFile error: %s", err.Error())
		return
	}

	for ok {
		select {
		case <-ctx.Done():
			ok = false
		case <-tmr.C:

		case it, ok = <-d.c:
			log_list = append(log_list, it)
			if ok {
				continue
			}
		}

		for _, it = range log_list {
			if len(it) < 1 {
				continue
			}
			log.Println(it)
			log.SetOutput(lg_file)
			log.Println(it)
			log.SetOutput(os.Stdout)
		}

		clear(log_list)
	}

	return
}

func (d *LoggerStruct) Wait() {
	d.wg.Wait()
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func log_file_name() string {

	d := time.Now().Format(time.RFC3339)
	d1 := strings.ReplaceAll(d, " ", "_")
	d2 := strings.ReplaceAll(d1, ":", "")
	d3 := strings.Split(d2, "T")
	return d3[0]
}
