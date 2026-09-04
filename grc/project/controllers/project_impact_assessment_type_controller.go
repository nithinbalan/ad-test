package controllers

import (
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectImpactAssessmentTypeController struct {
	service services.ProjectImpactAssessmentTypeServiceInterface
	db      *gorm.DB
}

func NewProjectImpactAssessmentTypeController(service services.ProjectImpactAssessmentTypeServiceInterface) *ProjectImpactAssessmentTypeController {
	return &ProjectImpactAssessmentTypeController{
		service: service,
	}
}

func (piatc *ProjectImpactAssessmentTypeController) GetAllProjectImpactAssessmentTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	impactTypes, err := piatc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(piatc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, impactTypes)
}

func (piatc *ProjectImpactAssessmentTypeController) GetProjectImpactAssessmentTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	impactType, err := piatc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, impactType)
}
