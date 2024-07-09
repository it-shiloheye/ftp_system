package networkpeer

import (
	ftp_context "github.com/it-shiloheye/ftp_system/v3/lib/context"
	"github.com/it-shiloheye/ftp_system/v3/lib/logging"

	server_config "github.com/it-shiloheye/ftp_system/v3/peer/config"

	"github.com/it-shiloheye/ftp_system/v3/peer/network-peer/api"
	"github.com/it-shiloheye/ftp_system/v3/peer/network-peer/server"
)

var ServerConfig = server_config.ServerConfig
var Logger = logging.Logger

func CreatePeerServer(ctx ftp_context.Context) *server.ServerType {
	// loc := log_item.Locf(`CreatePeerServer(ctx ftp_context.Context) *server.ServerType`)

	Srvr := &server.ServerType{
		Port: ServerConfig.PeerPort,
	}

	Srvr.InitServer(ServerConfig.TLS_Cert, "peer")

	return Srvr

}

func CreateBrowserServer(ctx ftp_context.Context) <-chan error {
	Srvr := &server.ServerType{
		Port: ServerConfig.BrowserPort,
	}

	Srvr.InitServer(nil, "browser client")

	api.RegisterApi(ctx, Srvr.Engine)
	err_c := make(chan error, 1)
	go Srvr.ServerRun(ctx.Add(), err_c)

	return err_c
}
