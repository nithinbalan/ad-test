package utility

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDBFromContext(c *gin.Context) (*gorm.DB, error) {
	db, exists := c.Get("db")
	if !exists {
		return nil, fmt.Errorf("database connection not found in context")
	}

	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("context contains a non-GORM database connection")
	}

	return gormDB, nil
}
