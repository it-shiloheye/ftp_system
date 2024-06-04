package initialiseclient

type ClientConfigStruct struct {
	Schema      string `json:"$schema"`
	SchemaId    string `json:"$id"`
	LocalIp     string `json:"local_ip"`
	WebIp       string `json:"web_ip"`
	CA_Location string `json:"ca_location"`
	ClientId    string `json:"client_id"`
	CommonName  string `json:"common_name"`
	DirConfig   `json:"all_dirs"`

	Directories []DirConfig `json:"directories"`

	ServerIps []string `json:"server_ips"`
	LogDir    string   `json:"log_dir"`

	StopOnError bool `json:"stop_on_error"`
}

type DirConfig struct {
	Id            string   `json:"id"`
	Path          string   `json:"path"`
	ExcludeDirs   []string `json:"exclude_dir"`
	ExcluedFile   []string `json:"exclude_file"`
	ExcludeRegex  []string `json:"exclude_regex"`
	FollowSymlink bool     `json:"follow_symlink"`
	IncludeDir    []string `json:"include_dir"`
	IncludeExt    []string `json:"include_ext"`
	IncludeFile   []string `json:"include_file"`
}
