package main

import (
	"jethur/config"
	"jethur/middleware"
	"jethur/utility"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv() // Load environment variables

	port := config.GetEnv("PORT", "8080")
	rootPath := config.GetEnv("ROOT_PATH", "")
	if rootPath == "" {
		log.Fatal("ROOT_PATH is not set")
	}
	bundle := utility.LoadBundle(rootPath) //loading language bundles

	r := gin.Default()

	//middlewares
	r.Use(middleware.DatabaseMiddleware())
	r.Use(middleware.LocalizerMiddleware(bundle))
	r.Use(middleware.RecoveryWithDatabaseLog())
	r.Use(middleware.ThrottleMiddleware(1, 5))

	r.GET("/compliance", func(c *gin.Context) {
		// Sample data to encrypt
		plaintext := "Sensitive information here"

		// Encrypt the data
		encryptedData, err := utility.Encrypt([]byte(plaintext))
		if err != nil {
			log.Panic("Encryption error:", err)
			return
		}

		// Decrypt the data
		decryptedData, err := utility.Decrypt("sadfasdf")
		if err != nil {
			log.Panic("Decryption error:", err)
			return
		}

		// Respond with the original, encrypted, and decrypted data
		c.JSON(200, gin.H{
			"original":  plaintext,
			"encrypted": encryptedData,
			"decrypted": string(decryptedData),
		})
	})

	// Start the server on specified port
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
