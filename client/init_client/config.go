package initialiseclient

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"time"

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
	log.Println("loading config")
	*ClientConfig = BlankClientConfigStruct()

	b, err1 := os.ReadFile("./config.json")
	if err1 != nil {
		if errors.Is(err1, os.ErrNotExist) {

			tmp, err2 := json.MarshalIndent(ClientConfig, " ", "\t")
			if err2 != nil {
				log.Fatalln(err2.Error())
			}
			err3 := os.WriteFile("./config.json", tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
			if err3 != nil {
				log.Fatalln(err3.Error())
			}
			log.Fatalln("fill in config")
			return
		}
		log.Fatalln(err1.Error())
	}

	err3 := json.Unmarshal(b, ClientConfig)
	if err3 != nil {
		log.Fatalln(err3.Error())
	}

	log.Println("successfull loaded config")
}

func WriteConfig() (err ftp_context.LogErr) {
	loc := "WriteConfig() (err ftp_context.LogErr)"
	tmp, err1 := json.MarshalIndent(ClientConfig, " ", "\t")
	if err1 != nil {
		return &ftp_context.LogItem{Location: loc, Time: time.Now(),
			Err:       true,
			After:     `tmp, err1 := json.MarshalIndent(ClientConfig, " ", "\t")`,
			Message:   err1.Error(),
			CallStack: []error{err1},
		}
	}
	err2 := os.WriteFile("./config.json", tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
	if err2 != nil {
		return &ftp_context.LogItem{Location: loc, Time: time.Now(),
			Err:       true,
			After:     `err2 := os.WriteFile("./config.json", tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))`,
			Message:   err2.Error(),
			CallStack: []error{err2},
		}
	}

	return
}
