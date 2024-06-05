package main

import (
	"log"
	"time"

	initialiseclient "github.com/it-shiloheye/ftp_system/client/init_client"
	mainthread "github.com/it-shiloheye/ftp_system/client/main_thread"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

var ClientConfig = initialiseclient.ClientConfig

func main() {

	log.Println("new client started: ", ClientConfig.Id)
	ctx := ftp_context.CreateNewContext()
	defer ctx.Wait()
	ctx.Add()
	go UpdateConfig(ctx)
	mainthread.MainThread(ctx.Add())

}

func UpdateConfig(ctx ftp_context.Context) {
	defer ctx.Finished()
	tc := time.NewTicker(time.Minute)
	for ok := true; ok; {
		select {
		case <-tc.C:
		case _, ok = <-ctx.Done():
		}

		err := initialiseclient.WriteConfig()
		if err != nil {
			log.Println(err)
		}
		log.Println("updated config successfully")
	}
}
