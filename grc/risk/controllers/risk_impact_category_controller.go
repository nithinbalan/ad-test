package controllers

import (
	"jethur/grc/risk/models"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskImpactCategoryController struct {
	service services.RiskImpactCategoryServiceInterface
	db      *gorm.DB
}

func NewRiskImpactCategoryController(service services.RiskImpactCategoryServiceInterface) *RiskImpactCategoryController {
	return &RiskImpactCategoryController{
		service: service,
	}
}

func (ricc *RiskImpactCategoryController) GetAllRiskImpactCategories(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskImpactCategories, err := ricc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskImpactCategories)
}

func (ricc *RiskImpactCategoryController) GetRiskImpactCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskImpactCategory, err := ricc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskImpactCategory)
}

func (ricc *RiskImpactCategoryController) CreateRiskImpactCategory(c *gin.Context) {

	var category *models.RiskImpactCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	createdCategory, err := ricc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ricc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (ricc *RiskImpactCategoryController) UpdateRiskImpactCategory(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedCategory *models.RiskImpactCategory
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	updatedCategory, err = ricc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ricc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (ricc *RiskImpactCategoryController) DeleteRiskImpactCategories(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	failedIDs, err := ricc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ricc *RiskImpactCategoryController) ActivateRiskImpactCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	err = ricc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ricc *RiskImpactCategoryController) DeactivateRiskImpactCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	err = ricc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ricc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
