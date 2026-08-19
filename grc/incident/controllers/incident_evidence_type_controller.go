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

type IncidentEvidenceTypeController struct {
	service services.IncidentEvidenceTypeServiceInterface
	db      *gorm.DB
}

func NewIncidentEvidenceTypeController(service services.IncidentEvidenceTypeServiceInterface) *IncidentEvidenceTypeController {
	return &IncidentEvidenceTypeController{
		service: service,
	}
}

func (ietc *IncidentEvidenceTypeController) GetAllIncidentEvidenceTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentEvidenceTypes, err := ietc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentEvidenceTypes)
}

func (ietc *IncidentEvidenceTypeController) GetIncidentEvidenceTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentEvidenceType, err := ietc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentEvidenceType)
}

func (ietc *IncidentEvidenceTypeController) CreateIncidentEvidenceType(c *gin.Context) {
	var evidenceType *models.IncidentEvidenceType

	if err := c.ShouldBindJSON(&evidenceType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	createdEvidenceType, err := ietc.service.Create(currentUserID, evidenceType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ietc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdEvidenceType)
}

func (ietc *IncidentEvidenceTypeController) UpdateIncidentEvidenceType(c *gin.Context) {
	var updatedEvidenceType *models.IncidentEvidenceType

	evidenceTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedEvidenceType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	updatedEvidenceType, err = ietc.service.Update(uint(evidenceTypeID), currentUserID, updatedEvidenceType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ietc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedEvidenceType)
}

func (ietc *IncidentEvidenceTypeController) DeleteIncidentEvidenceTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	failedIDs, err := ietc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ietc *IncidentEvidenceTypeController) ActivateIncidentEvidenceTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	err = ietc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ietc *IncidentEvidenceTypeController) DeactivateIncidentEvidenceTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	err = ietc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ietc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
