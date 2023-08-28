package bitbucket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/sifterstudios/bitbucket-notifier/data"
)

var client = resty.New()

func GetCurrentPullRequestsByUser(config data.Config) (data.ActivePullRequestsResponse, error) {

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
	r := resp.Body()

	var jsonData data.ActivePullRequestsResponse

	jsonErr := json.Unmarshal(r, &jsonData)
	if jsonErr != nil {
		return data.ActivePullRequestsResponse{}, jsonErr
	}

	return jsonData, nil
}

func GetPullRequestsActivity(prs []data.PullRequest) (response [][]data.Activity, err error) {
	username := string(data.UserConfig.Credentials.Username)
	password := string(data.UserConfig.Credentials.Password)
	client.SetBasicAuth(username, password)

	for _, pr := range prs {
		url := getActivityUrl(pr.FromRef.Repository.Project.Key,
			pr.FromRef.Repository.Name,
			pr.ID)
		r, err := client.R().Get(url)
		if err != nil {
			panic(err)
		}

		jsonData := data.PullRequestActivityResponse{}
		err = json.Unmarshal(r.Body(), &jsonData)

		response = append(response, jsonData.Values)
	}

	return response, nil
}

func getActivityUrl(key string, name string, id int) (url string) {
	url = data.UserConfig.Bitbucket.ServerUrl + data.PullRequestActivitiesApiPath

	replacements := map[string]string{
		"projectname": key,
		"reponame":    name,
		"PR-id":       strconv.Itoa(id),
	}

	for placeholder, replacement := range replacements {
		url = strings.ReplaceAll(url, placeholder, replacement)
	}
	return url
}
