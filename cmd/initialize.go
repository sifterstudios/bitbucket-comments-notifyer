package main

import (
	"crypto/rand"
	"fmt"
	"github.com/sifterstudios/bitbucket-notifier/auth"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/sifterstudios/bitbucket-notifier/data"
)

func initialize() {
	if fileExists(data.SecurityFile) {
		getSecretKey()
	} else {
		getRandomKey()
		createAndSaveSecurityFile()
	}
	if fileExists(data.ConfigFile) {
		data.UserConfig = getConfig()
	} else {
		createAndSaveConfigFile()
	}
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
		&secretKey)
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

func createAndSaveSecurityFile() {
	err := os.WriteFile(data.SecurityFile, secretKey[:], 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Security file created.")
}

func getSecretKey() {
	secretData, err := os.ReadFile(data.SecurityFile)
	if err != nil {
		secretData = nil
	}

	copy(secretKey[:], secretData)
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func getRandomKey() {
	if _, err := rand.Read(secretKey[:]); err != nil {
		panic(err)
	}
}

func getConfig() data.Config {
	fileData, err := os.ReadFile(data.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	var config data.Config
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		log.Fatal(err)
	}

	decryptedUsername, decryptedPassword, err := auth.DecryptCredentials(&secretKey, config.Credentials.Username, config.Credentials.Password)

	config.Credentials.Username = decryptedUsername
	config.Credentials.Password = decryptedPassword

	return config
}
