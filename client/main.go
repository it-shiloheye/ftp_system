package main

import (
	"log"

	mainthread "github.com/ftp_system_client/main_thread"

	configuration "github.com/ftp_system_client/main_thread/config"
	ftp_context "github.com/ftp_system_client/main_thread/context"
)

var cfg = configuration.Config

func main() {
	log.Println("new", cfg.Identity, "started: ", cfg.Id)
	ctx := ftp_context.CreateNewContext()
	ctx.Set("config", cfg)
	mainthread.MainThread(ctx.Add())

}
