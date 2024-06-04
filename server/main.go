package main

import (
	"log"

	ginserver "github.com/ftp_system_server/gin_server"
	initialiseserver "github.com/ftp_system_server/initialise_server"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

var ServerConfig = initialiseserver.ServerConfig

func main() {

	log.Println("server started")
	ctx := ftp_context.CreateNewContext()
	defer ctx.Wait()
	ginserver.NewServer(ctx)
}
