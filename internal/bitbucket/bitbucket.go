package bitbucket

import (
	"github.com/go-resty/resty/v2"
)

type BitbucketClient struct {
	client *resty.Client
}

func NewBitbucketClient(baseUrl, username, password string) *BitbucketClient {
	// TODO: Initialize and configure the client
	return nil
}

func (bb *BitbucketClient) FetchEvents() error {
	// TODO: Implement event fetching
	return nil
}
