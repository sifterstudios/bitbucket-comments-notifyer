package data

import (
	"crypto/rand"
	"fmt"
	"github.com/sifterstudios/bitbucket-notifier/auth"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var SecretKey [32]byte

func GetPersistentData() []PersistentPullRequest {
	if !FileOrFolderExists(LogbookFile) {
		err := os.WriteFile(LogbookFile, []byte{}, 0600)
		if err != nil {
			println("Error creating logbook file")
		}
	}

	fileData, err := os.ReadFile(LogbookFile)
	if err != nil {
		println("Error reading logbook file")
	}

	var persistentPrs []PersistentPullRequest
	if err := yaml.Unmarshal(fileData, &persistentPrs); err != nil {
		println("Error unmarshalling logbook file")
	}
	return persistentPrs
}

func SavePersistentData() {
	data, err := yaml.Marshal(Logbook)
	if err != nil {
		print("Error marshalling logbook file")
	}
	err = os.WriteFile(LogbookFile, data, 0600)
	if err != nil {
		print("Error writing logbook file")
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
		panic(err)
	}
}

func GetConfig() Config {
	fileData, err := os.ReadFile(ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := yaml.Unmarshal(fileData, &config); err != nil {
		log.Fatal(err)
	}

	decryptedUsername, decryptedPassword, err := auth.DecryptCredentials(&SecretKey, config.Credentials.Username, config.Credentials.Password)

	config.Credentials.Username = decryptedUsername
	config.Credentials.Password = decryptedPassword

	return config
}
func CreateAndSaveSecurityFile() {
	err := os.WriteFile(SecurityFile, SecretKey[:], 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Security file created.")
}
func GetSecretKey() {
	secretData, err := os.ReadFile(SecurityFile)
	if err != nil {
		secretData = nil
	}

	copy(SecretKey[:], secretData)
}
