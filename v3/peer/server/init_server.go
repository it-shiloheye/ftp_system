package server

import ( /*
		"encoding/json"
		"errors"
		"io/fs"
		"log"
		"net"
		"os"
		"time"

		ftp_base "github.com/it-shiloheye/ftp_system/v3/lib/base"
		"github.com/it-shiloheye/ftp_system/v3/lib/logging/log_item"
		ftp_tlshandler "github.com/it-shiloheye/ftp_system/v3/lib/tls_handler/v2" */
	"github.com/it-shiloheye/ftp_system/v3/peer/config"
)

var ServerConfig = server_config.ServerConfig

/*
var C_loc = &CertsLocation{
	CertsDirectory: ServerConfig.CertsDirectory,
	caPem:          &ftp_tlshandler.CAPem{},
	cert_d:         ftp_tlshandler.CertData{},
	tlsCert:        &ftp_tlshandler.TLSCert{},
}



func (c *CertsLocation) Cert() *ftp_tlshandler.TLSCert {
	return c.tlsCert
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

func init() {
	start := time.Now()
	loc := log_item.Loc("server/gin_server/init_server.go init()")
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
			Days: 31,
		},
		IPAddrresses: []net.IP{
			net.IPv4(127, 0, 0, 1),
			net.IPv6loopback,
			ServerConfig.LocalIp().To4(),
		},
	}

	ca_buf, err1 := os.ReadFile(C_loc.CA())
	if err1 != nil {
		if !errors.Is(err1, os.ErrNotExist) {
			log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

				After:   "ca_buf, err1 := os.ReadFile(C_loc.CA())",
				Message: err1.Error(),
				Level:   log_item.LogLevelError02, CallStack: []error{err1},
			})
		}

		err2 := os.MkdirAll(C_loc.CertsDirectory, f_mode)
		if err2 != nil && !errors.Is(err2, os.ErrExist) {
			log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

				After:   "err2 := os.MkdirAll(C_loc.CertsDirectory, f_mode)",
				Message: err2.Error(),
				Level:   log_item.LogLevelError02, CallStack: []error{err2, err1},
			})
		}

		cd_buf, err3 := os.ReadFile(C_loc.CertData())
		if err3 != nil {
			if !errors.Is(err3, os.ErrNotExist) {
				log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

					After:   "cd_buf, err3 := os.ReadFile(C_loc.CertData())",
					Message: err3.Error(),
					Level:   log_item.LogLevelError02, CallStack: []error{err3, err1},
				})
			}

			cd_buf, err4 := json.MarshalIndent(&template_cd, " ", "\t")
			if err4 != nil {
				log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

					After:   `cd_buf, err4 := json.MarshalIndent(&template_cd," ","\t")`,
					Message: err4.Error(),
					Level:   log_item.LogLevelError02, CallStack: []error{err3, err1},
				})
			}
			err5 := os.WriteFile(C_loc.CertData(), cd_buf, f_mode)
			if err5 != nil {
				log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

					After:   `err5 := os.WriteFile(C_loc.CertData(),cd_buf,f_mode)`,
					Message: err5.Error(),
					Level:   log_item.LogLevelError02, CallStack: []error{err3, err1},
				})
			}

			log.Fatalf(`please fill the Organisation and CertificateData in: %s`, C_loc.CertData())
		}

		err4 := json.Unmarshal(cd_buf, &C_loc.cert_d)
		if err4 != nil {
			log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

				After:   `err4 := json.Unmarshal(cd_buf,&C_loc.cert_d)`,
				Message: err4.Error(),
				Level:   log_item.LogLevelError02, CallStack: []error{err1},
			})
		}

		tmp_x509 := ftp_tlshandler.ExampleCACert(C_loc.cert_d)

		tmp, err5 := ftp_tlshandler.GenerateCAPem(tmp_x509)
		if err5 != nil {
			log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

				After:   `tmp, err5 := ftp_tlshandler.GenerateCAPem(tmp_x509)`,
				Message: err5.Error(),
				Level:   log_item.LogLevelError02, CallStack: []error{err1},
			})
		}

		*C_loc.caPem = tmp

		ca_buf_, err6 := json.MarshalIndent(&tmp, " ", "\t")
		if err6 != nil {
			log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

				After:   `ca_buf_, err6 := json.MarshalIndent(&tmp," ","\t")`,
				Message: err6.Error(),
				Level:   log_item.LogLevelError02, CallStack: []error{err1},
			})
		}
		ca_buf = ca_buf_

		err7 := ftp_base.WriteFile(C_loc.CA(), ca_buf)
		if err7 != nil {
			log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),
				Level:     log_item.LogLevelError02,
				After:     `err7 := ftp_base.WriteFile(C_loc.CA(),ca_buf)`,
				Message:   err7.Error(),
				CallStack: []error{err1},
			})
		}
	} else {
		// I expect to have a ca_buf with the caPEM data in bytes
		err2 := json.Unmarshal(ca_buf, C_loc.caPem)
		if err2 != nil {
			log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),
				Level:     log_item.LogLevelError02,
				After:     `err2 := json.Unmarshal(ca_buf, C_loc.caPem)`,
				Message:   err2.Error(),
				CallStack: []error{err2},
			})
		}
	}

	// simple time guard, update cert every 7 days, server restarts every day at least once
	ServerConfig.TLS_Cert_Creation = time.Now()

	// generate new tls each time
	x509_tls_cert := ftp_tlshandler.ExampleTLSCert(template_cd)
	tmp, err3 := ftp_tlshandler.GenerateTLSCert(C_loc.caPem, x509_tls_cert)
	if err3 != nil {
		log.Fatalln(&log_item.LogItem{Location: loc, Time: time.Now(),

			After:   "tmp, err3 := ftp_tlshandler.GenerateTLSCert(*C_loc.caPem,x509_tls_cert)",
			Message: err3.Error(),
			Level:   log_item.LogLevelError02, CallStack: []error{err3},
		})
	}
	*C_loc.tlsCert = tmp

}
*/
