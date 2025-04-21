package controllers

import (
	"errors"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	processServices "jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessImpactAssessmentController struct {
	service                  services.BusinessImpactAssessmentServiceInterface
	processAssociatedItemSvc processServices.ProcessAssociatedItemServiceInterface
	processSvc               processServices.ProcessServiceInterface
	processDependencySvc     processServices.ProcessDependencyServiceInterface
	db                       *gorm.DB
}

func NewBusinessImpactAssessmentController(service services.BusinessImpactAssessmentServiceInterface,
	processAssociatedItemSvc processServices.ProcessAssociatedItemServiceInterface,
	processDependencySvc processServices.ProcessDependencyServiceInterface,
	processSvc processServices.ProcessServiceInterface) *BusinessImpactAssessmentController {
	return &BusinessImpactAssessmentController{
		service:                  service,
		processAssociatedItemSvc: processAssociatedItemSvc,
		processDependencySvc:     processDependencySvc,
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
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Export:      export,
	}

	processes, err := biac.processSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	if export {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(biac.db, c, err.Error())
			return
		}

		excelFile, err := biac.processSvc.ExportProcesses(processes.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(biac.db, c, err.Error())
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
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
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
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
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
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
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

func (biac *BusinessImpactAssessmentController) GetAllProcessRisks(c *gin.Context) {
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
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := biac.processAssociatedItemSvc.GetAllRisks(uint(processID), &queryParams)

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

func (biac *BusinessImpactAssessmentController) GetAllProcessDependencies(c *gin.Context) {
	processIDStr := c.Param("processId")

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	allowedParams := []string{
		"process_dependency_sub_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := organizationDTOs.ProcessDependencyListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Filters: filters,
	}
	result, err := biac.processDependencySvc.GetAll(uint(processID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (biac *BusinessImpactAssessmentController) GetAllBIAImpactAnalysis(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "process_title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_ids",
		"risk_level_ids",
		"responsible_unit_ids",
		"owner_ids",
		"process_criticality_level_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	validSortFields := map[string]string{
		"reference_number":        "business_impact_assessments.reference_number",
		"process_title_en":        "processes.title_en",
		"process_title_ar":        "processes.title_ar",
		"asset_register_title_en": "asset_register.title_en",
		"asset_register_title_ar": "asset_register.title_ar",
		"risk_level_title_en":     "risk_levels.title_en",
		"risk_level_title_ar":     "risk_levels.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
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

func (biac *BusinessImpactAssessmentController) GetProcessDetails(c *gin.Context) {
	processID, err := strconv.Atoi(c.Param("processId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	process, err := biac.processSvc.GetFullDetailsByID(uint(processID))
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, process)
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

func (biac *BusinessImpactAssessmentController) UpdateProcessCriticality(c *gin.Context) {
	var businessImpactAssessmentProcessCriticalityFormDTO businessContinuityDTOs.BusinessImpactAssessmentProcessCriticalityFormDTO

	if err := c.ShouldBindJSON(&businessImpactAssessmentProcessCriticalityFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	err = biac.service.UpdateProcessCriticality(businessImpactAssessmentProcessCriticalityFormDTO, currentUserID)
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

func (biac *BusinessImpactAssessmentController) DeleteBusinessImpactAssessments(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	failedIDs, err := biac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (biac *BusinessImpactAssessmentController) GetBIAMatrix(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "business_function_title")
	keyword := c.Query("keyword")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid export parameter")
		return
	}

	validSortFields := map[string]string{
		"business_function_title": "business_functions.title",
		"process_title_en":        "processes.title_en",
		"process_title_ar":        "processes.title_ar",
		"risk_level_title_en":     "risk_levels.title_en",
		"risk_level_title_ar":     "risk_levels.title_ar",
		"rto":                     "processes.rto",
		"mtp":                     "processes.mtpd",
		"rpo":                     "processes.rpo",
		"rto_unit_title_en":       "time_units_rto.title_en",
		"rto_unit_title_ar":       "time_units_rto.title_ar",
		"mtp_unit_title_en":       "time_units_mtpd.title_en",
		"mtp_unit_title_ar":       "time_units_mtpd.title_ar",
		"rpo_unit_title_en":       "time_units_rpo.title_en",
		"rpo_unit_title_ar":       "time_units_rpo.title_ar",
		"risk_level_id":           "processes.risk_level_id",
		"created_at":              "processes.created_at",
		"updated_at":              "processes.updated_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	filters := utility.ParseQueryFilters([]string{"risk_level_id"}, c.Request.URL.Query())

	queryParams := businessContinuityDTOs.BIAMatrixQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
		Filters: filters,
		Export:  export,
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

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(biac.db, c, err.Error())
			return
		}

		excelFile, err := biac.service.ExportBiaMatrix(matrix, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(biac.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="BiaMatrix.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, matrix)
}

func (biac *BusinessImpactAssessmentController) GetBiaMatrixExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	biac.GetBIAMatrix(c)
}

func (biac *BusinessImpactAssessmentController) GetAllProcessesExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	biac.GetAllProcesses(c)
}
