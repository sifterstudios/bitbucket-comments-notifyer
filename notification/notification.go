package notification

import (
	"fmt"
	"os/exec"
)

func SendNotification(headline, message string) error {
	cmd := exec.Command("notify-send", headline, message, "-t", "0")

	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("Notification sent: ", headline, message)

	return nil
}
func NotifyAboutOpenedPr() {
	fmt.Println("New PR opened!")
}
func NotifyAboutNewAnswer(authorName string, message string, filePath, prTitle string) {
	fmt.Printf("New answer by %s on PR %s: %s\n", authorName, filePath, message)
	err := SendNotification(fmt.Sprintf("New answer by %s on PR %s", authorName, prTitle), fmt.Sprintf("%s/n %s", filePath, message))
	if err != nil {
		fmt.Println(err)
	}
}
func NotifyAboutNewComment(authorName string, message string, filePath, prTitle string) {
	fmt.Printf("New comment by %s on PR %s: %s\n", authorName, filePath, message)
	err := SendNotification(fmt.Sprintf("New comment by %s on PR %s", authorName, prTitle), fmt.Sprintf("%s/n %s", filePath, message))
	if err != nil {
		fmt.Println(err)
	}
}

func NotifyAboutNewAmend() {
	fmt.Println("New amend!")
}

func NotifyAboutNewCommit() {
	fmt.Println("New commit!")
}

func NotifyAboutApprovedPr() {
	fmt.Println("PR approved!")
}

func NotifyAboutDeclinedPr() {
	fmt.Println("PR declined!")
}

func NotifyAboutDeletedPr() {
	fmt.Println("PR deleted!")
}

func NotifyAboutMergedPr() {
	fmt.Println("PR merged!")
}

func NotifyAboutReopenedPr() {
	fmt.Println("PR reopened!")
}

func NotifyAboutUnapprovedPr() {
	fmt.Println("PR unapproved!")
}

func NotifyAboutReviewCommented() {
	fmt.Println("Review commented!")
}

func NotifyAboutReviewDiscarded() {
	fmt.Println("Review discarded!")
}

func NotifyAboutReviewFinished() {
	fmt.Println("Review finished!")
}

func NotifyAboutReviewed() {
	fmt.Println("Reviewed!")
}

func NotifyAboutNewTask() {
	fmt.Println("New task!")
}
