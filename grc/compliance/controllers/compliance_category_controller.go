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

	sort := c.Query("sort")
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
		utility.SendNotFoundErrResp(ccc.db, c)
		return
	}

	complianceCategory, err := ccc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ccc.db, c)
		return
	}

	// If everything is okay, return the compliance category.
	c.JSON(http.StatusOK, complianceCategory)
}

func (ccc *ComplianceCategoryController) CreateCategory(c *gin.Context) {

	var category models.ComplianceCategory

	// Bind the JSON payload to the ComplianceCategory struct
	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(ccc.db, c, err)
		return
	}

	// Call the service to create the category
	createdCategory, err := ccc.service.Create(category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ccc.db, c, e)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (ccc *ComplianceCategoryController) UpdateCategory(c *gin.Context) {

	// Extract the category ID from the URL parameter
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(ccc.db, c)
		return
	}

	// Bind the incoming JSON to a ComplianceCategory struct
	var updatedCategory models.ComplianceCategory
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		utility.SendInvalidErrResp(ccc.db, c, err)
		return
	}

	// Update the category using the service layer
	updatedCategory, err = ccc.service.Update(uint(categoryID), updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ccc.db, c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ccc.db, c, e)
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
		utility.SendInvalidErrResp(ccc.db, c, err)
		return
	}

	failedIDs, err := ccc.service.Delete(request.IDs)
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
		utility.SendInvalidErrResp(ccc.db, c, err)
		return
	}

	err := ccc.service.Activate(request.IDs)
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
		utility.SendInvalidErrResp(ccc.db, c, err)
		return
	}

	err := ccc.service.Deactivate(request.IDs)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
