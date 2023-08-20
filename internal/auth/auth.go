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

func EncryptCredentials(username, password []byte, secretKey *[32]byte) ([]byte, []byte, error) {
	encryptedUsername, err := Encrypt([]byte(username), secretKey)
	if err != nil {
		return nil, nil, err
	}

	encryptedPassword, err := Encrypt([]byte(password), secretKey)
	if err != nil {
		return nil, nil, err
	}

	return encryptedUsername, encryptedPassword, nil
}

func DecryptCredentials(secretKey *[32]byte, encryptedUsername, encryptedPassword []byte) ([]byte, []byte, error) {
	decryptedUsername, err := Decrypt(encryptedUsername, secretKey)
	if err != nil {
		return nil, nil, err
	}

	decryptedPassword, err := Decrypt(encryptedPassword, secretKey)
	if err != nil {
		return nil, nil, err
	}

	return decryptedUsername, decryptedPassword, nil
}
