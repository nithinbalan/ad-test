package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceCategoryController struct {
	service services.ComplianceCategoryServiceInterface
	db      *gorm.DB
}

func NewComplianceCategoryController(service services.ComplianceCategoryServiceInterface) *ComplianceCategoryController {
	return &ComplianceCategoryController{
		service: service,
	}
}

func (ccc *ComplianceCategoryController) GetAllComplianceCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	complianceCategories, err := ccc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceCategories)
}

func (ccc *ComplianceCategoryController) GetComplianceCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceCategory, err := ccc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, complianceCategory)
}

func (ccc *ComplianceCategoryController) CreateComplianceCategory(c *gin.Context) {
	var category *models.ComplianceCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	createdCategory, err := ccc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (ccc *ComplianceCategoryController) UpdateComplianceCategory(c *gin.Context) {
	var updatedCategory *models.ComplianceCategory

	categoryID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	updatedCategory, err = ccc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (ccc *ComplianceCategoryController) DeleteComplianceCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	failedIDs, err := ccc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ccc *ComplianceCategoryController) ActivateComplianceCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ccc *ComplianceCategoryController) DeactivateComplianceCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
