package main

import (
	"fmt"

	data "github.com/sifterstudios/bitbucket-comments-notifyer/data"
	"github.com/sifterstudios/bitbucket-comments-notifyer/web"
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
	fmt.Println("Starting Web Server...")

	web.StartWebServer()
}
