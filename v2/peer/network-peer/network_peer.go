package networkpeer

import (
	ftp_context "github.com/it-shiloheye/ftp_system/v2/lib/context"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging"

	server_config "github.com/it-shiloheye/ftp_system/v2/peer/config"

	"github.com/it-shiloheye/ftp_system/v2/peer/server"
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
