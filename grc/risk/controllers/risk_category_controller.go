package controllers

import (
	"jethur/grc/models"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskCategoryController struct {
	service services.RiskCategoryServiceInterface
	db      *gorm.DB
}

func NewRiskCategoryController(service services.RiskCategoryServiceInterface) *RiskCategoryController {
	return &RiskCategoryController{
		service: service,
	}
}

func (rcc *RiskCategoryController) GetAllRiskCategories(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskCategories, err := rcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskCategories)
}

func (rcc *RiskCategoryController) GetRiskCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskCategory, err := rcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskCategory)
}

func (rcc *RiskCategoryController) CreateRiskCategory(c *gin.Context) {

	var category *models.RiskCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	createdCategory, err := rcc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (rcc *RiskCategoryController) UpdateRiskCategory(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedCategory *models.RiskCategory
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	updatedCategory, err = rcc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (rcc *RiskCategoryController) DeleteRiskCategories(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	failedIDs, err := rcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rcc *RiskCategoryController) ActivateRiskCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rcc *RiskCategoryController) DeactivateRiskCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	err = rcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
