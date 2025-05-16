package controllers

import (
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyAssessmentStatusController struct {
	service services.ThirdPartyAssessmentStatusServiceInterface
	db      *gorm.DB
}

func NewThirdPartyAssessmentStatusController(service services.ThirdPartyAssessmentStatusServiceInterface) *ThirdPartyAssessmentStatusController {
	return &ThirdPartyAssessmentStatusController{
		service: service,
	}
}

func (lsc *ThirdPartyAssessmentStatusController) GetAllThirdPartyAssessmentStatuses(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	thirdPartyAssessmentStatuses, err := lsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(lsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, thirdPartyAssessmentStatuses)
}

func (lsc *ThirdPartyAssessmentStatusController) GetThirdPartyAssessmentStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	thirdPartyAssessmentStatus, err := lsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, thirdPartyAssessmentStatus)
}
