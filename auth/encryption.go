package auth

import "crypto/rand"

func GenerateKey() (*[]byte, error) {
	var key []byte
	_, err := rand.Read(key[:])
	if err != nil {
		return nil, err
	}
	return &key, nil
}
