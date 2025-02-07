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

type AssetSubCategoryController struct {
	service services.AssetSubCategoryServiceInterface
	db      *gorm.DB
}

func NewAssetSubCategoryController(service services.AssetSubCategoryServiceInterface) *AssetSubCategoryController {
	return &AssetSubCategoryController{
		service: service,
	}
}

func (ascc *AssetSubCategoryController) GetAllAssetSubCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{
		"asset_category_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := assetDTOs.AssetSubCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	assetSubCategories, err := ascc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, assetSubCategories)
}

func (ascc *AssetSubCategoryController) GetAssetSubCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetSubCategory, err := ascc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetSubCategory)
}

func (ascc *AssetSubCategoryController) CreateAssetSubCategory(c *gin.Context) {
	var category *assetDTOs.AssetSubCategoryFormDTO

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}
	id, err := ascc.service.Create(currentUserID, category)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ascc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, id)
}

func (ascc *AssetSubCategoryController) UpdateAssetSubCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedAssetSubCategory *assetDTOs.AssetSubCategoryFormDTO
	if err := c.ShouldBindJSON(&updatedAssetSubCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}

	err = ascc.service.Update(uint(categoryID), currentUserID, updatedAssetSubCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ascc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ascc *AssetSubCategoryController) DeleteAssetSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}

	failedIDs, err := ascc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ascc *AssetSubCategoryController) ActivateAssetSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}

	err = ascc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ascc *AssetSubCategoryController) DeactivateAssetSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}

	err = ascc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ascc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
