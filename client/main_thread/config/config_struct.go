package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/ftp_system_client/base"
)

type Identity string

const (
	IdentityClient Identity = "client"
	IdentityServer Identity = "server"
)

type ConfigStruct struct {
	Version        string   `json:"version"`
	Identity       Identity `json:"identity"`
	Ip             string   `json:"ip"`
	LocalServerIp  string   `json:"local_server_ip"`
	RemoteServerIp string   `json:"remote_server_ip"`
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Include        []string `json:"include"`
	Exclude        []string `json:"exclude"`
	UpdateRate     int      `json:"update_rate"` // number of seconds between updates
	FileTreeFile   string   `json:"file_tree_file"`
	DataDir        string   `json:"data_dir"`
}

func newConfig() *ConfigStruct {

	return &ConfigStruct{
		Version:        "0.0.0",
		Identity:       "client",
		Ip:             "",
		LocalServerIp:  "",
		RemoteServerIp: "",
		Id:             "",
		Name:           "",
		Include:        []string{},
		Exclude:        []string{"node_modules", ".next", "~"},
		UpdateRate:     5,
		FileTreeFile:   "file_tree.json",
		DataDir:        "./data",
	}
}

var Config *ConfigStruct = newConfig()

const client_config = "client.json"

func get_path() string {
	return strings.Join([]string{Config.DataDir, client_config}, "/")
}

func init() {
	err := base.ReadJson(get_path(), &Config)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			err = Config.Write()
			log.Fatalln("please fill out the Config properly ", err)
		}

		log.Fatalln("error during Config init", err)

	}

}

func (cfg *ConfigStruct) Write() (err error) {
	config_text, err := json.MarshalIndent(cfg, " ", "\t")
	if err != nil {
		err = fmt.Errorf("json.MarshalIndent %s", err.Error())
		return
	}

	err = os.MkdirAll(cfg.DataDir, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
	if err != nil {
		err = fmt.Errorf("os.MkdirAll %s", err.Error())
		return
	}
	err = base.WriteFile(get_path(), config_text)
	if err != nil {
		err = fmt.Errorf("base.WriteFile %s", err.Error())
		return
	}

	return
}
