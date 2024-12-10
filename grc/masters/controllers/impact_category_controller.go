package controllers

import (
	"fmt"
	"jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ImpactCategoryController struct {
	service services.ImpactCategoryServiceInterface
	db      *gorm.DB
}

func NewImpactCategoryController(service services.ImpactCategoryServiceInterface) *ImpactCategoryController {
	return &ImpactCategoryController{
		service: service,
	}
}

func (icc *ImpactCategoryController) GetAllImpactCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	fmt.Println("1")
	impactCategories, err := icc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, impactCategories)
}

func (icc *ImpactCategoryController) GetImpactCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	impactCategory, err := icc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, impactCategory)
}

func (icc *ImpactCategoryController) CreateImpactCategory(c *gin.Context) {
	var category *models.ImpactCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	createdCategory, err := icc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (icc *ImpactCategoryController) UpdateImpactCategory(c *gin.Context) {
	var updatedCategory *models.ImpactCategory

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
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	updatedCategory, err = icc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (icc *ImpactCategoryController) DeleteImpactCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	failedIDs, err := icc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (icc *ImpactCategoryController) ActivateImpactCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	err = icc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (icc *ImpactCategoryController) DeactivateImpactCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	err = icc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
