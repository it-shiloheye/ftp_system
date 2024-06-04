package ginserver

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/ftp_system_server/main_thread/actions"
	"github.com/gin-gonic/gin"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_lib/tls_handler/v2"
)

func NewServer(ctx ftp_context.Context) (ftp_err ftp_context.LogErr) {
	x509_ca_cert := ftp_tlshandler.ExampleCACert(*cert_d)
	ca_pem, err1 := ftp_tlshandler.GenerateCAPem(x509_ca_cert)
	if err1 != nil {
		log.Fatalln(err1.Error())
	}

	err2 := actions.WriteJson("./data", "ca_cert", &ca_pem)
	if err2 != nil {
		log.Fatalln(err2)
	}
	x509_tls_cert := ftp_tlshandler.ExampleTLSCert(*cert_d)
	server_cert, err3 := ftp_tlshandler.GenerateTLSCert(ca_pem, x509_tls_cert)
	if err3 != nil {
		log.Fatalln(err3.Error())
	}
	select {
	case ftp_err = <-gin_server_main_thread(ctx, &server_cert):
		break
	case <-ctx.Done():
	}

	return
}

var cert_d *ftp_tlshandler.CertData

func init() {
	cert_d = &ftp_tlshandler.CertData{
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
	}
}

func gin_server_main_thread(ctx ftp_context.Context, server_cert *ftp_tlshandler.TLSCert) <-chan ftp_context.LogErr {
	loc := "gin_server_main_thread(ctx ftp_context.Context) (err ftp_context.LogErr)"

	err_c := make(chan ftp_context.LogErr, 1)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	ctx.Add()
	go func() {
		defer ctx.Finished()
		log.Println("Starting: gin_server_main_thread")

		srv := http.Server{
			Addr:      ":" + os.Getenv("PORT"),
			Handler:   r,
			TLSConfig: ftp_tlshandler.ServerTLSConf(server_cert.TlsCert),
		}

		log.Println("https://127.0.0.1", srv.Addr)

		if err_ := srv.ListenAndServeTLS("", ""); err_ != nil {
			err_c <- ftp_context.NewLogItem(loc, true).
				SetAfter(`srv.ListenAndServeTLS("","")`).
				SetMessage("server failed").
				AppendParentError(err_)
		}
		close(err_c)
	}()

	return err_c
}
