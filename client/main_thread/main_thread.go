package mainthread

import (
	"context"
	"net/http"
	"os"
	// "strings"

	"log"
	"time"

	initialiseclient "github.com/it-shiloheye/ftp_system/client/init_client"
	// "github.com/it-shiloheye/ftp_system/client/main_thread/actions"
	dir_handler "github.com/it-shiloheye/ftp_system/client/main_thread/filehandler"
	"github.com/it-shiloheye/ftp_system/client/main_thread/logging"
	netclient "github.com/it-shiloheye/ftp_system/client/main_thread/network_client"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	// filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
	// filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var ClientConfig = initialiseclient.ClientConfig
var Logger = logging.Logger
var FileTree = dir_handler.FileTree

func ticker(loc string, i int) {

	// Logger.Logf(loc, "%d", i)
}

func MainThread(ctx ftp_context.Context) context.Context {
	loc := "MainThread(ctx ftp_context.Context) context.Context "

	lock, ERR := dir_handler.Lock(ClientConfig.DataDir + "/index.lock")
	defer lock.Unlock()
	if ERR != nil {
		log.Println(ERR)
		log.Fatalln("cannot obtain lock on data/dir")
	}

	ticker(loc, 1)
	defer ctx.Wait()

	if len(ClientConfig.IncludeDir) < 1 {
		if len(ClientConfig.DirConfig.Path) > 0 {
			ClientConfig.IncludeDir = append(ClientConfig.IncludeDir, ClientConfig.DirConfig.Path)
		} else {
			log.Fatalln("add at least one file to include list or directory path")
		}
	}

	if ClientConfig.UpdateRate < 1 {
		ClientConfig.UpdateRate = time.Minute * 5
	} else {
		ClientConfig.UpdateRate = time.Duration(ClientConfig.UpdateRate)
	}

	tick := ClientConfig.UpdateRate
	tckr := time.NewTicker(ClientConfig.UpdateRate)

	client, err_ := netclient.NewNetworkClient(ctx)
	if err_ != nil {
		log.Fatalln(err_)
	}
	base_server := "https://127.0.0.1:8080"
	tyc := &TestServerConnection{
		tmp: map[string]string{},
		tc:  time.NewTicker(time.Second * 5),
	}

	test_server_connection(client, base_server, tyc)

	go UpdateFileTree(ctx.Add())
	for ok := true; ok; {

		child_ctx := ctx.NewChild()
		child_ctx.SetDeadline(tick)
		Logger.Logf(loc, "starting client cycle")
		/**
		* five tasks:
		*	1. Read all files in directory
		*		- list all files (exclude .git) [done]
		*		- create a printout of list of files (current timestamped - incase of crash)
		* 	2. Check for any changes in directory compared to last scan
		*		- store past "ModTime" in special format
		*		- compare present and past mod-time for changes
		*	3. Add and commit all changes
		*   4. Hash all files in .git folder
		*		- read all files in .git
		*		- check for any changes in mod time (or new files)
		*		- hash where necessary
		*	5. Transmit over network any new changes where necessary
		 */

		rd, err1 := dir_handler.ReadDir(child_ctx.Add(), ClientConfig.DirConfig)
		ticker(loc, 2)
		if err1 != nil {
			Logger.Logf(loc, "error occured, shutdown: %s", ClientConfig.StopOnError)
			if ClientConfig.StopOnError {
				log.Fatalln(err1.Error())
			}
			Logger.LogErr(loc, err1)
			_, ok = <-tckr.C
			continue
		}

		for _, file_ := range rd.ToRehash {
			Logger.Logf("hashing: %s", file_)
		}
		// child_ctx.Cancel()
		select {
		case _, ok = <-ctx.Done():

		case <-child_ctx.Done():
			Logger.Logf(loc, "new tick")

		}
	}

	return ctx
}

type TestServerConnection struct {
	tmp   map[string]string
	tc    *time.Ticker
	count int
}

func test_server_connection(client *http.Client, host string, tsc *TestServerConnection) {
	loc := " test_server_connection(client *http.Client, host string, tsc *TestServerConnection)"
	Logger.Logf(loc, "test_server_connection")
	rc := netclient.Route{
		BaseUrl:  host,
		Pathname: "/ping",
	}
	_, err1 := netclient.MakeGetRequest(client, rc, &tsc.tmp)
	if err1 != nil {
		Logger.Logf(loc, "error here")
		tsc.count += 1
		if tsc.count < 5 {
			Logger.LogErr(loc, err1)

		} else {
			Logger.LogErr(loc, err1)
			os.Exit(1)
		}
		<-tsc.tc.C
		test_server_connection(client, host, tsc)
	}

	Logger.Logf(loc, "server connected successfully: %s", host)
}
func UpdateFileTree(ctx ftp_context.Context) {
	loc := "UpdateFileTree(ctx ftp_context.Context)"
	defer ctx.Finished()
	tc := time.NewTicker(time.Minute)
	for ok := true; ok; {
		select {
		case <-tc.C:
		case _, ok = <-ctx.Done():
		}

		err := dir_handler.WriteFileTree(ctx)
		if err != nil {
			Logger.LogErr(loc, err)
		}
		Logger.Logf(loc, "updated filetree successfully")
	}
}

func hashing_piston(ctx ftp_context.Context, to_hash string) (done_hashing string, err error) {
	loc := "hashing_piston(ctx ftp_context.Context, to_hash string) (done_hashing string, err error)"
	Logger.Logf(loc, "entering hashing piston, hashing: %s", to_hash)

	Logger.Logf(loc, "exiting hashing piston, successful hashing: %s", to_hash)
	return

}
