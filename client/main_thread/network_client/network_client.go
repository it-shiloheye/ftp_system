package netclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_lib/tls_handler"
)

func NewNetworkClient(ctx ftp_context.Context) *http.Client {

	client := &http.Client{}

	res, err := client.Get("http://localhost:3000/cert")
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBuffer(make([]byte, res.ContentLength+1))

	io.Copy(buf, res.Body)
	res.Body.Close()

	tmp := map[string]string{}

	err = json.Unmarshal([]byte(buf.String()), &tmp)
	if err != nil {
		log.Fatalln(err)
	}

	pem := []byte(tmp["ca_pem"])
	log.Println("\n", string(pem))

	return ftp_tlshandler.TLSClient(bytes.NewBuffer(pem))
}
