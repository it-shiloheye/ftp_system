package mainthread

import (
	"context"
	"fmt"
	"log"
	"time"

	initialiseclient "github.com/it-shiloheye/ftp_system/client/init_client"
	"github.com/it-shiloheye/ftp_system/client/main_thread/actions"
	netclient "github.com/it-shiloheye/ftp_system/client/main_thread/network_client"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler"
	githandler "github.com/it-shiloheye/ftp_system_lib/git_handler"
)

var ClientConfig = initialiseclient.ClientConfig

func MainThread(ctx ftp_context.Context) context.Context {
	defer ctx.Wait()

	if len(ClientConfig.IncludeDir) < 1 {
		log.Fatalln("add at least one file to include list")
	}

	gte := githandler.GitEngine{}
	gte.Init(ctx.NewChild())

	tick := time.Duration(ClientConfig.UpdateRate) * time.Minute
	tckr := time.NewTicker(tick)
	client, err_ := netclient.NewNetworkClient(ctx)
	if err_ != nil {
		log.Fatalln(err_)
	}
	tmp := map[string]any{}

	o, err1 := netclient.MakeGetRequest(client, netclient.Route{
		BaseUrl:  "https://127.0.0.1:8080",
		Pathname: "/ping",
	}, &tmp)
	if err1 != nil {
		log.Fatalln(err1.Error())
	}

	log.Println("\n", tmp, "\n", o)

	for ok := true; ok; {

		log.Println("in loop")
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

		for _, directory := range ClientConfig.IncludeDir {
			log.Println("loading ", directory)
			ls, err := filehandler.ReadDir(child_ctx, directory, append(append(ClientConfig.ExcludeDirs, ".git"), ClientConfig.ExcluedFile...))
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
