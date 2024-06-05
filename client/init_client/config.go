package initialiseclient

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"

	"os"

	"github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

var ClientConfig = &ClientConfigStruct{}

func (sc ClientConfigStruct) WriteConfig(file_path string) error {

	tmp, err1 := json.MarshalIndent(&sc, " ", "\t")
	if err1 != nil {
		return err1
	}

	err2 := os.WriteFile(file_path, tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
	if err2 != nil {
		return err2
	}
	return nil
}
func ReadConfig(file_path string) (sc ClientConfigStruct, err error) {

	b, err1 := os.ReadFile(file_path)
	if err1 != nil {
		return sc, err1
	}

	err2 := json.Unmarshal(b, &sc)
	if err2 != nil {
		err = err2
	}

	return
}

func init() {

	*ClientConfig = BlankClientConfigStruct()

	b, err := os.ReadFile("./config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			tmp, err1 := json.MarshalIndent(ClientConfig, " ", "\t")
			if err1 != nil {
				log.Fatalln(err1)
			}
			err2 := os.WriteFile("./config.json", tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
			if err2 != nil {
				log.Fatalln(err2)
			}
			log.Fatalln("fill in config")
			return
		}
		log.Fatalln(err)
	}

	err3 := json.Unmarshal(b, ClientConfig)
	if err3 != nil {
		log.Fatalln(err)
	}
}

func WriteConfig() (err ftp_context.LogErr) {
	loc := "WriteConfig() (err ftp_context.LogErr)"
	tmp, err1 := json.MarshalIndent(ClientConfig, " ", "\t")
	if err1 != nil {
		return &ftp_context.LogItem{
			Location:  loc,
			Err:       true,
			After:     `tmp, err1 := json.MarshalIndent(ClientConfig, " ", "\t")`,
			Message:   err1.Error(),
			CallStack: []error{err1},
		}
	}
	err2 := os.WriteFile("./config.json", tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
	if err2 != nil {
		return &ftp_context.LogItem{
			Location:  loc,
			Err:       true,
			After:     `err2 := os.WriteFile("./config.json", tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))`,
			Message:   err2.Error(),
			CallStack: []error{err2},
		}
	}

	return
}
