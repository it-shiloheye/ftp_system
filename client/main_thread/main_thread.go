package mainthread

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ftp_system_client/main_thread/actions"
	configuration "github.com/ftp_system_client/main_thread/config"
	ftp_context "github.com/ftp_system_client/main_thread/context"
	filehandler "github.com/ftp_system_client/main_thread/file_handler"
	githandler "github.com/ftp_system_client/main_thread/git_handler"
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
