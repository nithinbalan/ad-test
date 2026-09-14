package controllers

import (
	"errors"
	projectDTOs "jethur/grc/project/dtos"
	projectServices "jethur/grc/project/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectReportController struct {
	service              services.ProjectReportServiceInterface
	projectSvc           projectServices.ProjectServiceInterface
	projectActionPlanSvc projectServices.ProjectActionPlanServiceInterface

	db *gorm.DB
}

func NewProjectReportController(
	service services.ProjectReportServiceInterface,
	projectSvc projectServices.ProjectServiceInterface,
	projectActionPlanSvc projectServices.ProjectActionPlanServiceInterface,
) *ProjectReportController {
	return &ProjectReportController{
		service:              service,
		projectSvc:           projectSvc,
		projectActionPlanSvc: projectActionPlanSvc,
	}
}

func (prc *ProjectReportController) GetAllProjects(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"project_type_ids",
		"project_status_ids",
		"project_phase_ids",
		"project_classification_ids",
		"responsible_user_ids",
	}

	validSortFields := map[string]string{
		"title_en":                  "projects.title_en",
		"title_ar":                  "projects.title_ar",
		"description_en":            "projects.description_en",
		"description_ar":            "projects.description_ar",
		"start_date":                "projects.start_date",
		"end_date":                  "projects.end_date",
		"created_at":                "projects.created_at",
		"updated_at":                "projects.updated_at",
		"project_type_en":           "project_types.title_en",
		"project_type_ar":           "project_types.title_ar",
		"project_status_en":         "project_statuses.title_en",
		"project_status_ar":         "project_statuses.title_ar",
		"project_phase_en":          "project_phases.title_en",
		"project_phase_ar":          "project_phases.title_ar",
		"project_classification_en": "project_classifications.title_en",
		"project_classification_ar": "project_classifications.title_ar",
		"budget":                    "projects.budget",
		"number_of_resources":       "projects.number_of_resources",
		"completion_percentage":     "projects.completion_percentage",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page_size parameter")
		return
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := projectDTOs.ProjectListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	projects, err := prc.projectSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}

		excelFile, err := prc.projectSvc.ExportProjects(projects.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Projects.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (prc *ProjectReportController) GetProjectExport(c *gin.Context) {
	c.Set("export", true)
	prc.GetAllProjects(c)
}

func (prc *ProjectReportController) GetProjectByType(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectByType(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectByPhase(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectByPhase(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectByClassification(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectByClassification(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectByYear(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectByYear(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetAllProjectActionPlans(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	if projectIDStr == "" {
		projectIDStr = c.Param("id")
	}

	var projectID uint64
	var err error

	if projectIDStr != "" {
		projectID, err = strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		projectID = 0
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"action_plan_status_ids",
		"responsible_user_ids",
		"review_user_ids",
		"severity_level_ids",
		"project_ids",
	}

	validSortFields := map[string]string{
		"title":               "project_action_plans.title",
		"description":         "project_action_plans.description",
		"project_title":       "projects.title_en",
		"severity_level_en":   "severity_levels.title_en",
		"severity_level_ar":   "severity_levels.title_ar",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"target_date":         "project_action_plans.target_date",
		"percentage":          "project_action_plans.percentage",
		"status_en":           "action_plan_statuses.title_en",
		"status_ar":           "action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		dbField = "title"
	}

	queryParams := projectDTOs.ProjectActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         dbField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		Export:         export,
	}

	projectActionPlans, err := prc.projectActionPlanSvc.GetAll(uint(projectID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}

		excelFile, err := prc.projectActionPlanSvc.ExportProjectActionPlans(projectActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Project Action Plans.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, projectActionPlans)
}

func (prc *ProjectReportController) GetProjectActionPlanExport(c *gin.Context) {
	c.Set("export", true)
	prc.GetAllProjectActionPlans(c)
}

func (prc *ProjectReportController) GetProjectActionPlanByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectActionPlanByStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectActionPlanByResponsibleUser(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectActionPlanByResponsibleUser(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (prc *ProjectReportController) GetProjectActionPlanByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := prc.service.GetProjectActionPlanByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}
