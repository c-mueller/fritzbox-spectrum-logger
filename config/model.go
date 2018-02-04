package config

const defaultBindAddress = ":8080"
const defaultInterval = 15
const defaultAskForPassword = false
const defaultAutoLaunch = false
const defaultDbPath = "spectra.db"
const defaultEndpoint = "192.168.178.1"

type Configuration struct {
	Credentials    RouterCredentials `yaml:"credentials" json:"credentials"`
	DatabasePath   string            `yaml:"database_path" json:"database_path"`
	UpdateInterval int               `yaml:"update_interval" json:"update_interval"`
	AskForPassword bool              `yaml:"ask_for_password" json:"ask_for_password"`
	Autolaunch     bool              `yaml:"autolaunch" json:"autolaunch"`
	BindAddress    string            `yaml:"bind_address" json:"bind_address"`
	cfgPath        string
}

type RouterCredentials struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}
