package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// EncryptAES256GCMHex is Encrypt AES256 GCM
func EncryptAES256GCMHex(data string, key []byte) (string, error) {
	text := []byte(data)

	if len(key) == 0 {
		return "", errors.New("Error Get Config Key")
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	seal := gcm.Seal(nil, nonce, text, nil)

	output := fmt.Sprintf("%x%x", nonce, seal)
	return output, nil
}

// DecryptAES256GCMHex is Encrypt AES256 GCM
func DecryptAES256GCMHex(encrypt string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(encrypt)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
