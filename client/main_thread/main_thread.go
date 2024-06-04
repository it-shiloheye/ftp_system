package mainthread

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ftp_system_client/main_thread/actions"
	netclient "github.com/ftp_system_client/main_thread/network_client"
	configuration "github.com/it-shiloheye/ftp_system_lib/config"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler"
	githandler "github.com/it-shiloheye/ftp_system_lib/git_handler"
)

func MainThread(ctx ftp_context.Context) context.Context {
	defer ctx.Wait()
	cfg, ok := ctx.Get("config")
	if !ok {
		log.Fatalln("no config provided")
	}
	config, ok := cfg.(*configuration.ConfigStruct)
	if !ok {
		log.Fatalln("invalid config type")
	}

	if len(config.Include) < 1 {
		log.Fatalln("add at least one file to include list")
	}

	gte := githandler.GitEngine{}
	gte.Init(ctx.NewChild())

	tick := time.Duration(config.UpdateRate) * time.Minute
	tckr := time.NewTicker(tick)
	client, err_ := netclient.NewNetworkClient(ctx)
	if err_ != nil {
		log.Fatalln(err_)
	}
	tmp := map[string]any{}
	o := ""
	var err error = &ftp_context.LogItem{}
	for ; err != nil; err = nil {
		res, err := MakeGetRequest(client, "https://127.0.0.1:8080/cert", tmp)
		if err != nil {
			log.Println(err)
			<-time.After(time.Second * 10)
			continue
		}
		o = string(res)

	}
	log.Println("\n", tmp, "\n", o)

	for ok {

		child_ctx := ctx.NewChild()
		child_ctx.SetDeadline(tick)
		log.Println("starting git cycle")
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

		for _, directory := range config.Include {
			log.Println("loading ", directory)
			ls, err := filehandler.ReadDir(child_ctx, directory, append(config.Exclude, ".git"))
			if err != nil {
				log.Fatalln(err.Error())
			}
			for _, f := range ls[:5] {
				fmt.Println(f.Path, " found")
			}
			act_err := actions.Write_directory_files_list(directory, ls)
			if act_err != nil {
				log.Fatalln(act_err)
			}

			err = gte.Commit(directory)
			if err != nil {
				log.Fatalln(err)
			}
		}
		// child_ctx.Cancel()
		select {
		case _, ok = <-ctx.Done():

		case <-child_ctx.Done():

		case _, ok = <-tckr.C:
			log.Println("new tick")
		}
	}

	return ctx
}

func MakeGetRequest(client *http.Client, route string, tmp any) (out []byte, err ftp_context.LogErr) {
	loc := "MakeRequest(method, route string, tmp any ) (out []byte, err ftp_context.LogErr)"
	var eror error
	log.Println(loc)
	res, eror := client.Get(route)
	if eror != nil {
		log.Println(err)
		return nil, ftp_context.NewLogItem(loc, true).
			SetAfter("client.Get").
			AppendParentError(eror)

	}
	out = make([]byte, res.ContentLength+1)
	res.Body.Read(out)
	eror = json.Unmarshal(out, tmp)
	if eror != nil {
		log.Println(err)
		return nil, ftp_context.NewLogItem(loc, true).
			SetAfter("client.Get").
			AppendParentError(eror)

	}

	return
}
