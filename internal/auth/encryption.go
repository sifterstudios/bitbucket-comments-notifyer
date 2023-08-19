package auth

import "crypto/rand"

func GenerateKey() (*[32]byte, error) {
	var key [32]byte
	_, err := rand.Read(key[:])
	if err != nil {
		return nil, err
	}
	return &key, nil
}
