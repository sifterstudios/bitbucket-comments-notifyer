package data

import (
	"fmt"

	"github.com/sifterstudios/bitbucket-notifier/notification"
)

var Logbook []PersistentPullRequest

func HandlePrActivity(activePrs []PullRequest, allSlicesOfActivities [][]Activity) {
	if len(activePrs) != len(allSlicesOfActivities) {
		fmt.Println("Error: Mismatch of PRs and slices of activities returned")
		return
	}
	for i, sliceOfActivities := range allSlicesOfActivities {
		for _, a := range sliceOfActivities {
			handleDifference(activePrs[i], a)
		}
	}
	SavePersistentData()
}

func handleDifference(pr PullRequest, activity Activity) {
	if !containsActivity(activity.ID) {
		handleNotifying(pr, activity)
		updateCurrentPrActivities(pr, activity, 0, 0)
	}
}

func handleNotifying(pr PullRequest, activity Activity) {
	if authorIsYou(activity) { // TODO: Add option to negate this if debugging
		return
	}
	switch activity.Action {
	case "OPENED":
		notification.NotifyAboutOpenedPr(pr.FromRef.Repository.Name, activity.User.DisplayName, pr.Title, pr.Description)
		updateCurrentPrActivities(pr, activity, activity.CreatedDate, 0)
		break
	case "COMMENTED":
		if prIsClosed(pr) {
			break
		}
		handleCommentLogic(pr, activity)
		break
	case "RESCOPED":
		notification.NotifyAboutNewAmend(pr.FromRef.Repository.Name, activity.User.DisplayName, pr.Title, activity.Diff.Destination.Name)
		break
	case "UPDATED":
		notification.NotifyAboutNewCommit(pr.FromRef.Repository.Name, activity.User.DisplayName, pr.Title, activity.Diff.Destination.Name)
		break
	case "APPROVED":
		notification.NotifyAboutApprovedPr(pr.FromRef.Repository.Name, activity.User.DisplayName, pr.Title)
		break
	case "DECLINED":
		notification.NotifyAboutDeclinedPr(pr.FromRef.Repository.Name, activity.User.DisplayName, pr.Title)
		updateCurrentPrActivities(pr, activity, 0, activity.CreatedDate)
		break
	case "MERGED":
		notification.NotifyAboutMergedPr(pr.FromRef.Repository.Name, activity.User.DisplayName, pr.Title)
		updateCurrentPrActivities(pr, activity, 0, activity.CreatedDate)
		break
	case "REVIEWED":
		notification.NotifyAboutReviewed(pr.FromRef.Repository.Name, activity.User.DisplayName, pr.Title)
		break
	}
}

func handleCommentLogic(pr PullRequest, activity Activity) {
	commentThread := activity.Comment.CommentThread
	if len(commentThread) != 0 {
		for _, answer := range commentThread {
			if containsActivity(answer.ID) || authorIsYou(activity) {
				return
			}
			notification.NotifyAboutComment(answer.Author.DisplayName, answer.Text, activity.CommentAnchor.Path, pr.Title)
		}
	}
	if isTask(activity) {
		if isTaskClosed(activity) {
			notification.NotifyAboutClosedTask(activity.User.DisplayName, activity.Comment.Text, activity.CommentAnchor.Path, pr.Title)
		}
		notification.NotifyAboutNewTask(activity.User.DisplayName, activity.Comment.Text, activity.CommentAnchor.Path, pr.Title)
	}
	notification.NotifyAboutComment(activity.User.DisplayName, activity.Comment.Text, activity.CommentAnchor.Path, pr.Title)
}

func isTaskClosed(activity Activity) bool {
	return activity.Comment.State == "RESOLVED"
}

func isTask(activity Activity) bool {
	return activity.Comment.Severity == "BLOCKER"
}

func prIsClosed(pr PullRequest) bool {
	return pr.State == "DECLINED" || pr.State == "MERGED" || pr.State == "UNAPPROVED" || pr.State == "DELETED"
}

func authorIsYou(activity Activity) bool { // NOTE: Different servers use email/username to authenticate
	configUsername := string(UserConfig.Credentials.Username)
	slug := activity.Comment.Author.Slug
	email := activity.Comment.Author.EmailAddress

	if activity.Comment.Text == "" {
		return false
	}

	if len(activity.Comment.CommentThread) != 0 {
		slug = activity.Comment.CommentThread[len(activity.Comment.CommentThread)-1].Author.Slug
		email = activity.Comment.CommentThread[len(activity.Comment.CommentThread)-1].Author.EmailAddress
	}

	return slug == configUsername || // BUG: This is still letting through a comment in my PR that I made, overview
		email == configUsername
}

func updateCurrentPrActivities(pr PullRequest, newActivity Activity, timeOpened int64, timeClosed int64) {
	idxOfLogbook := getIdxOfLogbook(pr.ID)

	if idxOfLogbook == -1 { // NOTE: PR not found in logbook
		Logbook = append(Logbook, PersistentPullRequest{
			Id:                   pr.ID,
			NotifiedActivityIds:  []int{newActivity.ID},
			TimeOpened:           timeOpened,
			TimeFinished:         timeClosed,
			DurationOpenToFinish: timeClosed - timeOpened,
		})
		return
	}

	if isActivityNew(idxOfLogbook, newActivity.ID) {
		Logbook[idxOfLogbook].NotifiedActivityIds = append(Logbook[idxOfLogbook].NotifiedActivityIds, newActivity.ID)
	}
	if len(newActivity.Comment.CommentThread) != 0 {
		appendAnswers(idxOfLogbook, newActivity.Comment.CommentThread)
	}
	if timeOpened != 0 {
		Logbook[idxOfLogbook].TimeOpened = timeOpened
	}
	if timeClosed != 0 {
		Logbook[idxOfLogbook].TimeFinished = timeClosed
	}
	if timeOpened != 0 && timeClosed != 0 {
		Logbook[idxOfLogbook].DurationOpenToFinish = timeClosed - timeOpened
	}
}

