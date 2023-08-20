package bitbucket

import (
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/sifterstudios/bitbucket-comments-notifyer/data"
)

func getActivePullRequestsByUser(config data.Config) {
	client := resty.New()

	apiUrl := config.Bitbucket.ServerUrl + data.ActivePullRequestsApiPath
	username := string(config.Credentials.Username)
	password := string(config.Credentials.Password)

	client.SetBasicAuth(username, password)

	resp, err := client.R().Get(apiUrl)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() == 200 {
		fmt.Println("Success!")
	} else {
		fmt.Printf("GET request failed with status code: %d\n", resp.StatusCode())
	}

	fmt.Println("Response Body: ", resp.String())
}

func getPullRequestActivity(data.Value) {
	// Original string
	url := "/rest/api/latest/projects/projectname/repos/reponame/pull-requests/PR-id/activities"

	// Define replacements
	replacements := map[string]string{
		"projectname": "newproject",
		"reponame":    "newrepo",
		"PR-id":       "newPRid",
	}

	// Replace each placeholder
	for placeholder, replacement := range replacements {
		url = strings.ReplaceAll(url, placeholder, replacement)
	}

	// Print the updated string
	fmt.Println(url)
}
