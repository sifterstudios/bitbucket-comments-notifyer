package data

var UserConfig Config

type Config struct {
	Bitbucket struct {
		ServerUrl string `yaml:"server_url"`
	} `yaml:"bitbucket"`
	Notification struct {
		PollingInterval int `yaml:"polling_interval"`
	} `yaml:"notification"`
	Credentials struct {
		Username []byte `yaml:"username"`
		Password []byte `yaml:"password"`
	} `yaml:"credentials"`
}
