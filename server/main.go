package main

import (
	"log"

	ginserver "github.com/it-shiloheye/ftp_system/server/gin_server"
	initialiseserver "github.com/it-shiloheye/ftp_system/server/initialise_server"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

var ServerConfig = initialiseserver.ServerConfig

func main() {

	log.Println("server started")
	ctx := ftp_context.CreateNewContext()
	defer ctx.Wait()
	ginserver.NewServer(ctx)
}
