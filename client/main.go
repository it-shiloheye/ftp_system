package main

import (
	"log"

	initialiseclient "github.com/ftp_system_client/init_client"
	mainthread "github.com/ftp_system_client/main_thread"

	configuration "github.com/it-shiloheye/ftp_system_lib/config"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

var cfg = configuration.Config
var ClientConfig = initialiseclient.ClientConfig

func main() {
	log.Println("new", cfg.Identity, "started: ", cfg.Id)
	ctx := ftp_context.CreateNewContext()
	ctx.Set("config", cfg)
	mainthread.MainThread(ctx.Add())

}
