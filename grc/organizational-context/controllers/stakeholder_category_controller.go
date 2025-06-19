package controllers

import (
	"jethur/grc/organizational-context/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StakeholderCategoryController struct {
	service services.StakeholderCategoryServiceInterface
	db      *gorm.DB
}

func NewStakeholderCategoryController(service services.StakeholderCategoryServiceInterface) *StakeholderCategoryController {
	return &StakeholderCategoryController{
		service: service,
	}
}

func (esc *StakeholderCategoryController) GetAllStakeholderCategoryes(c *gin.Context) {
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

func (esc *StakeholderCategoryController) GetStakeholderCategoryDetails(c *gin.Context) {
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
