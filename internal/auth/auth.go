package auth

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/nacl/secretbox"
)

func Encrypt(data []byte, secretKey *[32]byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	encrypted := secretbox.Seal(nonce[:], data, &nonce, secretKey)
	return encrypted, nil
}

func Decrypt(data []byte, secretKey *[32]byte) ([]byte, error) {
	var nonce [24]byte
	copy(nonce[:], data[:24])

	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, secretKey)
	if !ok {
		return nil, errors.New("Unable to decrypt data")
	}
	return decrypted, nil
}

func SaveEncryptedCredentials(username, password string, secretKey *[32]byte) error {
	encryptedUsername, err := Encrypt([]byte(username), secretKey)
	if err != nil {
		return err
	}

	encryptedPassword, err := Encrypt([]byte(password), secretKey)
	if err != nil {
		return err
	}

	return nil
}

func GetDecryptedCredentials(secretKey *[32]byte) (string, string, error) {
	decryptedUsername, err := Decrypt(encryptedUsername, secretKey)
	if err != nil {
		return "", "", err
	}

	decryptedPassword, err := Decrypt(encryptedPassword, secretKey)
	if err != nil {
		return "", "", err
	}

	return string(decryptedUsername), string(decryptedPassword), nil
}
