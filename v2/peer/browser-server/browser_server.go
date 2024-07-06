package browserserver

import (
	server_config "github.com/it-shiloheye/ftp_system/v2/peer/config"
	"github.com/it-shiloheye/ftp_system/v2/peer/server"
)

var ServerConfig = server_config.ServerConfig

func CreateBrowserServer() *server.ServerType {

	Srvr := &server.ServerType{
		Port: ServerConfig.BrowserPort,
	}

	Srvr.InitServer(nil, "browser")

	return Srvr

}
