package controllers

import (
	"errors"
	"net/http"
	"strconv"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	processServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessImpactAssessmentController struct {
	service                  services.BusinessImpactAssessmentServiceInterface
	processAssociatedItemSvc processServices.ProcessAssociatedItemServiceInterface
	processSvc               processServices.ProcessServiceInterface
	db                       *gorm.DB
}

func NewBusinessImpactAssessmentController(service services.BusinessImpactAssessmentServiceInterface,
	processAssociatedItemSvc processServices.ProcessAssociatedItemServiceInterface,
	processSvc processServices.ProcessServiceInterface,
) *BusinessImpactAssessmentController {
	return &BusinessImpactAssessmentController{
		service:                  service,
		processAssociatedItemSvc: processAssociatedItemSvc,
		processSvc:               processSvc,
	}
}

func (biac *BusinessImpactAssessmentController) GetAllProcesses(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid export parameter")
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
		utility.SendInternalErrResp(biac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page_size parameter")
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

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
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	processes, err := biac.processSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, processes)
}

func (biac *BusinessImpactAssessmentController) GetAllProcessIncidents(c *gin.Context) {
	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ProcessAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	incidents, err := biac.processAssociatedItemSvc.GetAllIncidents(uint(processID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidents)
}

func (biac *BusinessImpactAssessmentController) GetAllProcessDocuments(c *gin.Context) {
	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ProcessAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	documents, err := biac.processAssociatedItemSvc.GetAllDocuments(uint(processID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}

func (biac *BusinessImpactAssessmentController) GetAllProcessAssets(c *gin.Context) {
	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ProcessAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	assets, err := biac.processAssociatedItemSvc.GetAllAssets(uint(processID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, assets)
}

func (biac *BusinessImpactAssessmentController) GetAllBIAImpactAnalysis(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_ids",
		"risk_level_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO{
		Sort:    sort,
		Keyword: keyword,
		Filters: filters,
	}
	businessImpactAssessments, err := biac.service.GetAllImpactAnalysis(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, businessImpactAssessments)
}

func (biac *BusinessImpactAssessmentController) GetAllBIARecoveryObjectives(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO{
		Sort:    sort,
		Keyword: keyword,
		Filters: filters,
	}
	businessImpactAssessments, err := biac.service.GetAllRecoveryObjectives(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, businessImpactAssessments)
}

func (biac *BusinessImpactAssessmentController) CreateBusinessImpactAssessment(c *gin.Context) {
	var businessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactFormDTO

	if err := c.ShouldBindJSON(&businessImpactAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	id, err := biac.service.Create(currentUserID, &businessImpactAssessmentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (biac *BusinessImpactAssessmentController) UpdateImpactAnalysis(c *gin.Context) {
	var BusinessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO

	if err := c.ShouldBindJSON(&BusinessImpactAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	err = biac.service.UpdateImpactAnalysis(BusinessImpactAssessmentFormDTO, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (biac *BusinessImpactAssessmentController) UpdateRecoveryObjective(c *gin.Context) {
	var businessImpactAssessmentRecoveryObjectiveFormDTO businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveFormDTO

	if err := c.ShouldBindJSON(&businessImpactAssessmentRecoveryObjectiveFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	err = biac.service.UpdateRecoveryObjectives(businessImpactAssessmentRecoveryObjectiveFormDTO, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (biac *BusinessImpactAssessmentController) GetBIAMatrix(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "business_function_title")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page parameter")
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"business_function_title": "business_functions.title",
		"process_title_en":        "processes.title_en",
		"risk_level_title_en":     "risk_levels.title_en",
		"rto":                     "processes.rto",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	filters := utility.ParseQueryFilters([]string{"risk_level_id"}, c.Request.URL.Query())

	queryParams := businessContinuityDTOs.BIAMatrixQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	matrix, err := biac.service.GetBIAMatrix(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, matrix)
}
