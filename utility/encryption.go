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
	"time"
)

func Encrypt(plaintext []byte) (string, error) {
	key := getKeyFromEnv()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string) ([]byte, error) {
	key := getKeyFromEnv()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return nil, err
	}
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return ciphertextBytes, nil
}

func EncryptUser(userID int, currentTime time.Time) (string, error) {
	key := getKeyFromEnv()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := fmt.Sprintf("%d:%d", userID, currentTime.Unix())
	plaintextBytes := []byte(plaintext)

	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return hex.EncodeToString(ciphertext), nil
}

func DecryptUser(encryptedData string) (int, time.Time, error) {
	key := getKeyFromEnv()
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, time.Time{}, err
	}

	ciphertext, err := hex.DecodeString(encryptedData)
	if err != nil {
		return 0, time.Time{}, err
	}

	if len(ciphertext) < aes.BlockSize {
		return 0, time.Time{}, fmt.Errorf("ciphertext is too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	data := string(ciphertext)
	var userID int
	var timestamp int64
	_, err = fmt.Sscanf(data, "%d:%d", &userID, &timestamp)
	if err != nil {
		return 0, time.Time{}, err
	}

	return userID, time.Unix(timestamp, 0).UTC(), nil
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
