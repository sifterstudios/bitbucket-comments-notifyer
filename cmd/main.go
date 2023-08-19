package main

import (
	"crypto/rand"
	"fmt"
	"os"
)

const (
	securityFile = "data/.securi.ty"
)

var secretKey [32]byte

func main() {
	if securityFileExists() {
		getSecretKey()
	} else {
		getRandomKey()
		createAndSaveSecurityFile()
	}
}

func createAndSaveSecurityFile() {
	os.Create(securityFile)
	err := os.WriteFile(securityFile, secretKey[:], 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Security file created.")
}

func getSecretKey() {
	data, err := os.ReadFile(securityFile)
	if err != nil {
		data = []byte("")
	}

	secretKey = [32]byte(data)
}

func securityFileExists() bool {
	if _, err := os.Stat(securityFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func getRandomKey() {
	if _, err := rand.Read(secretKey[:]); err != nil {
		panic(err)
	}
}
