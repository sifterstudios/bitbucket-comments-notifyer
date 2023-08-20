package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	data "github.com/sifterstudios/bitbucket-comments-notifyer/data"
	auth "github.com/sifterstudios/bitbucket-comments-notifyer/internal/auth"
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
	fmt.Scanln(&username)
	fmt.Println("Please enter your Bitbucket password:")
	var password string
	fmt.Scanln(&password)

	fmt.Println("Please enter the address for the bitbucket server:")
	var address string
	fmt.Scanln(&address)

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

func createAndSaveSecurityFile() {
	err := os.WriteFile(data.SecurityFile, secretKey[:], 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Security file created.")
}

func getSecretKey() {
	data, err := os.ReadFile(data.SecurityFile)
	if err != nil {
		data = nil
	}

	copy(secretKey[:], data)
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
