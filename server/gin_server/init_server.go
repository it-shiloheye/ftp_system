package ginserver

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_lib/tls_handler"
)

func NewServer(ctx ftp_context.Context) (ftp_err ftp_context.LogErr) {

	cert_d = ftp_tlshandler.NewCA(cert_d.CertData)
	cad := cert_d.NewCAJson()
	_, err_ := cad.ToJSON()
	if err_ != nil {
		log.Fatalln(err_)
	}
	// log.Fatalln(d)
	server_cert := cert_d.NewServerCert(cert_d.CertData)
	select {
	case ftp_err = <-gin_server_http_thread(ctx, cad.PEM):
		break
	case ftp_err = <-gin_server_main_thread(ctx, server_cert, cad.PEM):
		break
	case <-ctx.Done():
	}

	return
}

var cert_d *ftp_tlshandler.CertSetup

func init() {
	cert_d = &ftp_tlshandler.CertSetup{
		CertData: &ftp_tlshandler.CertData{
			Organization:  "Shiloh Eye, Ltd",
			Country:       "KE",
			Province:      "Coast",
			Locality:      "Mombasa",
			StreetAddress: "2nd Floor, SBM Bank Centre, Nyerere Avenue, Mombasa",
			PostalCode:    "80100",
			NotAfter: ftp_tlshandler.NotAfterStruct{
				Days: 7,
			},
			IPAddrresses: []net.IP{
				net.IPv4(127, 0, 0, 1),
				net.IPv6loopback,
			},
		}}

}

func gin_server_http_thread(ctx ftp_context.Context, caPEM string) <-chan ftp_context.LogErr {
	loc := "func gin_server_http_thread(ctx ftp_context.Context, caPEM string) (err ftp_context.LogErr)"

	err_c := make(chan ftp_context.LogErr, 1)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/cert", func(ctx *gin.Context) {

		ctx.JSON(200, gin.H{
			"ca_pem": caPEM,
			"error":  "",
		})
	})
	ctx.Add()
	go func() {
		defer ctx.Finished()
		log.Println("Starting: gin_server_http_thread")
		if err_ := r.Run(":3000"); err_ != nil {
			err_c <- ftp_context.NewLogItem(loc, true).
				SetAfter(`r.Run()`).
				SetMessage("server failed").
				AppendParentError(err_)
		}
		close(err_c)
	}()
	return err_c
}

func gin_server_main_thread(ctx ftp_context.Context, server_cert *ftp_tlshandler.CertSetup, caPEM string) <-chan ftp_context.LogErr {
	loc := "gin_server_main_thread(ctx ftp_context.Context) (err ftp_context.LogErr)"

	err_c := make(chan ftp_context.LogErr, 1)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/cert", func(ctx *gin.Context) {

		ctx.JSON(200, gin.H{
			"ca_pem": caPEM,
			"error":  "",
		})
	})

	ctx.Add()
	go func() {
		defer ctx.Finished()
		log.Println("Starting: gin_server_main_thread")
		if err_ := ftp_tlshandler.GinHandler(r, server_cert, ":8080"); err_ != nil {
			err_c <- ftp_context.NewLogItem(loc, true).
				SetAfter(`ftp_tlshandler.GinHandler(r, server_cert, "8081")`).
				SetMessage("server failed").
				AppendParentError(err_)
		}
		close(err_c)
	}()

	return err_c
}
