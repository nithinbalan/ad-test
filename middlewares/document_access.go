package middlewares

import (
	"fmt"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DocumentAccessCheckingMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("id")
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user models.User
		result := db.
			Preload("UserRoles").
			Preload("UserRoles.Role").
			First(&user, currentUserID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user"})
			c.Abort()
			return
		}

		// --- Super Admin Shortcut ---
		isUserSuperAdmin := false
		fmt.Println("isUserSuperAdmin =", result)
		fmt.Println("user =", user)
		fmt.Println("userrole =", user.UserRoles)
		if user.UserRoles != nil {
			for _, userRole := range *user.UserRoles {
				if userRole.Role.IsSuperAdmin {
					isUserSuperAdmin = true
					break
				}
			}
		}
		fmt.Println("isUserSuperAdmin =", isUserSuperAdmin)
		if isUserSuperAdmin {
			fmt.Println("insissUserSuperAdmin =", isUserSuperAdmin)
			c.Next()
			return
		}

		// --- Load Primary Org Level ---
		var primaryLevel models.OrganizationLevel
		if err := db.Where("is_primary = ?", true).First(&primaryLevel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch primary org level"})
			c.Abort()
			return
		}

		var primaryOrgUnitField string
		switch primaryLevel.Level {
		case 1:
			primaryOrgUnitField = "organization_unit_1_id"
		case 2:
			primaryOrgUnitField = "organization_unit_2_id"
		case 3:
			primaryOrgUnitField = "organization_unit_3_id"
		case 4:
			primaryOrgUnitField = "organization_unit_4_id"
		case 5:
			primaryOrgUnitField = "organization_unit_5_id"
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsupported primary level"})
			c.Abort()
			return
		}

		// --- Document Access Query ---
		var documents []models.Document
		query := db.Model(&models.Document{}).
			Where("id = ?", documentID).
			Where(`
            (
                document_status_id = (SELECT id FROM document_statuses WHERE type = ?) 
                AND (created_by_id = ? OR responsible_user_id = ? 
                     OR EXISTS (SELECT 1 FROM document_shared_users dsu WHERE dsu.document_id = documents.id AND dsu.user_id = ? AND dsu.type = ?))
            )
            OR
            (
                document_status_id = (SELECT id FROM document_statuses WHERE type = ?)
                AND (
                    document_access_type_id = (SELECT id FROM document_access_types WHERE type = ?)
                    OR EXISTS (SELECT 1 FROM document_shared_users dsu WHERE dsu.document_id = documents.id AND dsu.user_id = ? AND dsu.type IN (?, ?))
                    OR EXISTS (
                        SELECT 1 
                        FROM document_shared_units dsu
                        JOIN organization_units ou ON ou.id = dsu.organization_unit_id
                        JOIN users u ON u.id = ?
                        WHERE dsu.document_id = documents.id
                        AND ou.id = u.`+primaryOrgUnitField+`
                    )
                )
            )
            OR
            (
                document_status_id NOT IN (SELECT id FROM document_statuses WHERE type IN (?, ?))
                AND (created_by_id = ? OR responsible_user_id = ?)
            )
        `,
				"in-review", currentUserID, currentUserID, currentUserID, "review-user",
				"published", "public", currentUserID, "shared-user", "review-user", currentUserID,
				"in-review", "published", currentUserID, currentUserID,
			)

		if err := query.Find(&documents).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check document permissions"})
			c.Abort()
			return
		}

		if len(documents) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions to this document"})
			c.Abort()
			return
		}

		c.Next()
	}
}
