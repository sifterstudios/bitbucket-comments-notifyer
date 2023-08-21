package data

var UserConfig Config

type Config struct {
	Bitbucket struct {
		ServerUrl string `yaml:"server_url"`
	} `yaml:"bitbucket"`
	Notification struct {
		PollingInterval int  `yaml:"polling_interval"`
		Comments        bool `yaml:"comments"`
		Tasks           bool `yaml:"tasks"`
		StatusChanges   bool `yaml:"status_changes"`
		CompletionTime  bool `yaml:"completion_time"`
	} `yaml:"notification"`
	Credentials struct {
		Username []byte `yaml:"username"`
		Password []byte `yaml:"password"`
	} `yaml:"credentials"`
}
