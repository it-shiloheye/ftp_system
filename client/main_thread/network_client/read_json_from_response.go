package netclient

import (
	
	"encoding/json"

	"io"
	"log"
	
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	
)

func read_json_from_response(r io.ReadCloser, tmp any) (out []byte, err ftp_context.LogErr) {
	loc := "read_json_from_buffer(r io.Reader,tmp any)(out []byte, err ftp_context.LogErr)"
	out, eror := io.ReadAll(r)
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
