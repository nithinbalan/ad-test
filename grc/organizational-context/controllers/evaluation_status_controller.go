package controllers

import (
	"jethur/grc/organizational-context/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EvaluationStatusController struct {
	service services.EvaluationStatusServiceInterface
	db      *gorm.DB
}

func NewEvaluationStatusController(service services.EvaluationStatusServiceInterface) *EvaluationStatusController {
	return &EvaluationStatusController{
		service: service,
	}
}

func (esc *EvaluationStatusController) GetAllEvaluationStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	requirementStatuses, err := esc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(esc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, requirementStatuses)
}

func (esc *EvaluationStatusController) GetEvaluationStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementStatus, err := esc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, requirementStatus)
}
