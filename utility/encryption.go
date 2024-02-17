package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"jethur/config"
	"log"
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
