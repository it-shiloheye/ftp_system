package netclient

import (
	"bytes"
	"encoding/json"

	"io"
	"log"
	"net/http"
	"os"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_lib/tls_handler/v2"
)

// 1MB * 10
var BufferStore *BufferStoreStruct = func() (b *BufferStoreStruct) {
	b = &BufferStoreStruct{}

	return
}()

type BufferStoreStruct struct {
	buffers chan *bytes.Buffer
}

func (bss *BufferStoreStruct) Read(r io.ReadCloser) (out []byte, err ftp_context.LogErr) {
	loc := "(bss *BufferStoreStruct) Read(r io.ReadCloser) (out []byte, err ftp_context.LogErr)"
	if bss.buffers == nil {
		bss.buffers = make(chan *bytes.Buffer, 11)
		for i := 0; i < 10; i += 1 {
			bss.buffers <- bytes.NewBuffer(make([]byte, 1024))
		}
	}

	b := <-bss.buffers
	b.Reset()
	n, err_ := b.ReadFrom(r)
	if err_ != nil {
		err = ftp_context.NewLogItem(loc, true).SetAfter("_,err_ := io.ReadAll(r)").AppendParentError(err_)
		return
	}

	r.Close()
	out = (b.Bytes()[:n])
	go func() { bss.buffers <- b }()
	return
}

func NewNetworkClient(ctx ftp_context.Context) (cl *http.Client, err ftp_context.LogErr) {
	loc := "NewNetworkClient(ctx ftp_context.Context)(cl *http.Client, err ftp_context.LogErr )"
	cl = &http.Client{}
	tmp, err1 := os.ReadFile("./data/certs/ca_certs.json")
	if err1 != nil {
		err = ftp_context.NewLogItem(loc, true).SetAfterf("tmp, err1 := os.ReadFile(%s)", "./certs/ca_certs.json").SetMessage(err1.Error()).AppendParentError(err1)
		return
	}

	ca := ftp_tlshandler.CAPem{}
	err2 := json.Unmarshal(tmp, &ca)
	if err2 != nil {
		err = ftp_context.NewLogItem(loc, true).SetAfterf("err2 := json.Unmarshal(tmp,&ca)").SetMessage(err2.Error()).AppendParentError(err2)
		return
	}

	client_tls_config := ftp_tlshandler.ClientTLSConf(ca)
	cl.Transport = &http.Transport{
		TLSClientConfig: client_tls_config,
	}

	return

}

func make_get_request(client *http.Client, route string, tmp any) (res *http.Response, out []byte, err ftp_context.LogErr) {
	loc := "make_get_request(client *http.Client, route string, tmp any) (res *http.Response, out []byte, err ftp_context.LogErr)"
	var eror error

	res, eror = client.Get(route)
	if eror != nil {
		log.Println(eror.Error())
		return res, nil, ftp_context.NewLogItem(loc, true).
			SetAfter("client.Get").
			AppendParentError(eror)

	}
	log.Println(loc, "client.Get(route)", "done", res)
	out, eror = BufferStore.Read(res.Body)
	if eror != nil {
		log.Println(eror.Error())
		return res, nil, ftp_context.NewLogItem(loc, true).
			SetAfter("out, eror = BufferStore.Read(res.Body)").
			SetMessage(eror.Error()).
			AppendParentError(eror)
	}
	log.Println(loc, out)
	eror = json.Unmarshal(out, tmp)
	if eror != nil {
		log.Println(eror.Error())
		return res, nil, ftp_context.NewLogItem(loc, true).
			SetAfter("json.Unmarshal(out, tmp)").
			AppendParentError(eror)

	}

	return
}

func read_json_from_response(r io.ReadCloser, tmp any) (out []byte, err ftp_context.LogErr) {
	loc := "read_json_from_buffer(r io.Reader,tmp any)(out []byte, err ftp_context.LogErr)"
	out, eror := BufferStore.Read(r)
	if eror != nil {
		log.Println(eror.Error())
		err = ftp_context.NewLogItem(loc, true).
			SetAfter("out, eror = BufferStore.Read(res.Body)").
			SetMessage(eror.Error()).
			AppendParentError(eror)
		return
	}
	log.Println(loc, out)
	eror_ := json.Unmarshal(out, tmp)
	if eror_ != nil {
		log.Println(eror.Error())
		err = ftp_context.NewLogItem(loc, true).
			SetAfter("json.Unmarshal(out, tmp)").
			SetMessage(eror.Error()).
			AppendParentError(eror)
		return
	}

	return
}
