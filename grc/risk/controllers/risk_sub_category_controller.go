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

type RiskSubCategoryController struct {
	service services.RiskSubCategoryServiceInterface
	db      *gorm.DB
}

func NewRiskSubCategoryController(service services.RiskSubCategoryServiceInterface) *RiskSubCategoryController {
	return &RiskSubCategoryController{service: service}
}

func (rsc *RiskSubCategoryController) GetAllRiskSubCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{
		"risk_category_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	queryParams := dtos.RiskSubCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	items, err := rsc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

func (rsc *RiskSubCategoryController) GetRiskSubCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	item, err := rsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (rsc *RiskSubCategoryController) CreateRiskSubCategory(c *gin.Context) {
	var item *dtos.RiskSubCategoryFormDTO
	if err := c.ShouldBindJSON(&item); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	id, err := rsc.service.Create(currentUserID, item)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (rsc *RiskSubCategoryController) UpdateRiskSubCategory(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedItem *dtos.RiskSubCategoryFormDTO
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	err = rsc.service.Update(uint(itemID), currentUserID, updatedItem)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rsc *RiskSubCategoryController) DeleteRiskSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	failedIDs, err := rsc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rsc *RiskSubCategoryController) ActivateRiskSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	if err := rsc.service.Activate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rsc *RiskSubCategoryController) DeactivateRiskSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	if err := rsc.service.Deactivate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
