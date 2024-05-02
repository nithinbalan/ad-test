package controllers

import (
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewFrequencyController struct {
	service services.ReviewFrequencyServiceInterface
	db      *gorm.DB
}

func NewReviewFrequencyController(service services.ReviewFrequencyServiceInterface) *ReviewFrequencyController {
	return &ReviewFrequencyController{
		service: service,
	}
}

func (rfc *ReviewFrequencyController) GetAllReviewFrequencies(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	reviewFrequencies, err := rfc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rfc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, reviewFrequencies)
}

func (rfc *ReviewFrequencyController) GetReviewFrequencyDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reviewFrequency, err := rfc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, reviewFrequency)
}
