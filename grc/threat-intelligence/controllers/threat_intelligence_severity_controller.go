package controllers

import (
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/grc/threat-intelligence/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThreatIntelligenceSeverityController struct {
	service services.ThreatIntelligenceSeverityServiceInterface
	db      *gorm.DB
}

func NewThreatIntelligenceSeverityController(service services.ThreatIntelligenceSeverityServiceInterface) *ThreatIntelligenceSeverityController {
	return &ThreatIntelligenceSeverityController{
		service: service,
	}
}

func (tisc *ThreatIntelligenceSeverityController) GetAllThreatIntelligenceSeverities(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	severities, err := tisc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, severities)
}

func (tisc *ThreatIntelligenceSeverityController) GetThreatIntelligenceSeverityDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	severity, err := tisc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, severity)
}

func (tisc *ThreatIntelligenceSeverityController) CreateThreatIntelligenceSeverity(c *gin.Context) {
	var severity *tiDTOs.ThreatIntelligenceSeverityFormDTO

	if err := c.ShouldBindJSON(&severity); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	id, err := tisc.service.Create(currentUserID, severity)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tisc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (tisc *ThreatIntelligenceSeverityController) UpdateThreatIntelligenceSeverity(c *gin.Context) {
	severityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSeverity *tiDTOs.ThreatIntelligenceSeverityFormDTO
	if err := c.ShouldBindJSON(&updatedSeverity); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	err = tisc.service.Update(uint(severityID), currentUserID, updatedSeverity)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tisc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tisc *ThreatIntelligenceSeverityController) DeleteThreatIntelligenceSeverities(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	failedIDs, err := tisc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tisc *ThreatIntelligenceSeverityController) ActivateThreatIntelligenceSeverity(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	err = tisc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tisc *ThreatIntelligenceSeverityController) DeactivateThreatIntelligenceSeverity(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	err = tisc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tisc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
