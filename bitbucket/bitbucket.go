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

const (
	placeholderProjectKey = "projectname"
	placeholderRepoName   = "reponame"
	placeholderPrId       = "PR-id"
)

func GetCurrentPullRequestsByUser(config data.Config) (data.ActivePullRequestsResponse, error) {
	apiUrl := config.Bitbucket.ServerUrl + data.CurrentPullRequestsApiPath
	username := string(config.Credentials.Username)
	password := string(config.Credentials.Password)

	client.SetBasicAuth(username, password)

	resp, err := client.R().Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode() == 200 {
		fmt.Println("Success!")
	} else {
		fmt.Printf("GET request failed with status code: %d\n", resp.StatusCode())
	}
	r := resp.Body()

	var jsonData data.ActivePullRequestsResponse

	err = json.Unmarshal(r, &jsonData)
	if err != nil {
		return data.ActivePullRequestsResponse{}, err
	}

	return jsonData, nil
}

func GetPullRequestsActivity(prs []data.PullRequest, getCount *int) (response [][]data.Activity, err error) {
	username := string(data.UserConfig.Credentials.Username)
	password := string(data.UserConfig.Credentials.Password)
	client.SetBasicAuth(username, password)

	for _, pr := range prs {
		*getCount++
		url := getActivityUrl(pr.FromRef.Repository.Project.Key,
			pr.FromRef.Repository.Name,
			pr.ID)
		r, err := client.R().Get(url)
		if err != nil {
			print("Error: Couldn't get activity for PR " + strconv.Itoa(pr.ID) + "\n" + err.Error())
		}

		jsonData := data.PullRequestActivityResponse{}
		err = json.Unmarshal(r.Body(), &jsonData)
		if err != nil {
			print("Error: Couldn't unmarshal activity for PR " + strconv.Itoa(pr.ID) + "\n" + err.Error())
		}

		response = append(response, jsonData.Values)
	}

	return response, nil
}

func getActivityUrl(key string, name string, id int) (url string) {
	url = data.UserConfig.Bitbucket.ServerUrl + data.PullRequestActivitiesApiPath

	replacements := map[string]string{
		placeholderProjectKey: key,
		placeholderRepoName:   name,
		placeholderPrId:       strconv.Itoa(id),
	}

	for placeholder, replacement := range replacements {
		url = strings.ReplaceAll(url, placeholder, replacement)
	}
	return url
}
