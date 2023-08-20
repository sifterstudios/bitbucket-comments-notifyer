package main

import (
	"fmt"

	data "github.com/sifterstudios/bitbucket-comments-notifyer/data"
)

var (
	secretKey [32]byte
	config    data.Config
)

func main() {
	fmt.Println("Welcome!")
	fmt.Println("Looking up config file...")
	initialize()
	fmt.Println("Config file loaded!")
	fmt.Println(string(config.Credentials.Username))
}
