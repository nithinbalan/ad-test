package controllers

import (
	"jethur/grc/libraries/models"
	"jethur/grc/libraries/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRICategoryController struct {
	service services.KRICategoryServiceInterface
	db      *gorm.DB
}

func NewKRICategoryController(service services.KRICategoryServiceInterface) *KRICategoryController {
	return &KRICategoryController{
		service: service,
	}
}

func (kcc *KRICategoryController) GetAllKRICategories(c *gin.Context) {

	sort := c.Query("sort")
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

func (kcc *KRICategoryController) GetKRICategoryDetails(c *gin.Context) {
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

func (kcc *KRICategoryController) CreateKRICategory(c *gin.Context) {

	var category *models.KRICategory

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

func (kcc *KRICategoryController) UpdateKRICategory(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedCategory *models.KRICategory
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

func (kcc *KRICategoryController) DeleteKRICategories(c *gin.Context) {

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

func (kcc *KRICategoryController) ActivateKRICategories(c *gin.Context) {
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

func (kcc *KRICategoryController) DeactivateKRICategories(c *gin.Context) {
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
