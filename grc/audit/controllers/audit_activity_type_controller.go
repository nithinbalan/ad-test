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

type AuditActivityTypeController struct {
	service services.AuditActivityTypeServiceInterface
	db      *gorm.DB
}

func NewAuditActivityTypeController(service services.AuditActivityTypeServiceInterface) *AuditActivityTypeController {
	return &AuditActivityTypeController{
		service: service,
	}
}

func (aatc *AuditActivityTypeController) GetAllAuditActivityTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditActivityTypes, err := aatc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditActivityTypes)
}

func (aatc *AuditActivityTypeController) GetAuditActivityTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditCategory, err := aatc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditCategory)
}

func (aatc *AuditActivityTypeController) CreateAuditActivityType(c *gin.Context) {
	var auditActivityType *models.AuditActivityType

	if err := c.ShouldBindJSON(&auditActivityType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	createdCategory, err := aatc.service.Create(currentUserID, auditActivityType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aatc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (aatc *AuditActivityTypeController) UpdateAuditActivityType(c *gin.Context) {
	var updatedCategory *models.AuditActivityType

	auditActivityTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	updatedCategory, err = aatc.service.Update(uint(auditActivityTypeID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aatc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (aatc *AuditActivityTypeController) DeleteAuditActivityTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	failedIDs, err := aatc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aatc *AuditActivityTypeController) ActivateAuditActivityTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	err = aatc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (aatc *AuditActivityTypeController) DeactivateAuditActivityTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	err = aatc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
