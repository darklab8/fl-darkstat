package web

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// https://www.codingexplorations.com/blog/understanding-encryption-in-go-a-developers-guide
// Trivial ssymetric encryption with AES

func encrypt(plaintext []byte, keyString string) (string, error) {
	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize] // Initialization vector
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Return the encoded hex string
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(ciphertext string, keyString string) ([]byte, error) {
	ciphertextBytes, _ := hex.DecodeString(ciphertext)

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		return nil, err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return ciphertextBytes, nil
}
