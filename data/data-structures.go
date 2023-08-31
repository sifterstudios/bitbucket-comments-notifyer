package data

var (
	UserConfig Config
)

type Config struct {
	Bitbucket struct {
		ServerUrl string `yaml:"server_url"`
	} `yaml:"bitbucket"`
	ConfigNotifications `yaml:"notification"`
	Credentials         struct {
		Username []byte `yaml:"username"`
		Password []byte `yaml:"password"`
	} `yaml:"credentials"`
}

type Notification struct {
	Title string
	Body  string
}

type ConfigNotifications struct {
	PollingInterval int  `yaml:"polling_interval"`
	Comments        bool `yaml:"comments"`
	Tasks           bool `yaml:"tasks"`
	StatusChanges   bool `yaml:"status_changes"`
	CompletionTime  bool `yaml:"completion_time"`
}

type PersistentData struct {
	PersistentPullRequests []PersistentPullRequest `yaml:"persistent_pull_requests"`
}

type PersistentPullRequest struct {
	Id                   int   `yaml:"id"`
	NotifiedActivityIds  []int `yaml:"notified_activities"`
	TimeOpened           int64 `yaml:"time_opened"`
	TimeFinished         int64 `yaml:"time_finished"`
	DurationOpenToFinish int64 `yaml:"duration_open_to_finish"`
}
