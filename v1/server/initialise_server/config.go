package initialiseserver

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	initclient "github.com/it-shiloheye/ftp_system_client/init_client"

	"github.com/it-shiloheye/ftp_system_lib/base"
	ftp_base "github.com/it-shiloheye/ftp_system_lib/base"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

var ServerConfig = NewConfigStruct()

type ServerConfigStruct struct {
	Schema         string `json:"$schema"`
	SchemaId       string `json:"$id"`
	ServerId       string `json:"server_id"`
	LocalIp        string `json:"local_ip"`
	WebIp          string `json:"web_ip"`
	CertsDirectory string `json:"certs_dir"`
	initclient.DirConfig
	TmpDirectory      string                               `json:"tmp_directory"`
	RemoteRepository  string                               `json:"remote_git_repo"`
	Clients           ftp_base.MutexedMap[*ClientIDStruct] `json:"clients"`
	TLS_Cert_Creation time.Time
}

type ClientIDStruct struct {
	IpAddress   string `json:"ip_address"`
	ClientId    string `json:"client_id"`
	CommonName  string `json:"common_name"`
	StoragePath string `json:"config_path"`
}

func (sc ServerConfigStruct) WriteConfig(file_path string, i ...int) error {
	lock_file_p := ServerConfig.DirConfig.Path + "config.lock"

	l, err1 := filehandler.Lock(lock_file_p)
	if err1 != nil {
		return err1
	}
	defer l.Unlock()

	tmp, err1 := json.MarshalIndent(&sc, " ", "\t")
	if err1 != nil {
		i_ := 0
		if len(i) > 0 {
			i_ = i[0]
			i_ += 1
		}

		if i_ < 5 {
			<-time.After(time.Second)
			return sc.WriteConfig(file_path, i_)
		}

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
			if err1 := WriteConfigToFile(); err1 != nil {
				log.Fatalln("error writing config: ", err1.Error())
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

func NewConfigStruct() (svfg *ServerConfigStruct) {
	svfg = &ServerConfigStruct{
		Schema:   "https://json-schema.org/draft/2020-12/schema",
		ServerId: uuid.NewString(),
		DirConfig: initclient.DirConfig{
			DirId: uuid.NewString(),

			Path:          "./data/storage",
			ExcludeDirs:   []string{},
			ExcluedFile:   []string{},
			ExcludeRegex:  []string{},
			FollowSymlink: false,
			IncludeDir:    []string{"./data/storage"},
			IncludeExt:    []string{},
			IncludeFile:   []string{},
			UpdateRate:    time.Hour,
			PathSeparator: string(os.PathSeparator),
		},
		Clients: ftp_base.NewMutexedMap[*ClientIDStruct](),
	}

	return
}

func WriteConfigToFile(i ...int) error {
	l, err1 := filehandler.Lock("./config.index.lock")
	if err1 != nil {
		i_ := 0
		if len(i) > 0 {
			i_ = i[0]
		}
		if i_ < 5 {
			<-time.After(time.Second)
			return WriteConfigToFile(i_ + 1)
		}
		return err1
	}
	defer l.Unlock()

	tmp, err2 := json.MarshalIndent(ServerConfig, " ", "\t")
	if err2 != nil {
		return err2
	}
	err3 := os.WriteFile("./config.json", tmp, fs.FileMode(base.S_IRWXU|base.S_IRWXO))
	if err2 != nil {
		return err3
	}

	return nil
}

func (sc *ServerConfigStruct) RegisterClient(client_id string, client_ip string, common_name string) (storage_path string) {
	c, ok := sc.Clients.Get(client_id)
	if ok {
		return c.StoragePath
	}

	fs_mode := fs.FileMode(ftp_base.S_IRWXU | ftp_base.S_IRWXO)
	storage_path = sc.Path + string(os.PathSeparator) + common_name + string(os.PathSeparator) + client_id
	sc.Clients.Set(client_id, &ClientIDStruct{
		ClientId:    client_id,
		IpAddress:   client_ip,
		CommonName:  common_name,
		StoragePath: storage_path,
	})

	os.Mkdir(storage_path, fs_mode)
	return
}

func (sc *ServerConfigStruct) GetStoragePath(client_id string) (storage_path string, exists bool) {
	c, exists := sc.Clients.Get(client_id)
	if !exists {
		return
	}

	storage_path = c.StoragePath
	return
}
