package main

import (
	"log"

	ginserver "github.com/ftp_system_server/gin_server"
	ftp_context "github.com/ftp_system_server/main_thread/context"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load() // 👈 load .env file
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.Println("server started")
	ctx := ftp_context.CreateNewContext()
	defer ctx.Wait()
	ginserver.NewServer(ctx)
}
