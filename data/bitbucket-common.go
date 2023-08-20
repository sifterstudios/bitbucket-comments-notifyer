package data

type BitbucketResponse struct {
	Size       int           `json:"size"`
	Limit      int           `json:"limit"`
	IsLastPage bool          `json:"isLastPage"`
	Values     []PullRequest `json:"values"`
	Start      int           `json:"start"`
}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	Links        Links  `json:"links"`
}

type Links struct {
	Self []Link `json:"self"`
}

type Link struct {
	Href string `json:"href"`
	Name string `json:"name"`
}
