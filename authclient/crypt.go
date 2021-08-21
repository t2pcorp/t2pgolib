package authclient

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

type encryptAESReturn struct {
	Key    []byte
	IV     []byte
	Output []byte
}

func encryptAES(input []byte) (*encryptAESReturn, error) {
	key, _ := randomHex(16)
	iv, _ := randomHex(8)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	input = PKCS5Padding(input, 16)

	output := make([]byte, len(input))
	mode.CryptBlocks(output, input)

	output = []byte(base64.StdEncoding.EncodeToString(output))

	key = []byte(base64.StdEncoding.EncodeToString(key))
	iv = []byte(base64.StdEncoding.EncodeToString(iv))

	return &encryptAESReturn{
		Key:    key,
		IV:     iv,
		Output: output,
	}, nil
}

func decryptAES(key, iv, input []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(input, input)
	input = PKCS5Unpadding(input)

	return input, nil
}

// PKCS5Padding is a function for padding string
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5Unpadding is a function for unpadding string
func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func randomHex(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(bytes)), nil
}

func getPublicKeyFromString(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
		// panic("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse DER encoded public key: " + err.Error())
		// panic("failed to parse DER encoded public key: " + err.Error())
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("unknow type of public key")
	}
}

func getPrivateKeyFromString(priStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(priStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("private key is not rsa")
	}

	pv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pv, nil
}

func getDecryptMessageInfo(encMessage string) (encryptedKey, encryptedMsg, error) {
	e := strings.Split(encMessage, ":")
	if len(e) < 2 {
		return nil, nil, errors.New("encrypted message is invalid format")
	}

	c, err := base64.StdEncoding.DecodeString(e[0])
	if err != nil {
		return nil, nil, errors.New((err.Error() + ":0:" + e[0]))
	}

	i, err := base64.StdEncoding.DecodeString(e[1])
	if err != nil {
		return nil, nil, errors.New((err.Error() + ":0:" + e[0]))
	}

	return c, i, nil
}

func encryptAESKey(publicKey string, key, iv []byte) ([]byte, error) {
	pubKey, err := getPublicKeyFromString(publicKey)
	if err != nil {
		return nil, err
	}
	message := []byte(fmt.Sprintf("%s:%s", key, iv))
	return encryptRSA(pubKey, message)
}

func encryptRSA(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func decryptAESKey(privateKey string, keyIvEnc []byte) ([]byte, []byte, error) {
	pv, err := getPrivateKeyFromString(privateKey)
	if err != nil {
		return nil, nil, err
	}

	aesKeyIv, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, pv, keyIvEnc, nil)
	if err != nil {
		return nil, nil, err
	}

	aesKeyIvs := strings.Split(string(aesKeyIv), ":")
	if len(aesKeyIvs) < 2 {
		return nil, nil, errors.New("encrypted message is invalid format")
	}

	key, err := base64.StdEncoding.DecodeString(aesKeyIvs[0])
	if err != nil {
		return nil, nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(aesKeyIvs[1])
	if err != nil {
		return nil, nil, err
	}

	return key, iv, nil
}

type (
	encryptedMsg []byte
	encryptedKey []byte
)

func Encrypt(plainText string, keyContent string) (string, error) {
	c, err := extractKey(keyContent)

	if err != nil {
		return "", nil
	}

	encText, err := encryptMessage(plainText, c.publicKey)
	if err != nil {
		return "", err
	}

	return encText, nil
}

func Decrypt(encryptedText string, keyContent string) (string, error) {
	c, err := extractKey(keyContent)

	if err != nil {
		return "", nil
	}

	encText, err := decryptMessage(c, encryptedText)
	if err != nil {
		return "", err
	}

	return encText, nil
}
