package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ftp_context "github.com/it-shiloheye/ftp_system/v2/lib/context"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging"
	"github.com/it-shiloheye/ftp_system/v2/lib/logging/log_item"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system/v2/lib/tls_handler/v2"
)

var Logger = logging.Logger

type ServerType struct {
	Port string `json:"port"`

	HttpsServer *http.Server
	*gin.Engine

	ServerRun func(ctx ftp_context.Context, err_c chan error)
}

func (st *ServerType) InitServer(server_cert *ftp_tlshandler.TLSCert, srv_name string) {

	loc := log_item.Locf("func (st *ServerType) InitServer(server_cert *ftp_tlshandler.TLSCert, srv_name: %s) ", srv_name)

	st.Engine = gin.Default()
	r := st.Engine

	server_ip, ip_net, err1 := net.ParseCIDR("192.168.0.0/24")
	if err1 != nil {
		log.Fatalln(&log_item.LogItem{
			Location:  loc,
			After:     `ip, ip_net, err1  := net.ParseCIDR("192.168.0.0/24")`,
			Message:   err1.Error(),
			Level:     log_item.LogLevelError02,
			CallStack: []error{err1},
		})
	}

	r.Use(func(ctx *gin.Context) {
		req_ip := ctx.RemoteIP()

		if valid_ip(req_ip, ip_net, server_ip) {

			ctx.Next()
			return
		}
		ctx.Status(400)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	st.HttpsServer = &http.Server{
		Addr:              st.Port,
		Handler:           r,
		ReadHeaderTimeout: time.Millisecond,
		ReadTimeout:       time.Minute * 5,
	}
	if server_cert != nil {
		st.HttpsServer.TLSConfig = ftp_tlshandler.ServerTLSConf(server_cert.TlsCert)
	}

	st.ServerRun = func(ctx ftp_context.Context, err_c chan error) {
		defer ctx.Finished()
		log.Println("Starting:", srv_name, "thread")

		if st.HttpsServer.TLSConfig != nil {
			if err_ := st.HttpsServer.ListenAndServeTLS("", ""); err_ != nil {
				err_c <- Logger.LogErr(loc, &log_item.LogItem{
					Message:   "server failed",
					Level:     log_item.LogLevelError01,
					CallStack: []error{err_},
				})
			}
		} else {
			if err_ := st.HttpsServer.ListenAndServe(); err_ != nil {
				err_c <- Logger.LogErr(loc, &log_item.LogItem{
					Message:   "server failed",
					Level:     log_item.LogLevelError01,
					CallStack: []error{err_},
				})
			}
		}
	}

}

func valid_ip(ip string, ip_net *net.IPNet, server_ip net.IP) bool {
	req_ip := net.ParseIP(ip)

	if req_ip.Equal(server_ip) || net.IPv6loopback.Equal(req_ip) {
		return true
	}

	if req_ip.IsLoopback() {
		return true
	}

	if ip_net.Contains(req_ip) {
		return true
	}

	return false
}
