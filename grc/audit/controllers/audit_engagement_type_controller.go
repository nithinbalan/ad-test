package controllers

import (
	"jethur/grc/audit/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditEngagementTypeController struct {
	service services.AuditEngagementTypeServiceInterface
	db      *gorm.DB
}

func NewAuditEngagementTypeController(service services.AuditEngagementTypeServiceInterface) *AuditEngagementTypeController {
	return &AuditEngagementTypeController{
		service: service,
	}
}

func (aetc *AuditEngagementTypeController) GetAllAuditEngagementTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditEngagementTypes, err := aetc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditEngagementTypes)
}

func (aetc *AuditEngagementTypeController) GetAuditEngagementTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditEngagementType, err := aetc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditEngagementType)
}

func (aetc *AuditEngagementTypeController) CreateAuditEngagementType(c *gin.Context) {
	var auditEngagementType *models.AuditEngagementType

	if err := c.ShouldBindJSON(&auditEngagementType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	createdEngagementType, err := aetc.service.Create(currentUserID, auditEngagementType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aetc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdEngagementType)
}

func (aetc *AuditEngagementTypeController) UpdateAuditEngagementType(c *gin.Context) {
	var updatedEngagementType *models.AuditEngagementType

	auditEngagementTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedEngagementType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	updatedEngagementType, err = aetc.service.Update(uint(auditEngagementTypeID), currentUserID, updatedEngagementType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aetc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedEngagementType)
}

func (aetc *AuditEngagementTypeController) DeleteAuditEngagementTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	failedIDs, err := aetc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aetc *AuditEngagementTypeController) ActivateAuditEngagementTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	err = aetc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (aetc *AuditEngagementTypeController) DeactivateAuditEngagementTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	err = aetc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aetc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
