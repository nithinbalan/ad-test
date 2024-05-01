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

type IncidentTypeController struct {
	service services.IncidentTypeServiceInterface
	db      *gorm.DB
}

func NewIncidentTypeController(service services.IncidentTypeServiceInterface) *IncidentTypeController {
	return &IncidentTypeController{
		service: service,
	}
}

func (itc *IncidentTypeController) GetAllIncidentTypes(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentTypes, err := itc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentTypes)
}

func (itc *IncidentTypeController) GetIncidentTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentType, err := itc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentType)
}

func (itc *IncidentTypeController) CreateIncidentType(c *gin.Context) {
	var incidentType *models.IncidentType

	if err := c.ShouldBindJSON(&incidentType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	createdType, err := itc.service.Create(currentUserID, incidentType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(itc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdType)
}

func (itc *IncidentTypeController) UpdateIncidentType(c *gin.Context) {
	var updatedType *models.IncidentType

	incidentTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	updatedType, err = itc.service.Update(uint(incidentTypeID), currentUserID, updatedType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(itc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedType)
}

func (itc *IncidentTypeController) DeleteIncidentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	failedIDs, err := itc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (itc *IncidentTypeController) ActivateIncidentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	err = itc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (itc *IncidentTypeController) DeactivateIncidentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	err = itc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(itc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
