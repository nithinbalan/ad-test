package middlewares

import (
	"jethur/utility"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	var JWTSecret = []byte("f0a3e3b61e7082367820944efeb9d81c7e242ec4d22dc6fa941aaf76d333f158")

	return func(c *gin.Context) {
		// Retrieve the stored license key from DB
		tokenString, err := GetLicenseToken(db)
		if err != nil || tokenString == "" {
			// c.JSON(http.StatusForbidden, gin.H{"error": "License key not found"})
			// c.Abort()
			c.Next()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil // Use the same secret key from service package
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired license key"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid license claims"})
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

		// Compare expiry date with current time
		if time.Now().After(expiryDate) {
			c.JSON(http.StatusForbidden, gin.H{"error": "License key expired"})
			c.Abort()
			return
		}

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
