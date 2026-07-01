package controllers

import (
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfidentialityClassificationController struct {
	service services.ConfidentialityClassificationServiceInterface
	db      *gorm.DB
}

func NewConfidentialityClassificationController(service services.ConfidentialityClassificationServiceInterface) *ConfidentialityClassificationController {
	return &ConfidentialityClassificationController{
		service: service,
	}
}

func (ccc *ConfidentialityClassificationController) GetAllConfidentialityClassifications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	confidentialityClassifications, err := ccc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, confidentialityClassifications)
}

func (ccc *ConfidentialityClassificationController) GetConfidentialityClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	confidentialityClassification, err := ccc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, confidentialityClassification)
}