package controllers

import (
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceController struct {
	service services.ComplianceServiceInterface
	db      *gorm.DB
}

func NewComplianceController(service services.ComplianceServiceInterface) *ComplianceController {
	return &ComplianceController{
		service: service,
	}
}

func (cc *ComplianceController) GetAllCompliances(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	compliances, err := cc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, compliances)
}

func (cc *ComplianceController) GetComplianceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	compliance, err := cc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, compliance)
}

func (cc *ComplianceController) CreateCompliance(c *gin.Context) {

	var compliance *models.ComplianceRegister

	if err := c.ShouldBindJSON(&compliance); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	createdCompliance, err := cc.service.Create(currentUserID, compliance)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCompliance)
}

func (cc *ComplianceController) UpdateCompliance(c *gin.Context) {

	complianceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedCompliance *models.ComplianceRegister
	if err := c.ShouldBindJSON(&updatedCompliance); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	updatedCompliance, err = cc.service.Update(uint(complianceID), currentUserID, updatedCompliance)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCompliance)
}

func (cc *ComplianceController) DeleteCompliances(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	failedIDs, err := cc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cc *ComplianceController) ActivateCompliances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cc *ComplianceController) DeactivateCompliances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
