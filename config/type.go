package config

type service struct {
	Name     string
	AddrList []string
	LB       bool `mapstructure:"load-balance"`
}

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type smtp struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	FromName string `mapstructure:"from_name"`
}

type redis struct {
	Addr     string
	Username string
	Password string
}

type elasticsearch struct {
	Addr string
	Host string
}

type openAI struct {
	ApiKey   string `mapstructure:"api-key"`
	ApiUrl   string `mapstructure:"api-url"`
	ApiModel string `mapstructure:"api-model"`
}
type oss struct {
	Bucket    string
	AccessKey string
	SecretKey string
	Domain    string
	Region    string
}
type config struct {
	MySQL         mySQL
	Redis         redis
	OSS           oss
	Elasticsearch elasticsearch
	Smtp          smtp
	Administrator administrator
	OpenAI        openAI
	Oss           oss
}

type administrator struct {
	Secret string
}
