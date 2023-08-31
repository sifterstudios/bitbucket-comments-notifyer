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
}

func handleDifference(pr PullRequest, activity Activity) {
	if !containsActivity(activity) {
		handleNotifying(pr, activity)
		updateCurrentPrActivities(pr, activity, 0, 0) // will reset if handled in HandleNotifying
	}
}

func handleNotifying(pr PullRequest, activity Activity) {
	if authorIsYou(activity) { // TODO: Add option to negate this if debugging
		return
	}
	switch activity.Action {
	case "OPENED":
		notification.NotifyAboutOpenedPr() // TODO: To be implemented User, Title, REPO
		break
	case "COMMENTED":
		if prIsClosed(pr) {
			break
		}
		if len(activity.Comment.CommentThread) != 0 {
			lastComment := activity.Comment.CommentThread[len(activity.Comment.CommentThread)-1]
			notification.NotifyAboutNewAnswer(lastComment.Author.DisplayName, lastComment.Text, activity.CommentAnchor.Path, pr.Title)
		} else {
			if activity.Comment.Severity == "BLOCKER" {
				if activity.Comment.State == "RESOLVED" {
					notification.NotifyAboutClosedTask() // TODO: To be implemented
					break
				}
				notification.NotifyAboutNewTask() // TODO: To be implemented
				break
			}

		}
		notification.NotifyAboutNewComment(activity.User.DisplayName, activity.Comment.Text, activity.CommentAnchor.Path, pr.Title)
		break
	case "RESCOPED":
		notification.NotifyAboutNewAmend() // User, Title, REPO. To be implemented
		break
	case "UPDATED":
		notification.NotifyAboutNewCommit() // User, Title, REPO. To be implemented
		break
	case "APPROVED":
		notification.NotifyAboutApprovedPr() // User, Title, REPO. To be implemented
		break
	case "DECLINED":
		notification.NotifyAboutDeclinedPr() // User, Title, REPO. To be implemented
		break
	case "DELETED":
		notification.NotifyAboutDeletedPr() // User, Title, REPO. To be implemented
		break
	case "MERGED":
		notification.NotifyAboutMergedPr() // User, Title, REPO. To be implemented
		break
	case "REOPENED":
		notification.NotifyAboutReopenedPr() // User, Title, REPO. To be implemented
		break
	case "UNAPPROVED":
		notification.NotifyAboutUnapprovedPr() // User, Title, REPO. To be implemented
		break
	//case "REVIEW_COMMENTED":
	//	notification.NotifyAboutReviewCommented() // TODO: To be implemented
	//	break
	//case "REVIEWED_DISCARDED":
	//	notification.NotifyAboutReviewDiscarded() // TODO: To be implemented
	//	break
	//case "REVIEW_FINISHED":
	//	notification.NotifyAboutReviewFinished() // TODO: To be implemented
	//	break
	case "REVIEWED":
		notification.NotifyAboutReviewed() // TODO: To be implemented
		break
	}
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

	return slug == configUsername ||
		email == configUsername
}

func updateCurrentPrActivities(pr PullRequest, newActivity Activity, timeOpened int64, timeClosed int64) {
	idxOfLogbook := getIdxOfLogbook(pr.ID)

	if idxOfLogbook == -1 { // NOTE: PR not found in logbook
		fmt.Println("Info: PR not found in logbook")
		Logbook = append(Logbook, PersistentPullRequest{
			Id:                   pr.ID,
			NotifiedActivityIds:  []int{newActivity.ID},
			TimeOpened:           timeOpened,
			TimeFinished:         timeClosed,
			DurationOpenToFinish: timeClosed - timeOpened,
		})
		return
	}

	Logbook[idxOfLogbook].NotifiedActivityIds = append(Logbook[idxOfLogbook].NotifiedActivityIds, newActivity.ID)
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

func getIdxOfLogbook(prId int) int {
	for i, prStruct := range Logbook {
		if prStruct.Id == prId {
			return i
		}
	}
	return -1
}

func containsActivity(newActivity Activity) bool {
	var foundComment bool
	var foundCommentThread bool

	for _, persistencePrStruct := range Logbook {
		for _, activityId := range persistencePrStruct.NotifiedActivityIds {
			if activityId == newActivity.ID {
				foundComment = true
			}
			if len(newActivity.Comment.CommentThread) != 0 && activityId == newActivity.Comment.CommentThread[0].ID {
				foundCommentThread = true
			}
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
