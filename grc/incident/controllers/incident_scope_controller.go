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

type IncidentScopeController struct {
	service services.IncidentScopeServiceInterface
	db      *gorm.DB
}

func NewIncidentScopeController(service services.IncidentScopeServiceInterface) *IncidentScopeController {
	return &IncidentScopeController{
		service: service,
	}
}

func (isc *IncidentScopeController) GetAllIncidentScopes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentScopes, err := isc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentScopes)
}

func (isc *IncidentScopeController) GetIncidentScopeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentScope, err := isc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentScope)
}

func (isc *IncidentScopeController) CreateIncidentScope(c *gin.Context) {
	var scope *models.IncidentScope

	if err := c.ShouldBindJSON(&scope); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	createdScope, err := isc.service.Create(currentUserID, scope)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(isc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdScope)
}

func (isc *IncidentScopeController) UpdateIncidentScope(c *gin.Context) {
	var updatedScope *models.IncidentScope

	scopeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedScope); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	updatedScope, err = isc.service.Update(uint(scopeID), currentUserID, updatedScope)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(isc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedScope)
}

func (isc *IncidentScopeController) DeleteIncidentScopes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	failedIDs, err := isc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (isc *IncidentScopeController) ActivateIncidentScopes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	err = isc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (isc *IncidentScopeController) DeactivateIncidentScopes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	err = isc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(isc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
