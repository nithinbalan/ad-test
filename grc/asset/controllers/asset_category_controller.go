package controllers

import (
	"jethur/grc/asset/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetCategoryController struct {
	service services.AssetCategoryServiceInterface
	db      *gorm.DB
}

func NewAssetCategoryController(service services.AssetCategoryServiceInterface) *AssetCategoryController {
	return &AssetCategoryController{
		service: service,
	}
}

func (acc *AssetCategoryController) GetAllAssetCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	assetCategories, err := acc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, assetCategories)
}

func (acc *AssetCategoryController) GetAssetCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetCategory, err := acc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetCategory)
}

func (acc *AssetCategoryController) CreateAssetCategory(c *gin.Context) {
	var assetCategory *models.AssetCategory

	if err := c.ShouldBindJSON(&assetCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	createdCategory, err := acc.service.Create(currentUserID, assetCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (acc *AssetCategoryController) UpdateAssetCategory(c *gin.Context) {
	var updatedCategory *models.AssetCategory

	assetCategoryID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	updatedCategory, err = acc.service.Update(uint(assetCategoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (acc *AssetCategoryController) DeleteAssetCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	failedIDs, err := acc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (acc *AssetCategoryController) ActivateAssetCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (acc *AssetCategoryController) DeactivateAssetCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
