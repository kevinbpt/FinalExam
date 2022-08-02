package model

type Auth struct {
	Email    string
	Password string
}

type Token struct {
	Token string
}

type ConfigData struct {
	Server   string `yaml:"server"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Secret   string `yaml:"secretkey"`
}
