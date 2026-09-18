package middlewares

import (
	"jethur/grc/models"
	"jethur/utility"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetLicenseToken fetches the license key from the system_preferences table
func GetLicenseToken(db *gorm.DB) (string, error) {
	var licenseKey *string

	err := db.Raw("SELECT license_key FROM system_preferences").Scan(&licenseKey).Error
	if err != nil {
		return "", err
	}

	if licenseKey == nil {
		return "", nil
	}

	return *licenseKey, nil
}

// isSuperAdminRequest returns true if the request carries a valid JWT for a super admin user.
// Used by LicenseMiddleware so super admins can access the app even when the license is expired.
func isSuperAdminRequest(c *gin.Context, db *gorm.DB) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return false
	}
	claims, err := utility.VerifyToken(tokenString)
	if err != nil {
		return false
	}

	var userRoles []models.UserRole
	if err := db.Where("user_id = ?", claims.UserID).Find(&userRoles).Error; err != nil {
		return false
	}
	for _, ur := range userRoles {
		var role models.Role
		if err := db.Where("id = ?", ur.RoleID).First(&role).Error; err != nil {
			continue
		}
		if role.IsSuperAdmin {
			return true
		}
	}
	return false
}

// LicenseMiddleware checks if the license key is expired or mismatched
func LicenseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define routes to bypass (unauthenticated routes that must work even with an expired license)
		bypassRoutes := map[string]bool{
			"/update-license-key": true,
			"/license-key":        true,
			"/auth/login":         true,
		}

		if _, exists := bypassRoutes[c.FullPath()]; exists {
			c.Next()
			return
		}

		// Super admins bypass the license check so they can renew an expired license
		if isSuperAdminRequest(c, db) {
			c.Next()
			return
		}

		// Get license token
		tokenString, err := GetLicenseToken(db)
		if err != nil || tokenString == "" {
			c.Next()
			return
		}

		// Decrypt token
		claims, err := utility.DecryptLicenseKey(tokenString)
		if err != nil {
			utility.SendUnauthorizedErrResp(c, "License is invalid")
			c.Abort()
			return
		}

		// Check expiry date
		expiryDateStr, exists := claims["expiry_date"].(string)
		if !exists {
			utility.SendUnauthorizedErrResp(c, "License expiry date missing")
			c.Abort()
			return
		}

		expiryDate, err := time.Parse("2006-01-02", expiryDateStr)
		if err != nil {
			utility.SendInternalErrResp(db, c, "Invalid expiry date format in license")
			c.Abort()
			return
		}

		if time.Now().After(expiryDate) {
			utility.SendLicenseExpiryErrResp(c, "Your license has expired. Please renew to continue.")
			c.Abort()
			return
		}

		// Validate API URL
		apiURL, exists := claims["api_url"].(string)
		if !exists || apiURL == "" {
			utility.SendUnauthorizedErrResp(c, "API URL missing in license key")
			c.Abort()
			return
		}

		appURL := utility.GetApiURL()
		if !strings.Contains(appURL, apiURL) {
			c.JSON(403, gin.H{
				"code":             "API_URL_MISMATCH",
				"error":            "API URL mismatch",
				"expected_api_url": apiURL,
				"received_app_url": appURL,
			})
			c.Abort()
			return
		}

		// License valid
		c.Next()
	}
}
