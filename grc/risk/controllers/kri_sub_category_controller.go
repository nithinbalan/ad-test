package controllers

import (
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KRISubCategoryController struct {
	service services.KRISubCategoryServiceInterface
	db      *gorm.DB
}

func NewKRISubCategoryController(service services.KRISubCategoryServiceInterface) *KRISubCategoryController {
	return &KRISubCategoryController{service: service}
}

func (ksc *KRISubCategoryController) GetAllKRISubCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{
		"kri_category_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	queryParams := dtos.KRISubCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	items, err := ksc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ksc *KRISubCategoryController) GetKRISubCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	item, err := ksc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ksc *KRISubCategoryController) CreateKRISubCategory(c *gin.Context) {
	var item *dtos.KRISubCategoryFormDTO
	if err := c.ShouldBindJSON(&item); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	id, err := ksc.service.Create(currentUserID, item)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ksc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (ksc *KRISubCategoryController) UpdateKRISubCategory(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedItem *dtos.KRISubCategoryFormDTO
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	err = ksc.service.Update(uint(itemID), currentUserID, updatedItem)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ksc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ksc *KRISubCategoryController) DeleteKRISubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	failedIDs, err := ksc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ksc *KRISubCategoryController) ActivateKRISubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	if err := ksc.service.Activate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ksc *KRISubCategoryController) DeactivateKRISubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	if err := ksc.service.Deactivate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ksc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
