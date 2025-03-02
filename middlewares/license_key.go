package middlewares

import (
	"jethur/utility"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetLicenseToken fetches the license key from the system_preferences table
func GetLicenseToken(db *gorm.DB) (string, error) {
	var licenseKey string

	err := db.Raw("select license_key from system_preferences").Scan(&licenseKey).Error
	if err != nil {
		return "", err
	}
	return licenseKey, nil
}

// LicenseMiddleware checks if the license key is expired
func LicenseMiddleware(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// Define routes to bypass
		bypassRoutes := map[string]bool{
			"/update-license-key": true,
			"/license-key":        true,
		}

		// Skip middleware for bypassed routes
		if _, exists := bypassRoutes[c.FullPath()]; exists {
			c.Next()
			return
		}

		// Retrieve the stored license key from DB
		tokenString, err := GetLicenseToken(db)
		if err != nil || tokenString == "" {
			c.Next()
			return
		}

		claims, err := utility.DecryptJWTToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "License is invalid"})
			c.Abort()
			return
		}

		// Check license expiry
		expiryDateStr, exists := claims["expiry_date"].(string)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "License expiry date missing"})
			c.Abort()
			return
		}

		// Convert expiry date string to time.Time
		expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid expiry date format"})
			c.Abort()
			return
		}

		today := time.Now().UTC().Truncate(24 * time.Hour)

		if today.Equal(expiryDate) {
			c.JSON(http.StatusForbidden, gin.H{"error": "License key expired"})
			c.Abort()
			return
		}

		// Compare expiry date with current time
		// if time.Now().After(expiryDate) {
		// 	c.JSON(http.StatusForbidden, gin.H{"error": "License key expired"})
		// 	c.Abort()
		// 	return
		// }

		// Validate api URL
		apiURL, exists := claims["api_url"].(string)
		if !exists || apiURL == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API URL missing in license key"})
			c.Abort()
			return
		}

		// Get the current app URL from request
		appURL := utility.GetApiURL() // Extracts host from request (e.g., "example.com:8080")
		// Compare api URL with app URL
		if !strings.Contains(appURL, apiURL) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":            "API URL mismatch",
				"expected_api_url": apiURL,
				"received_app_url": appURL,
			})
			c.Abort()
			return
		}

		// License is valid, proceed with the request
		c.Next()
	}
}
