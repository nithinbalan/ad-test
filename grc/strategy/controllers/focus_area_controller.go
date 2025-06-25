package controllers

import (
	"fmt"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FocusAreaController struct {
	service services.FocusAreaServiceInterface
	db      *gorm.DB
}

func NewFocusAreaController(service services.FocusAreaServiceInterface) *FocusAreaController {
	return &FocusAreaController{
		service: service,
	}
}

func (fac *FocusAreaController) GetAllFocusAreas(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, "Invalid page parameter")
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"id":       "focus_areas.id",
		"title_en": "focus_areas.title_en",
		"title_ar": "focus_areas.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid 'sort_by' parameter"))
		return
	}

	queryParams := dtos.FocusAreaListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	list, err := fac.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

func (fac *FocusAreaController) GetFocusAreaDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	focusArea, err := fac.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, focusArea)
}

func (fac *FocusAreaController) CreateFocusArea(c *gin.Context) {
	var focusAreaFormDTO *dtos.FocusAreaFormDTO

	if err := c.ShouldBindJSON(&focusAreaFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}
	createdFocusAreaID, err := fac.service.Create(currentUserID, focusAreaFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(fac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdFocusAreaID)
}

func (fac *FocusAreaController) UpdateFocusArea(c *gin.Context) {
	focusAreaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var focusAreaFormDTO *dtos.FocusAreaFormDTO
	if err := c.ShouldBindJSON(&focusAreaFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}

	err = fac.service.Update(uint(focusAreaID), currentUserID, focusAreaFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(fac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (fac *FocusAreaController) DeleteFocusAreas(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}

	failedIDs, err := fac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (fac *FocusAreaController) ActivateFocusAreas(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}

	err = fac.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (fac *FocusAreaController) DeactivateFocusAreas(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}

	err = fac.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(fac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
