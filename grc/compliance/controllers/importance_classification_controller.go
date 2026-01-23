package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImportanceClassificationController handles importance classifications
type ImportanceClassificationController struct {
	service services.ImportanceClassificationServiceInterface
	db      *gorm.DB
}

// NewImportanceClassificationController initializes the controller
func NewImportanceClassificationController(service services.ImportanceClassificationServiceInterface) *ImportanceClassificationController {
	return &ImportanceClassificationController{
		service: service,
	}
}

// GetAllImportanceClassifications fetches all classifications with optional filters
func (icc *ImportanceClassificationController) GetAllImportanceClassifications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	classifications, err := icc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, classifications)
}

// GetImportanceClassificationDetails fetches a single classification by ID
func (icc *ImportanceClassificationController) GetImportanceClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	classification, err := icc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, classification)
}
