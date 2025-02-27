package controllers

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
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
	allowedParams := []string{}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := assetDTOs.AssetCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	assetCategories, err := acc.service.GetAll(&queryParams)

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
	var category *assetDTOs.AssetCategoryFormDTO

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}
	id, err := acc.service.Create(currentUserID, category)

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
	c.JSON(http.StatusCreated, id)
}

func (acc *AssetCategoryController) UpdateAssetCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedAssetCategory *assetDTOs.AssetCategoryFormDTO
	if err := c.ShouldBindJSON(&updatedAssetCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Update(uint(categoryID), currentUserID, updatedAssetCategory)
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

	c.Status(http.StatusOK)
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
