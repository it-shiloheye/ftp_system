package netclient

import (
	
	"io"
	"log"
	"net/http"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

func make_post_request(client *http.Client, route string, contentType string, body io.Reader) (res *http.Response, out []byte, err ftp_context.LogErr) {
	loc := "make_get_request(client *http.Client, route string, tmp any) (res *http.Response, out []byte, err ftp_context.LogErr)"
	var eror error

	res, eror = client.Post(route, contentType, body)
	if eror != nil {
		log.Println(eror.Error())
		return res, nil, ftp_context.NewLogItem(loc, true).
			SetAfter("client.Get").
			AppendParentError(eror)

	}
	// log.Println(loc, "client.Get(route)", "done", res)
	out, eror = io.ReadAll(res.Body)
	if eror != nil {
		log.Println(eror.Error())
		return res, nil, ftp_context.NewLogItem(loc, true).
			SetAfter("out, eror = io.ReadAll(res.Body)").
			SetMessage(eror.Error()).
			AppendParentError(eror)
	}

	return
}
