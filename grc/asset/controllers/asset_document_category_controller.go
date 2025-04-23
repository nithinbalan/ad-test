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

type AssetDocumentCategoryController struct {
	service services.AssetDocumentCategoryServiceInterface
	db      *gorm.DB
}

func NewAssetDocumentCategoryController(service services.AssetDocumentCategoryServiceInterface) *AssetDocumentCategoryController {
	return &AssetDocumentCategoryController{
		service: service,
	}
}

func (adcc *AssetDocumentCategoryController) GetAllAssetDocumentCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := assetDTOs.AssetDocumentCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	assetDocumentCategories, err := adcc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, assetDocumentCategories)
}

func (adcc *AssetDocumentCategoryController) GetAssetDocumentCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetDocumentCategory, err := adcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetDocumentCategory)
}

func (adcc *AssetDocumentCategoryController) CreateAssetDocumentCategory(c *gin.Context) {
	var category *assetDTOs.AssetDocumentCategoryFormDTO

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}
	id, err := adcc.service.Create(currentUserID, category)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adcc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, id)
}

func (adcc *AssetDocumentCategoryController) UpdateAssetDocumentCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedAssetDocumentCategory *assetDTOs.AssetDocumentCategoryFormDTO
	if err := c.ShouldBindJSON(&updatedAssetDocumentCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	err = adcc.service.Update(uint(categoryID), currentUserID, updatedAssetDocumentCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adcc *AssetDocumentCategoryController) DeleteAssetDocumentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	failedIDs, err := adcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (adcc *AssetDocumentCategoryController) ActivateAssetDocumentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	err = adcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (adcc *AssetDocumentCategoryController) DeactivateAssetDocumentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	err = adcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
