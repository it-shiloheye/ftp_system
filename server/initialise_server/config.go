package initialiseserver

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/it-shiloheye/ftp_system_lib/base"
)

var ServerConfig = &ServerConfigStruct{}

type ServerConfigStruct struct {
	Schema           string           `json:"$schema"`
	SchemaId         string           `json:"$id"`
	LocalIp          string           `json:"local_ip"`
	WebIp            string           `json:"web_ip"`
	CA_Location      string           `json:"ca_location"`
	StorageDirectory string           `json:"storage_dir"`
	TmpDirectory     string           `json:"tmp_directory"`
	RemoteRepository string           `json:"remote_git_repo"`
	Clients          []ClientIDStruct `json:"clients"`
}

type ClientIDStruct struct {
	IpAddress  string `json:"ip_address"`
	Id         string `json:"id"`
	CommonName string `json:"common_name"`
	ConfigPath string `json:"config_path"`
}

func (sc ServerConfigStruct) WriteConfig(file_path string) error {

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
func ReadConfig(file_path string) (sc ServerConfigStruct, err error) {

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

	b, err := os.ReadFile("./config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ServerConfig.Schema = "https://json-schema.org/draft/2020-12/schema"
			ServerConfig.Clients = append(ServerConfig.Clients, ClientIDStruct{
				Id: uuid.New().String(),
			})
			tmp, err1 := json.MarshalIndent(ServerConfig, " ", "\t")
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

	err3 := json.Unmarshal(b, ServerConfig)
	if err3 != nil {
		log.Fatalln(err)
	}
}
