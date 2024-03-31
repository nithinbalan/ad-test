package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"jethur/config"
	"log"
)

func EncryptData(data []byte) (string, error) {
	key := getKeyFromEnv()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	// Convert the ciphertext to a hexadecimal string for easy storage/transmission.
	return hex.EncodeToString(ciphertext), nil
}

func DecryptData(encryptedData string) ([]byte, error) {
	key := getKeyFromEnv()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext, err := hex.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func getKeyFromEnv() []byte {
	hexKey := config.GetEnv("ENCRYPTION_KEY", "")
	if len(hexKey) != 64 { // 32 bytes in hexadecimal form
		log.Panic("Encryption key must be 64 characters long (32 bytes in hex)")
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Panic("Failed to decode encryption key from hex:", err)
	}

	return key
}
