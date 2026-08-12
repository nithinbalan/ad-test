package controllers

import (
	"jethur/grc/incident/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentImpactController struct {
	service services.IncidentImpactServiceInterface
	db      *gorm.DB
}

func NewIncidentImpactController(service services.IncidentImpactServiceInterface) *IncidentImpactController {
	return &IncidentImpactController{
		service: service,
	}
}

func (iic *IncidentImpactController) GetAllIncidentImpacts(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentImpacts, err := iic.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentImpacts)
}

func (iic *IncidentImpactController) GetIncidentImpactDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentImpact, err := iic.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentImpact)
}

func (iic *IncidentImpactController) CreateIncidentImpact(c *gin.Context) {
	var impact *models.IncidentImpact

	if err := c.ShouldBindJSON(&impact); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	createdImpact, err := iic.service.Create(currentUserID, impact)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdImpact)
}

func (iic *IncidentImpactController) UpdateIncidentImpact(c *gin.Context) {
	var updatedImpact *models.IncidentImpact

	impactID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedImpact); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	updatedImpact, err = iic.service.Update(uint(impactID), currentUserID, updatedImpact)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedImpact)
}

func (iic *IncidentImpactController) DeleteIncidentImpacts(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	failedIDs, err := iic.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iic *IncidentImpactController) ActivateIncidentImpacts(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	err = iic.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (iic *IncidentImpactController) DeactivateIncidentImpacts(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	err = iic.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iic.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
