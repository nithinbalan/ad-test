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

type ComplianceRequirementCategoryController struct {
	service services.ComplianceRequirementCategoryServiceInterface
	db      *gorm.DB
}

func NewComplianceRequirementCategoryController(service services.ComplianceRequirementCategoryServiceInterface) *ComplianceRequirementCategoryController {
	return &ComplianceRequirementCategoryController{
		service: service,
	}
}

func (crcc *ComplianceRequirementCategoryController) GetAllComplianceRequirementCategories(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	complianceRequirementCategories, err := crcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceRequirementCategories)
}

func (crcc *ComplianceRequirementCategoryController) GetComplianceRequirementCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceCategory, err := crcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, complianceCategory)
}

func (crcc *ComplianceRequirementCategoryController) CreateComplianceRequirementCategory(c *gin.Context) {

	var requirementCategory *models.ComplianceRequirementCategory

	if err := c.ShouldBindJSON(&requirementCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	createdRequirementCategory, err := crcc.service.Create(currentUserID, requirementCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdRequirementCategory)
}

func (crcc *ComplianceRequirementCategoryController) UpdateComplianceRequirementCategory(c *gin.Context) {

	requirementCategoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedRequirementCategory *models.ComplianceRequirementCategory
	if err := c.ShouldBindJSON(&updatedRequirementCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	updatedRequirementCategory, err = crcc.service.Update(uint(requirementCategoryID), currentUserID, updatedRequirementCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedRequirementCategory)
}

func (crcc *ComplianceRequirementCategoryController) DeleteComplianceRequirementCategories(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	failedIDs, err := crcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (crcc *ComplianceRequirementCategoryController) ActivateComplianceRequirementCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	err = crcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (crcc *ComplianceRequirementCategoryController) DeactivateComplianceRequirementCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	err = crcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
