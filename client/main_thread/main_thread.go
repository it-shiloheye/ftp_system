package mainthread

import (
	"context"
	"net/http"
	"strings"

	"log"
	"time"

	initialiseclient "github.com/it-shiloheye/ftp_system/client/init_client"
	// "github.com/it-shiloheye/ftp_system/client/main_thread/actions"
	dir_handler "github.com/it-shiloheye/ftp_system/client/main_thread/filehandler"
	netclient "github.com/it-shiloheye/ftp_system/client/main_thread/network_client"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	// filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var ClientConfig = initialiseclient.ClientConfig

func ticker(loc string, i int) {

	log.Println(loc, i)
}

func MainThread(ctx ftp_context.Context) context.Context {
	loc := "MainThread(ctx ftp_context.Context) context.Context "

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
		log.Println("starting client cycle")
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
			log.Println("error occured, shutdown: ", ClientConfig.StopOnError)
			if ClientConfig.StopOnError {
				log.Fatalln(err1.Error())
			}
			log.Println(err1.Error())
			_, ok = <-tckr.C
			continue
		}

		ticker(loc, 3)
		log.Println("to rehash:\n", strings.Join(rd.ToRehash, "\n"))
		log.Println("to upload:\n", strings.Join(rd.ToUpload, "\n"))

		// child_ctx.Cancel()
		select {
		case _, ok = <-ctx.Done():

		case <-child_ctx.Done():
			log.Println("new tick")

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
	log.Println("test_server_connection")
	rc := netclient.Route{
		BaseUrl:  host,
		Pathname: "/ping",
	}
	_, err1 := netclient.MakeGetRequest(client, rc, &tsc.tmp)
	if err1 != nil {
		log.Println("error here")
		tsc.count += 1
		if tsc.count < 5 {
			log.Println(err1.Error())
		} else {
			log.Fatalln(err1.Error())
		}
		<-tsc.tc.C
		test_server_connection(client, host, tsc)
	}

	log.Println("server connected successfully:", host)
}
func UpdateFileTree(ctx ftp_context.Context) {
	defer ctx.Finished()
	tc := time.NewTicker(time.Minute)
	for ok := true; ok; {
		select {
		case <-tc.C:
		case _, ok = <-ctx.Done():
		}

		err := dir_handler.WriteFileTree(ctx)
		if err != nil {
			log.Println(err)
		}
		log.Println("updated filetree successfully")
	}
}
