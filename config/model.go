package config

type Configuration struct {
    Credentials    RouterCredentials `yaml:"credentials"`
    DatabasePath   string            `yaml:"database_path"`
    UpdateInterval int               `yaml:"update_interval"`
    AskForPassword bool              `yaml:"ask_for_password"`
    cfgPath        string
}

type RouterCredentials struct {
    Endpoint string `yaml:"endpoint"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}
