package controllers

import (
	"jethur/grc/models"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPICategoryController struct {
	service services.KPICategoryServiceInterface
	db      *gorm.DB
}

func NewKPICategoryController(service services.KPICategoryServiceInterface) *KPICategoryController {
	return &KPICategoryController{
		service: service,
	}
}

func (kcc *KPICategoryController) GetAllKPICategories(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	categories, err := kcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (kcc *KPICategoryController) GetKPICategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriCategory, err := kcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, kriCategory)
}

func (kcc *KPICategoryController) CreateKPICategory(c *gin.Context) {

	var category *models.KPICategory

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	createdCategory, err := kcc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (kcc *KPICategoryController) UpdateKPICategory(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedCategory *models.KPICategory
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	updatedCategory, err = kcc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (kcc *KPICategoryController) DeleteKPICategories(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	failedIDs, err := kcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (kcc *KPICategoryController) ActivateKPICategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	err = kcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (kcc *KPICategoryController) DeactivateKPICategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	err = kcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(kcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