func appendAnswers(idxOfLogbook int, answers []Comment) {
	for _, answer := range answers {
		if isActivityNew(idxOfLogbook, answer.ID) {
			Logbook[idxOfLogbook].NotifiedActivityIds = append(Logbook[idxOfLogbook].NotifiedActivityIds, answer.ID)
		}
	}
}

func isActivityNew(idxOfLogbook int, newId int) bool {
	for _, id := range Logbook[idxOfLogbook].NotifiedActivityIds {
		if id == newId {
			return false
		}
	}
	return true
}

func getIdxOfLogbook(prId int) int {
	for i, prStruct := range Logbook {
		if prStruct.Id == prId {
			return i
		}
	}
	return -1
}

func containsActivity(id int) bool {
	for _, persistencePrStruct := range Logbook {
		for _, activityId := range persistencePrStruct.NotifiedActivityIds {
			if activityId == id {
				return true
			}
		}
	}
	return false
}

type PullRequestActivityResponse struct {
	Size       int        `json:"size"`
	Limit      int        `json:"limit"`
	IsLastPage bool       `json:"isLastPage"`
	Values     []Activity `json:"values"`
	Start      int        `json:"start"`
}
type Comment struct {
	Properties          CommentProperties `json:"properties"`
	ID                  int               `json:"id"`
	Version             int               `json:"version"`
	Text                string            `json:"text"`
	Author              User              `json:"author"`
	CreatedDate         int64             `json:"createdDate"`
	UpdatedDate         int64             `json:"updatedDate"`
	CommentThread       []Comment         `json:"comments"`
	Tasks               []Task            `json:"tasks"`
	Severity            string            `json:"severity"`
	State               string            `json:"state"`
	PermittedOperations struct {
		Editable       bool `json:"editable"`
		Transitionable bool `json:"transitionable"`
		Deletable      bool `json:"deletable"`
	} `json:"permittedOperations"`
	ResolvedDate int64 `json:"resolvedDate"`
	Resolver     User  `json:"resolver"`
}

type CommentProperties struct {
	RepositoryID int `json:"repositoryId"`
}

type Task struct {
	Properties          Properties `json:"properties"`
	ID                  int        `json:"id"`
	Version             int        `json:"version"`
	Text                string     `json:"text"`
	Author              User       `json:"author"`
	CreatedDate         int64      `json:"createdDate"`
	UpdatedDate         int64      `json:"updatedDate"`
	Comments            []Comment  `json:"comments"`
	Tasks               []Task     `json:"tasks"`
	Severity            string     `json:"severity"`
	State               string     `json:"state"`
	PermittedOperations struct {
		Editable       bool `json:"editable"`
		Transitionable bool `json:"transitionable"`
		Deletable      bool `json:"deletable"`
	} `json:"permittedOperations"`
}

type Diff struct {
	Source      interface{} `json:"source"`
	Destination struct {
		Components []string `json:"components"`
		Parent     string   `json:"parent"`
		Name       string   `json:"name"`
		Extension  string   `json:"extension"`
		ToString   string   `json:"toString"`
	} `json:"destination"`
	Hunks []struct {
		Context         string `json:"context"`
		SourceLine      int    `json:"sourceLine"`
		SourceSpan      int    `json:"sourceSpan"`
		DestinationLine int    `json:"destinationLine"`
		DestinationSpan int    `json:"destinationSpan"`
		Segments        []struct {
			Type  string `json:"type"`
			Lines []struct {
				Destination int    `json:"destination"`
				Source      int    `json:"source"`
				Line        string `json:"line"`
				Truncated   bool   `json:"truncated"`
			} `json:"lines"`
			Truncated bool `json:"truncated"`
		} `json:"segments"`
		Truncated bool `json:"truncated"`
	} `json:"hunks"`
	Truncated  bool `json:"truncated"`
	Properties struct {
		ToHash   string `json:"toHash"`
		Current  bool   `json:"current"`
		FromHash string `json:"fromHash"`
	} `json:"properties"`
}

type CommentAnchor struct {
	FromHash string `json:"fromHash"`
	ToHash   string `json:"toHash"`
	Line     int    `json:"line"`
	LineType string `json:"lineType"`
	FileType string `json:"fileType"`
	Path     string `json:"path"`
	DiffType string `json:"diffType"`
	Orphaned bool   `json:"orphaned"`
}

type Activity struct {
	ID            int           `json:"id"`
	CreatedDate   int64         `json:"createdDate"`
	User          User          `json:"user"`
	Action        string        `json:"action"`
	CommentAction string        `json:"commentAction"`
	Comment       Comment       `json:"comment"`
	CommentAnchor CommentAnchor `json:"commentAnchor"`
	Diff          Diff          `json:"diff"`
}
