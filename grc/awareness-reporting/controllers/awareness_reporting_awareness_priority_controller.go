package controllers

import (
	"jethur/grc/awareness-reporting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessPriorityController struct {
	service services.AwarenessPriorityServiceInterface
	db      *gorm.DB
}

func NewAwarenessPriorityController(service services.AwarenessPriorityServiceInterface) *AwarenessPriorityController {
	return &AwarenessPriorityController{
		service: service,
	}
}

func (apc *AwarenessPriorityController) GetAllAwarenessPriorities(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	awarenessPriorities, err := apc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, awarenessPriorities)
}

func (apc *AwarenessPriorityController) GetAwarenessPriorityDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	awarenessPriority, err := apc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, awarenessPriority)
}