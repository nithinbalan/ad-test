package middlewares

import (
	"jethur/grc/models"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RolePermissionMiddleware(db *gorm.DB, moduleIdentifier string, requiredPermissionType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the current user ID from context (assuming a JWT or session-based auth is being used)
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Fetch the user's roles
		var userRoles []models.UserRole
		if err := db.Where("user_id = ?", currentUserID).Find(&userRoles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user roles"})
			c.Abort()
			return
		}

		// Extract role IDs and check if the user is a super admin
		roleIDs := make([]uint, len(userRoles))
		isSuperAdmin := false
		for i, userRole := range userRoles {
			roleIDs[i] = userRole.RoleID

			// Check if any of the roles is a super admin role
			var role models.Role
			if err := db.Where("id = ?", userRole.RoleID).First(&role).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching role details"})
				c.Abort()
				return
			}
			if role.IsSuperAdmin {
				isSuperAdmin = true
				break
			}
		}

		// If the user is a super admin, bypass the permission check
		if isSuperAdmin {
			c.Next()
			return
		}

		// Fetch permissions from role_permissions table for the user's roles and module
		var rolePermissions []models.RolePermission
		if err := db.
			Preload("RolePermissionType").
			Joins("JOIN modules ON modules.id = role_permissions.module_id").
			Joins("JOIN role_permission_types ON role_permission_types.id = role_permissions.role_permission_type_id").
			Where("role_permissions.role_id IN (?) AND modules.identifier = ? AND role_permissions.is_granted = ?", roleIDs, moduleIdentifier, true).
			Find(&rolePermissions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching role permissions"})
			c.Abort()
			return
		}

		// Combine permissions (checking if user has requiredPermissionType by using the Type field)
		hasPermission := false
		for _, permission := range rolePermissions {
			if permission.RolePermissionType.Type == requiredPermissionType {
				hasPermission = true
				break
			}
		}

		// If the user does not have the required permission, deny access
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		// If the user has the required permission, proceed to the next handler
		c.Next()
	}
}
