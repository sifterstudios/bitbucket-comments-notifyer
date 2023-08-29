package main

import (
	"fmt"

	"github.com/sifterstudios/bitbucket-notifier/web"
)

var secretKey [32]byte

func main() {
	fmt.Println("Welcome!")
	fmt.Println("Looking up config file...")
	initialize()
	fmt.Println("Config file loaded!")
	fmt.Println("Starting Web Server...")

	web.StartWebServer()
	fmt.Println("Hello")
}
