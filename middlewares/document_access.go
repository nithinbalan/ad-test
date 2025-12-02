package middlewares

import (
	"jethur/grc/models"
	"jethur/utility"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Cache for organization level info

var (
	primaryOrgLevelCache struct {
		level uint8
		field string
		mu    sync.RWMutex
	}

	orgLevelCacheInitialized bool
)

func DocumentAccessCheckingMiddleware(db *gorm.DB) gin.HandlerFunc {

	// Initialize cache on first use

	initOrgLevelCache(db)
	return func(c *gin.Context) {

		documentID := c.Param("id")
		currentUserID, err := utility.GetCurrentUserID(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check super admin status first (most common case)
		isSuperAdmin, err := checkSuperAdminStatus(db, currentUserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user permissions"})
			c.Abort()
			return
		}

		if isSuperAdmin {
			c.Next()
			return
		}

		// Check document access
		hasAccess, err := checkDocumentAccess(db, documentID, currentUserID)

		if err != nil {
			log.Printf("checkDocumentAccess failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check document permissions"})
			c.Abort()
			return
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions to this document"})
			c.Abort()
			return
		}

		c.Next()
	}

}

func initOrgLevelCache(db *gorm.DB) {

	if orgLevelCacheInitialized {
		return
	}

	primaryOrgLevelCache.mu.Lock()

	defer primaryOrgLevelCache.mu.Unlock()

	if orgLevelCacheInitialized {
		return
	}

	var primaryLevel models.OrganizationLevel

	if err := db.Where("is_primary = ?", true).First(&primaryLevel).Error; err != nil {

		// Set default if query fails
		primaryOrgLevelCache.level = 1
		primaryOrgLevelCache.field = "organization_unit_1_id"

	} else {
		primaryOrgLevelCache.level = primaryLevel.Level
		primaryOrgLevelCache.field = getOrgUnitField(primaryLevel.Level)
	}
	orgLevelCacheInitialized = true

}

func getOrgUnitField(level uint8) string {

	switch level {
	case 1:
		return "organization_unit_1_id"
	case 2:
		return "organization_unit_2_id"
	case 3:
		return "organization_unit_3_id"
	case 4:
		return "organization_unit_4_id"
	case 5:
		return "organization_unit_5_id"
	default:
		return "organization_unit_1_id" // fallback
	}

}

func checkSuperAdminStatus(db *gorm.DB, userID uint) (bool, error) {

	var count int64
	err := db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.is_super_admin = ?", userID, true).
		Count(&count).Error
	return count > 0, err

}

func checkDocumentAccess(db *gorm.DB, documentID string, userID uint) (bool, error) {
	primaryOrgLevelCache.mu.RLock()
	orgUnitField := primaryOrgLevelCache.field
	primaryOrgLevelCache.mu.RUnlock()
	// SQL Server compatible query
	query := `
		SELECT CASE 
			WHEN EXISTS (
				SELECT 1 FROM documents d
				WHERE d.id = ?
				AND (
					-- In-review documents: creator, responsible, or shared for review
					(d.document_status_id = (SELECT id FROM document_statuses WHERE type = 'in-review')
					 AND (d.created_by_id = ? OR d.responsible_user_id = ? 
						  OR EXISTS (SELECT 1 FROM document_shared_users dsu 
						             WHERE dsu.document_id = d.id AND dsu.user_id = ? AND dsu.type = 'review-user')))
					OR
					-- Approved documents: creator, responsible, or shared for review
					(d.document_status_id = (SELECT id FROM document_statuses WHERE type = 'approved')
					 AND (d.created_by_id = ? OR d.responsible_user_id = ? 
						  OR EXISTS (SELECT 1 FROM document_shared_users dsu 
						             WHERE dsu.document_id = d.id AND dsu.user_id = ? AND dsu.type = 'review-user')))
					OR				 
					-- Rejected documents: creator, responsible, or shared for review
					(d.document_status_id = (SELECT id FROM document_statuses WHERE type = 'rejected')
					 AND (d.created_by_id = ? OR d.responsible_user_id = ? 
						  OR EXISTS (SELECT 1 FROM document_shared_users dsu 
						             WHERE dsu.document_id = d.id AND dsu.user_id = ? AND dsu.type = 'review-user')))
					OR
					-- Published documents: public, shared, or org unit access
					(d.document_status_id = (SELECT id FROM document_statuses WHERE type = 'published')
					 AND (d.document_access_type_id = (SELECT id FROM document_access_types WHERE type = 'public')
						  OR EXISTS (SELECT 1 FROM document_shared_users dsu 
						             WHERE dsu.document_id = d.id AND dsu.user_id = ? AND dsu.type IN ('shared-user','review-user'))
						  OR EXISTS (SELECT 1 
						             FROM document_shared_units dsu
						             JOIN organization_units ou ON ou.id = dsu.organization_unit_id
						             JOIN users u ON u.id = ?
						             WHERE dsu.document_id = d.id AND ou.id = u.` + orgUnitField + `)
						  OR d.created_by_id = ?
						  OR d.responsible_user_id = ?))
					OR
					-- Other documents: only creator or responsible
					(d.document_status_id NOT IN (SELECT id FROM document_statuses WHERE type IN ('in-review','published','approved'))
					 AND (d.created_by_id = ? OR d.responsible_user_id = ?))
				)
			) 
			THEN 1 ELSE 0 
		END AS hasAccess
	`

	var accessInt int
	err := db.Raw(query, documentID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID).Scan(&accessInt).Error

	if err != nil {
		return false, err
	}

	return accessInt == 1, nil
}
