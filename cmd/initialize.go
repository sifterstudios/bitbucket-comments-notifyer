package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/sifterstudios/bitbucket-notifier/auth"
	"github.com/sifterstudios/bitbucket-notifier/data"
)

func initialize() {
	if data.FileOrFolderExists(data.DataFolder) && data.FileOrFolderExists(data.SecurityFile) { // TODO: This should be data's responsibility
		data.GetSecretKey()
	} else {
		data.GetRandomKey()
		data.CreateAndSaveSecurityFile()
	}
	if data.FileOrFolderExists(data.ConfigFolder) && data.FileOrFolderExists(data.ConfigFile) {
		data.UserConfig = data.GetConfig()
	} else {
		createAndSaveConfigFile()
	}
	data.Logbook = data.GetPersistentData()
}

func createAndSaveConfigFile() {
	fmt.Println("Looks like you're new here! Let's get you set up.")
	fmt.Println("Please enter your Bitbucket username:")
	var username string
	_, err := fmt.Scanln(&username)
	if err != nil {
		return
	}
	fmt.Println("Please enter your Bitbucket password:")
	var password string
	_, err = fmt.Scanln(&password)
	if err != nil {
		return
	}

	fmt.Println("Please enter the full address for the bitbucket server(e.g: https://bitbucket.example.com):")
	var address string
	_, err = fmt.Scanln(&address)
	if err != nil {
		return
	}

	encryptedUsername, encryptedPassword, err := auth.EncryptCredentials(
		[]byte(username),
		[]byte(password),
		&data.SecretKey)
	if err != nil {
		panic(err)
	}

	data.UserConfig.Credentials.Username = encryptedUsername
	data.UserConfig.Credentials.Password = encryptedPassword
	data.UserConfig.Bitbucket.ServerUrl = address
	fmt.Println("Username: " + username)
	fmt.Println("Encrypted Username: " + string(encryptedUsername))

	getDefaultSettings(&data.UserConfig)

	configFile, err := yaml.Marshal(data.UserConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(data.ConfigFile, configFile, 0600)
	if err != nil {
		log.Fatal(err)
	}

	data.UserConfig.Credentials.Username = []byte(username)
	data.UserConfig.Credentials.Password = []byte(password)
}

func getDefaultSettings(config *data.Config) {
	config.ConfigNotifications.PollingInterval = 5
	config.ConfigNotifications.Comments = true
	config.ConfigNotifications.Tasks = true
	config.ConfigNotifications.StatusChanges = true
	config.ConfigNotifications.CompletionTime = true
}
