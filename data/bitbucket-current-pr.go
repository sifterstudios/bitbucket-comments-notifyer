package data

var CurrentPrs []PullRequest

func HandleCurrentPrs(newPrs []PullRequest) {
	newPrs = filterClosedPrs(newPrs)
	CurrentPrs = newPrs
}

func filterClosedPrs(newPrs []PullRequest) []PullRequest {
	var filteredPrs []PullRequest
	for _, pr := range newPrs {
		if prIsClosedAndNotified(pr) {
			continue
		}
		filteredPrs = append(filteredPrs, pr)
	}
	return filteredPrs
}

func prIsClosedAndNotified(pr PullRequest) bool {
	for _, persistentPr := range Logbook {
		if persistentPr.Id == pr.ID && persistentPr.TimeFinished != 0 {
			return true
		}
	}
	return false
}

type ActivePullRequestsResponse struct {
	Size       int           `json:"size"`
	Limit      int           `json:"limit"`
	IsLastPage bool          `json:"isLastPage"`
	Values     []PullRequest `json:"values"`
	Start      int           `json:"start"`
}
type PullRequest struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	State        string     `json:"state"`
	Open         bool       `json:"open"`
	Closed       bool       `json:"closed"`
	CreatedDate  int64      `json:"createdDate"`
	UpdatedDate  int64      `json:"updatedDate"`
	FromRef      Ref        `json:"fromRef"`
	ToRef        Ref        `json:"toRef"`
	Locked       bool       `json:"locked"`
	Author       User       `json:"author"`
	Reviewers    []User     `json:"reviewers"`
	Participants []User     `json:"participants"`
	Properties   Properties `json:"properties"`
	Links        Links      `json:"links"`
}

type Ref struct {
	ID           string     `json:"id"`
	DisplayID    string     `json:"displayId"`
	LatestCommit string     `json:"latestCommit"`
	Type         string     `json:"type"`
	Repository   Repository `json:"repository"`
}

type Repository struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	HierarchyId   string  `json:"hierarchyId"`
	ScmID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
	Links         Links   `json:"links"`
}

type Project struct {
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       Links  `json:"links"`
}

type Properties struct {
	MergeResult       MergeResult `json:"mergeResult"`
	ResolvedTaskCount int         `json:"resolvedTaskCount"`
	CommentCount      int         `json:"commentCount"`
	OpenTaskCount     int         `json:"openTaskCount"`
}

type MergeResult struct {
	Outcome string `json:"outcome"`
	Current bool   `json:"current"`
}
