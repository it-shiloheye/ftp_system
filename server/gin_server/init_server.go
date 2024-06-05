package ginserver

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	initialiseserver "github.com/it-shiloheye/ftp_system/server/initialise_server"

	"github.com/gin-gonic/gin"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_lib/tls_handler/v2"
)

var ServerConfig = initialiseserver.ServerConfig
var certs_loc = CertsLocation{
	CertsDirectory: ServerConfig.CertsDirectory,
	caPem:          &ftp_tlshandler.CAPem{},
	cert_d:         ftp_tlshandler.CertData{},
	tlsCert:        &ftp_tlshandler.TLSCert{},
}

type CertsLocation struct {
	CertsDirectory string
	cert_d         ftp_tlshandler.CertData
	caPem          *ftp_tlshandler.CAPem
	tlsCert        *ftp_tlshandler.TLSCert
}

func (cd CertsLocation) CA() string {
	return cd.CertsDirectory + "/ca_certs.json"
}

func (cd CertsLocation) TLS() string {
	return cd.CertsDirectory + "/tls_certs.json"
}

func (cd CertsLocation) CertData() string {
	return cd.CertsDirectory + "/certs_data.json"
}

func NewServer(ctx ftp_context.Context) (ftp_err ftp_context.LogErr) {

	select {
	case ftp_err = <-gin_server_main_thread(ctx, certs_loc.tlsCert):
		break
	case <-ctx.Done():
	}

	return
}

func init() {
	start := time.Now()
	loc := "server/gin_server/init_server.go init()"
	f_mode := fs.FileMode(ftp_base.S_IRWXU | ftp_base.S_IRWXO)
	defer func() {
		log.Printf(`server initialised certs, took: %03dms`, time.Since(start).Milliseconds())
	}()

	template_cd := ftp_tlshandler.CertData{
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

	ca_buf, err1 := os.ReadFile(certs_loc.CA())
	if err1 != nil {
		if !errors.Is(err1, os.ErrNotExist) {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     "ca_buf, err1 := os.ReadFile(certs_loc.CA())",
				Message:   err1.Error(),
				CallStack: []error{err1},
			})
		}

		err2 := os.MkdirAll(certs_loc.CertsDirectory, f_mode)
		if err2 != nil && !errors.Is(err2, os.ErrExist) {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     "err2 := os.MkdirAll(certs_loc.CertsDirectory, f_mode)",
				Message:   err2.Error(),
				CallStack: []error{err2, err1},
			})
		}

		cd_buf, err3 := os.ReadFile(certs_loc.CertData())
		if err3 != nil {
			if !errors.Is(err3, os.ErrNotExist) {
				log.Fatalln(&ftp_context.LogItem{
					Location:  loc,
					Err:       true,
					After:     "cd_buf, err3 := os.ReadFile(certs_loc.CertData())",
					Message:   err3.Error(),
					CallStack: []error{err3, err1},
				})
			}

			cd_buf, err4 := json.MarshalIndent(&template_cd, " ", "\t")
			if err4 != nil {
				log.Fatalln(&ftp_context.LogItem{
					Location:  loc,
					Err:       true,
					After:     `cd_buf, err4 := json.MarshalIndent(&template_cd," ","\t")`,
					Message:   err4.Error(),
					CallStack: []error{err3, err1},
				})
			}
			err5 := os.WriteFile(certs_loc.CertData(), cd_buf, f_mode)
			if err5 != nil {
				log.Fatalln(&ftp_context.LogItem{
					Location:  loc,
					Err:       true,
					After:     `err5 := os.WriteFile(certs_loc.CertData(),cd_buf,f_mode)`,
					Message:   err5.Error(),
					CallStack: []error{err3, err1},
				})
			}

			log.Fatalf(`please fill the Organisation and CertificateData in: %s`, certs_loc.CertData())
		}

		err4 := json.Unmarshal(cd_buf, &certs_loc.cert_d)
		if err4 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `err4 := json.Unmarshal(cd_buf,&certs_loc.cert_d)`,
				Message:   err4.Error(),
				CallStack: []error{err1},
			})
		}

		tmp_x509 := ftp_tlshandler.ExampleCACert(certs_loc.cert_d)

		tmp, err5 := ftp_tlshandler.GenerateCAPem(tmp_x509)
		if err5 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `tmp, err5 := ftp_tlshandler.GenerateCAPem(tmp_x509)`,
				Message:   err5.Error(),
				CallStack: []error{err1},
			})
		}

		*certs_loc.caPem = tmp

		ca_buf_, err6 := json.MarshalIndent(&tmp, " ", "\t")
		if err6 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `ca_buf_, err6 := json.MarshalIndent(&tmp," ","\t")`,
				Message:   err6.Error(),
				CallStack: []error{err1},
			})
		}
		ca_buf = ca_buf_

		err7 := ftp_base.WriteFile(certs_loc.CA(), ca_buf)
		if err7 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `err7 := ftp_base.WriteFile(certs_loc.CA(),ca_buf)`,
				Message:   err7.Error(),
				CallStack: []error{err1},
			})
		}
	} else {
		// I expect to have a ca_buf with the caPEM data in bytes
		err2 := json.Unmarshal(ca_buf, certs_loc.caPem)
		if err2 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `err2 := json.Unmarshal(ca_buf, certs_loc.caPem)`,
				Message:   err2.Error(),
				CallStack: []error{err2},
			})
		}
	}

	// simple time guard, update cert every 7 days, server restarts every day at least once
	latest_creation := ServerConfig.TLS_Cert_Creation
	current_datetime := time.Now()

	time_diff := current_datetime.Sub(latest_creation)
	if time_diff.Hours() > (time.Hour.Hours() * 6 * 24) {
		os.Remove(certs_loc.TLS())
	}

	tls_buf, err2 := os.ReadFile(certs_loc.TLS())
	if err2 != nil {
		if !errors.Is(err2, os.ErrNotExist) {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     "tls_buf, err2 := os.ReadFile(certs_loc.TLS()))",
				Message:   err2.Error(),
				CallStack: []error{err2},
			})
		}
		x509_tls_cert := ftp_tlshandler.ExampleTLSCert(template_cd)
		tmp, err3 := ftp_tlshandler.GenerateTLSCert(*certs_loc.caPem, x509_tls_cert)
		if err3 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     "tmp, err3 := ftp_tlshandler.GenerateTLSCert(*certs_loc.caPem,x509_tls_cert)",
				Message:   err3.Error(),
				CallStack: []error{err3, err2},
			})
		}
		*certs_loc.tlsCert = tmp

		tls_buf_, err4 := json.MarshalIndent(certs_loc.tlsCert, " ", "\t")
		if err4 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `tls_buf_, err4 := json.MarshalIndent(certs_loc.tlsCert," ","\t")`,
				Message:   err4.Error(),
				CallStack: []error{err4, err2},
			})
		}

		err5 := os.WriteFile(certs_loc.TLS(), tls_buf_, f_mode)
		if err5 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `err5 := os.WriteFile(certs_loc.TLS(),tls_buf_,f_mode)`,
				Message:   err5.Error(),
				CallStack: []error{err5, err2},
			})
		}
	} else {
		// I expect to have a tls_buf with the tls_cert data in bytes
		err3 := json.Unmarshal(tls_buf, certs_loc.tlsCert)
		if err3 != nil {
			log.Fatalln(&ftp_context.LogItem{
				Location:  loc,
				Err:       true,
				After:     `err3 := json.Unmarshal(ca_buf, certs_loc.tlsCert)`,
				Message:   err3.Error(),
				CallStack: []error{err3},
			})
		}
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
