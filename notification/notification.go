package notification

import (
	"fmt"
	"os/exec"
)

func SendNotification(headline, message string) {
	cmd := exec.Command("notify-send", headline, message, "-t", "0") // TODO: Add timing of notification as an option in front-end

	err := cmd.Run()
	if err != nil {
		print("Error: Not able to run command for sending notification:\n" + err.Error())
	}
}

func NotifyAboutOpenedPr(repo string, user string, prTitle string, prDesc string) {
	SendNotification(fmt.Sprintf(`New PR by %s on %s`, user, repo),
		fmt.Sprintf(`%s:
%s`, prTitle, prDesc))
}

func NotifyAboutComment(authorName string, message string, filePath, prTitle string) {
	filePath = parseFilePath(filePath)
	SendNotification(fmt.Sprintf(`New comment by %s on PR %s`, authorName, prTitle),
		fmt.Sprintf(`%s: 
%s`, filePath, message))
}

func NotifyAboutNewTask(authorName string, message string, filePath, prTitle string) {
	filePath = parseFilePath(filePath)
	SendNotification(fmt.Sprintf(`New TASK by %s on PR %s`, authorName, prTitle),
		fmt.Sprintf(`%s: 
%s`, filePath, message))
}

func NotifyAboutClosedTask(authorName string, message string, filePath, prTitle string) {
	filePath = parseFilePath(filePath)
	SendNotification(fmt.Sprintf(`Task CLOSED by %s on PR %s`, authorName, prTitle),
		fmt.Sprintf(`%s: 
%s`, filePath, message))
}

func NotifyAboutNewAmend(repo string, user string, prTitle string, commit string) {
	SendNotification(fmt.Sprintf(`New amend by %s in %s`, user, repo),
		fmt.Sprintf(`PR: %s:
Amend: %s`, prTitle, commit))
}

func NotifyAboutNewCommit(repo string, user string, prTitle string, commit string) {
	SendNotification(fmt.Sprintf(`New commit by %s in %s`, user, repo),
		fmt.Sprintf(`PR: %s:
Commit: %s`, prTitle, commit))
}

func NotifyAboutApprovedPr(repo string, user string, prTitle string) {
	SendNotification(fmt.Sprintf(`PR APPROVED by %s in %s`, user, repo),
		fmt.Sprintf(`PR: %s`, prTitle))
}

func NotifyAboutDeclinedPr(repo string, user string, prTitle string) {
	SendNotification(fmt.Sprintf(`PR declined by %s in %s`, user, repo),
		fmt.Sprintf(`PR: %s`, prTitle))
}

func NotifyAboutMergedPr(repo string, user string, prTitle string) {
	SendNotification(fmt.Sprintf(`PR merged by %s in %s`, user, repo),
		fmt.Sprintf(`PR: %s`, prTitle))
}

func NotifyAboutReviewed(repo string, user string, prTitle string) {
	SendNotification(fmt.Sprintf(`PR set to NEEDS WORK by %s in %s`, user, repo),
		fmt.Sprintf(`PR: %s`, prTitle))
}

func NotifyAboutUnreviewedPr(user string, prTitle string) {
	SendNotification(fmt.Sprintf(`PR by %s has no comments or status`, user),
		fmt.Sprintf(`PR: %s`, prTitle))
}

func parseFilePath(path string) string {
	if path == "" {
		return "In general"
	}
	lastSlashIndex := 0
	for i, c := range path {
		if i == len(path)-1 {
			break
		}
		if c == '/' {
			lastSlashIndex = i
		}
	}
	return path[lastSlashIndex+1:]
}
