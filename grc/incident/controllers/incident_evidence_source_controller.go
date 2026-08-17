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

type IncidentEvidenceSourceController struct {
	service services.IncidentEvidenceSourceServiceInterface
	db      *gorm.DB
}

func NewIncidentEvidenceSourceController(service services.IncidentEvidenceSourceServiceInterface) *IncidentEvidenceSourceController {
	return &IncidentEvidenceSourceController{
		service: service,
	}
}

func (iesc *IncidentEvidenceSourceController) GetAllIncidentEvidenceSources(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentEvidenceSources, err := iesc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentEvidenceSources)
}

func (iesc *IncidentEvidenceSourceController) GetIncidentEvidenceSourceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentEvidenceSource, err := iesc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentEvidenceSource)
}

func (iesc *IncidentEvidenceSourceController) CreateIncidentEvidenceSource(c *gin.Context) {
	var evidenceSource *models.IncidentEvidenceSource

	if err := c.ShouldBindJSON(&evidenceSource); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	createdEvidenceSource, err := iesc.service.Create(currentUserID, evidenceSource)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iesc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdEvidenceSource)
}

func (iesc *IncidentEvidenceSourceController) UpdateIncidentEvidenceSource(c *gin.Context) {
	var updatedEvidenceSource *models.IncidentEvidenceSource

	evidenceSourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedEvidenceSource); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	updatedEvidenceSource, err = iesc.service.Update(uint(evidenceSourceID), currentUserID, updatedEvidenceSource)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iesc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedEvidenceSource)
}

func (iesc *IncidentEvidenceSourceController) DeleteIncidentEvidenceSources(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	failedIDs, err := iesc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iesc *IncidentEvidenceSourceController) ActivateIncidentEvidenceSources(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	err = iesc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (iesc *IncidentEvidenceSourceController) DeactivateIncidentEvidenceSources(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	err = iesc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iesc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
