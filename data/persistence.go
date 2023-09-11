package data

import (
	"crypto/rand"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/sifterstudios/bitbucket-notifier/auth"
)

var SecretKey [32]byte

func GetPersistentData() []PersistentPullRequest {
	if !FileOrFolderExists(LogbookFile) {
		err := os.WriteFile(LogbookFile, []byte{}, 0600)
		if err != nil {
			println("Error creating logbook file")
			os.Exit(1)
		}
	}

	fileData, err := os.ReadFile(LogbookFile)
	if err != nil {
		println("Error reading logbook file")
		os.Exit(1)
	}

	var persistentPrs []PersistentPullRequest
	if err := yaml.Unmarshal(fileData, &persistentPrs); err != nil {
		println("Error unmarshalling logbook file")
		os.Exit(1)
	}
	return persistentPrs
}

func SavePersistentData() {
	data, err := yaml.Marshal(Logbook)
	if err != nil {
		print("Error marshalling logbook file")
		os.Exit(1)
	}
	err = os.WriteFile(LogbookFile, data, 0600)
	if err != nil {
		print("Error writing logbook file")
		os.Exit(1)
	}
}

func FileOrFolderExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetRandomKey() {
	if _, err := rand.Read(SecretKey[:]); err != nil {
		print("Error: Couldn't generate Random key..")
		os.Exit(1)
	}
}

func GetConfig() Config {
	fileData, err := os.ReadFile(ConfigFile)
	if err != nil {
		print("Error: Couldn't read config file. Do you have sufficient read permissions?\n" + err.Error())
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		print("Error: Couldn't unmarshal config file\n" + err.Error())
	}

	decryptedUsername, decryptedPassword, err := auth.DecryptCredentials(&SecretKey, config.Credentials.Username, config.Credentials.Password)
	if err != nil {
		print("Error: Couldn't decrypt credentials..\n" + err.Error())
		os.Exit(1)
	}

	config.Credentials.Username = decryptedUsername
	config.Credentials.Password = decryptedPassword

	return config
}

func CreateAndSaveSecurityFile() {
	err := os.WriteFile(SecurityFile, SecretKey[:], 0600)
	if err != nil {
		print("Error writing security file")
		os.Exit(1)
	}

	fmt.Println("Security file created.")
}

func GetSecretKey() {
	secretData, err := os.ReadFile(SecurityFile)
	if err != nil {
		print("Error reading security file. Do you have sufficient permissions?")
		os.Exit(1)
	}

	copy(SecretKey[:], secretData)
}
