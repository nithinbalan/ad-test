package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessActivityController struct {
	service services.ProcessActivityServiceInterface
	db      *gorm.DB
}

func NewProcessActivityController(service services.ProcessActivityServiceInterface) *ProcessActivityController {
	return &ProcessActivityController{
		service: service,
	}
}

func (pac *ProcessActivityController) GetAllProcessActivities(c *gin.Context) {
	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "sequence")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"created_user_ids",
		"updated_user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"title_en":   "process_activities.title_en",
		"title_ar":   "process_activities.title_ar",
		"sequence":   "process_activities.sequence",
		"created_at": "process_activities.created_at",
		"updated_at": "process_activities.updated_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.ProcessActivityListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	activities, err := pac.service.GetAll(uint(processID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (pac *ProcessActivityController) GetProcessActivityDetails(c *gin.Context) {
	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	activityID, err := strconv.Atoi(c.Param("activityId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := pac.service.GetFullDetailsByID(uint(processID), uint(activityID))
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (pac *ProcessActivityController) CreateProcessActivity(c *gin.Context) {
	var form dtos.ProcessActivityFormDTO

	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	id, err := pac.service.Create(uint(processID), currentUserID, &form)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pac *ProcessActivityController) UpdateProcessActivity(c *gin.Context) {
	var form dtos.ProcessActivityFormDTO

	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	activityID, err := strconv.Atoi(c.Param("activityId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	err = pac.service.Update(uint(processID), uint(activityID), currentUserID, &form)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pac *ProcessActivityController) DeleteProcessActivities(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	failedIDs, err := pac.service.Delete(uint(processID), req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
