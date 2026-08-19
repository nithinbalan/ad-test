package controllers

import (
	"jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentSLAEscalationMatrixController struct {
	service services.IncidentSLAEscalationMatrixServiceInterface
	db      *gorm.DB
}

func NewIncidentSLAEscalationMatrixController(service services.IncidentSLAEscalationMatrixServiceInterface) *IncidentSLAEscalationMatrixController {
	return &IncidentSLAEscalationMatrixController{
		service: service,
	}
}

func (isemc *IncidentSLAEscalationMatrixController) GetAllEscalationMatrices(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	matrices, err := isemc.service.GetAllEscalationMatrices(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(isemc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, matrices)
}

func (isemc *IncidentSLAEscalationMatrixController) GetAllResolutionMatrices(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	matrices, err := isemc.service.GetAllResolutionMatrices(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(isemc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, matrices)
}

func (isemc *IncidentSLAEscalationMatrixController) UpdateEscalationMatrix(c *gin.Context) {
	var req *dtos.IncidentEscalationMatrixUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(isemc.db, c, err.Error())
		return
	}

	updatedMatrices, err := isemc.service.UpdateEscalationMatrix(currentUserID, req)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(isemc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedMatrices)
}

func (isemc *IncidentSLAEscalationMatrixController) UpdateResolutionMatrix(c *gin.Context) {
	var req *dtos.IncidentResolutionMatrixUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(isemc.db, c, err.Error())
		return
	}

	updatedMatrices, err := isemc.service.UpdateResolutionMatrix(currentUserID, req)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(isemc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedMatrices)
}
