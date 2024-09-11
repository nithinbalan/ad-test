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
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_group_ids",
		"owner_unit_ids",
		"owner_ids",
		"lifecycle_status_ids",
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

	validSortFields := map[string]string{
		"reference_number":     "processes.reference_number",
		"title_en":             "processes.title_en",
		"title_ar":             "processes.title_ar",
		"description":          "processes.description",
		"process_group_en":     "process_groups.title_en",
		"process_group_ar":     "process_groups.title_ar",
		"organization_unit_en": "organization_units.title_en",
		"organization_unit_ar": "organization_units.title_ar",
		"owner_en":             "owner.first_name_en",
		"owner_ar":             "owner.first_name_ar",
		"status_en":            "lifecycle_statuses.title_en",
		"status_ar":            "lifecycle_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.ProcessListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	processes, err := pc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	if export {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(pc.db, c, err.Error())
			return
		}

		excelFile, err := pc.service.ExportProcesses(processes.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(pc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Processes.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

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

func (pc *ProcessController) GetControlSummary(c *gin.Context) {
	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	process, err := pc.service.GetControlSummary(uint(processID))
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, process)
}

func (pc *ProcessController) GetProcessesExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	pc.GetAllProcesses(c)
}
