package data

import (
	"fmt"
	"github.com/sifterstudios/bitbucket-notifier/notification"
)

var (
	CurrentPrActivity []Activity
)

func HandlePrActivity(activePrs []PullRequest, allSlicesOfActivities [][]Activity) {
	if len(activePrs) != len(allSlicesOfActivities) {
		fmt.Println("Error: Mismatch of PRs and slices of activities returned")
		return
	}
	for i, sliceOfActivities := range allSlicesOfActivities {
		for _, a := range sliceOfActivities {
			handleDifference(activePrs[i].Title, a)
		}
	}
	if len(CurrentPrActivity) == 0 { // TODO: Will this ever happen?
		CurrentPrActivity = flatten(allSlicesOfActivities)
	}
}

func flatten(activities [][]Activity) []Activity {
	var flattened []Activity
	for _, slice := range activities {
		for _, activity := range slice {
			flattened = append(flattened, activity)
		}
	}
	return flattened
}

func handleDifference(prTitle string, activity Activity) {
	if !containsActivity(CurrentPrActivity, activity) { // TODO: I think now every comment will be notified when there's an answer to that comment.
		handleNotifying(prTitle, activity)
		CurrentPrActivity = updateCurrentPrActivities(CurrentPrActivity, activity)
	}
}

func handleNotifying(prTitle string, activity Activity) {
	if authorIsYou(activity) { // TODO: Add option to negate this if debugging
		return
	}
	switch activity.Action {
	case "OPENED":
		notification.NotifyAboutOpenedPr()
		break
	case "COMMENTED":
		if len(activity.Comment.CommentThread) != 0 {
			lastComment := activity.Comment.CommentThread[len(activity.Comment.CommentThread)-1]
			notification.NotifyAboutNewAnswer(lastComment.Author.DisplayName, lastComment.Text, activity.CommentAnchor.Path, prTitle)
		} else {
			notification.NotifyAboutNewComment(activity.User.DisplayName, activity.Comment.Text, activity.CommentAnchor.Path, prTitle)
		}
		break
	case "RESCOPED":
		notification.NotifyAboutNewAmend()
		break
	case "UPDATED":
		notification.NotifyAboutNewCommit()
		break
	case "APPROVED":
		notification.NotifyAboutApprovedPr()
		break
	case "DECLINED":
		notification.NotifyAboutDeclinedPr()
		break
	case "DELETED":
		notification.NotifyAboutDeletedPr()
		break
	case "MERGED":
		notification.NotifyAboutMergedPr()
		break
	case "REOPENED":
		notification.NotifyAboutReopenedPr()
		break
	case "UNAPPROVED":
		notification.NotifyAboutUnapprovedPr()
		break
	case "REVIEW_COMMENTED":
		notification.NotifyAboutReviewCommented()
		break
	case "REVIEWED_DISCARDED":
		notification.NotifyAboutReviewDiscarded()
		break
	case "REVIEW_FINISHED":
		notification.NotifyAboutReviewFinished()
		break
	case "REVIEWED":
		notification.NotifyAboutReviewed()
		break
	}
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

	return slug == configUsername ||
		email == configUsername
}

func updateCurrentPrActivities(currentPrs []Activity, newActivity Activity) []Activity {
	var found bool
	for i, activity := range currentPrs {
		if activity.ID == newActivity.ID {
			currentPrs[i] = newActivity
			found = true
			break
		}
	}
	if !found {
		CurrentPrActivity = append(CurrentPrActivity, newActivity)
	}

	return currentPrs
}

func containsActivity(currentPrActivity []Activity, newActivity Activity) bool {
	var foundComment bool
	var foundCommentThread bool
	for _, activity := range currentPrActivity {
		if activity.ID == newActivity.ID {
			foundComment = true
		}
		if len(newActivity.Comment.CommentThread) > len(activity.Comment.CommentThread) {
			foundCommentThread = true
		}
	}
	return foundComment && foundCommentThread
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
