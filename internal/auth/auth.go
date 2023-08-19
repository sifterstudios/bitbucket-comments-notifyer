package auth

import (
	"golang.org/x/crypto/nacl/secretbox"
)

func Encrypt(data []byte, secretKey *[32]byte) ([]byte, error) {
	// TODO: Implement encryption here
	return nil, nil
}

func Decrypt(data []byte, secretKey *[32]byte) ([]byte, error) {
	// TODO: Implement decryption here
	return nil, nil
}
