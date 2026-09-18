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

type ProjectSecurityGateReportController struct {
	service                services.ProjectSecurityGateReportServiceInterface
	projectSecurityGateSvc projectServices.ProjectSecurityGateServiceInterface

	db *gorm.DB
}

func NewProjectSecurityGateReportController(
	service services.ProjectSecurityGateReportServiceInterface,
	projectSecurityGateSvc projectServices.ProjectSecurityGateServiceInterface,
) *ProjectSecurityGateReportController {
	return &ProjectSecurityGateReportController{
		service:                service,
		projectSecurityGateSvc: projectSecurityGateSvc,
	}
}

func (psgrc *ProjectSecurityGateReportController) GetAllProjectSecurityGates(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "project_title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"stage_ids",
		"readiness_status_ids",
		"classification_ids",
	}

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	// One entry per column of the excel export
	validSortFields := map[string]string{
		"project_title_en":             "project_security_gates.project_title_en",
		"project_title_ar":             "project_security_gates.project_title_ar",
		"checklist_count":              "checklist_count",
		"stage_title_en":               "stages.title_en",
		"stage_title_ar":               "stages.title_ar",
		"division_title_en":            "division.title_en",
		"division_title_ar":            "division.title_ar",
		"department_title_en":          "department.title_en",
		"department_title_ar":          "department.title_ar",
		"business_owner_first_name_en": "business_owner.first_name_en",
		"business_owner_first_name_ar": "business_owner.first_name_ar",
		"readiness_status_title_en":    "readiness_status.title_en",
		"readiness_status_title_ar":    "readiness_status.title_ar",
		"notes":                        "project_security_gates.notes",
		"created_user_first_name_en":   "created_user.first_name_en",
		"created_user_first_name_ar":   "created_user.first_name_ar",
		"created_at":                   "project_security_gates.created_at",
		"updated_user_first_name_en":   "updated_user.first_name_en",
		"updated_user_first_name_ar":   "updated_user.first_name_ar",
		"updated_at":                   "project_security_gates.updated_at",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, "Invalid page_size parameter")
		return
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid 'export' parameter"))
			return
		}
	}

	queryParams := projectDTOs.ProjectSecurityGateListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		QuickFilter: quickFilter,
		Filters:     utility.ParseQueryFilters(allowedParams, c.Request.URL.Query()),
		Page:        page,
		PageSize:    pageSize,
		Export:      export,
	}

	projectSecurityGates, err := psgrc.projectSecurityGateSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(psgrc.db, c, err.Error())
			return
		}

		excelFile, err := psgrc.projectSecurityGateSvc.ExportProjectSecurityGates(projectSecurityGates.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(psgrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Project Security Gates.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, projectSecurityGates)
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateExport(c *gin.Context) {
	c.Set("export", true)
	psgrc.GetAllProjectSecurityGates(c)
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateByStage(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := psgrc.service.GetProjectSecurityGateByStage(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateByReadinessStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := psgrc.service.GetProjectSecurityGateByReadinessStatus(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateByClassification(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := reportDTOs.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	results, err := psgrc.service.GetProjectSecurityGateByClassification(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}
