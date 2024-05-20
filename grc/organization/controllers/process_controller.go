package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProcessController struct {
	service services.ProcessServiceInterface
	db      *gorm.DB
}

func NewProcessController(service services.ProcessServiceInterface) *ProcessController {
	return &ProcessController{
		service: service,
	}
}

func (pc *ProcessController) GetAllProcesses(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "process_title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_group_ids",
		"owner_unit_ids",
		"owner_ids",
		"lifecycle_status_ids",
	}

	validSortFields := map[string]string{
		"process_title":   "processes.title_en",
		"refrence_number": "processes.reference_number",
		"description":     "processes.description",
		"process_group":   "process_groups.title_en",
		"status":          "lifecycle_statuses.title_enn",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.ProcessListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	processes, err := pc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, processes)
}

func (pc *ProcessController) GetProcessDetails(c *gin.Context) {
	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	process, err := pc.service.GetFullDetailsByID(uint(processID))
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, process)
}

func (pc *ProcessController) CreateProcess(c *gin.Context) {
	var processFormDTO organizationDTOs.ProcessFormDTO
	if err := c.ShouldBindJSON(&processFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	id, err := pc.service.Create(currentUserID, &processFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pc *ProcessController) UpdateProcess(c *gin.Context) {
	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var processFormDTO organizationDTOs.ProcessFormDTO
	if err := c.ShouldBindJSON(&processFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	err = pc.service.Update(uint(processID), uint(currentUserID), &processFormDTO)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pc *ProcessController) DeleteProcesses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	failedIDs, err := pc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pc *ProcessController) ArchiveProcesses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	successfulIds, err := pc.service.Archive(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}

func (pc *ProcessController) RestoreProcesses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	successfulIds, err := pc.service.Restore(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}
