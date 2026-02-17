package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/violation-report/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportController struct {
	service services.ViolationReportServiceInterface
	db      *gorm.DB
}

func NewViolationReportController(service services.ViolationReportServiceInterface) *ViolationReportController {
	return &ViolationReportController{
		service: service,
	}
}

func (vrc *ViolationReportController) GetAllViolationReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "reference_number")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Export flag
	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(vrc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"violation_report_category_ids",
		"violation_report_sub_category_ids",
		"violation_report_reason_ids",
		"violation_report_status_ids",
		"violation_report_channel_ids",
		"severity_level_ids",
		"reported_entity_ids",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"reported_at":                            "violation_reports.reported_at",
		"violation_report_channel_title_en":      "violation_reports.violation_report_channel_id",
		"violation_report_channel_title_ar":      "violation_reports.violation_report_channel_id",
		"violation_report_reason_title_en":       "violation_reports.violation_report_reason_id",
		"violation_report_reason_title_ar":       "violation_reports.violation_report_reason_id",
		"violation_report_category_title_en":     "violation_reports.violation_report_category_id",
		"violation_report_category_title_ar":     "violation_reports.violation_report_category_id",
		"violation_report_sub_category_title_en": "violation_reports.violation_report_sub_category_id",
		"violation_report_sub_category_title_ar": "violation_reports.violation_report_sub_category_id",
		"reference_number":                       "violation_reports.reference_number",
		"title_en":                               "violation_reports.title_en",
		"title_ar":                               "violation_reports.title_ar",
		"status_en":                              "violation_report_statuses.title_en",
		"status_ar":                              "violation_report_statuses.title_ar",
		"severity_level_en":                      "severity_levels.title_en",
		"severity_level_ar":                      "severity_levels.title_ar",
		"severity_level_color":                   "severity_levels.color_code",
		"reported_entity_title_en":               "organization_units.title_en",
		"reported_entity_title_ar":               "organization_units.title_ar",
		"created_at":                             "violation_reports.created_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	reportedAt := c.Query("reported_at")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, "Invalid page_size parameter")
		return
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid sort_by parameter"))
		return
	}

	queryParams := violationReportDTOs.ViolationReportListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		ReportedAt:  reportedAt,
		Export:      export,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	reports, err := vrc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(vrc.db, c, err.Error())
			return
		}

		excelFile, err := vrc.service.ExportViolationReports(reports.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(vrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Disposition", `attachment; filename="Violation Reports.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (vrc *ViolationReportController) GetViolationReportDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	report, err := vrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrc.db, c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, report)
}

func (vrc *ViolationReportController) CreateViolationReport(c *gin.Context) {
	var form violationReportDTOs.ViolationReportFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	id, err := vrc.service.Create(currentUserID, &form)
	if err != nil {
		switch err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, err)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrc.db, c, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (vrc *ViolationReportController) UpdateViolationReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form violationReportDTOs.ViolationReportFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	err = vrc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, err)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrc.db, c, err.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (vrc *ViolationReportController) DeleteViolationReports(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	failedIDs, err := vrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (vrc *ViolationReportController) GetViolationReportDocumentPreview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportDocumentPath, err := vrc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrc.db, c, err.Error())
		}
		return
	}
	c.File(violationReportDocumentPath)
}

func (vrc *ViolationReportController) DownloadViolationReportDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportDocumentPath, err := vrc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(violationReportDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(violationReportDocumentPath)
}

func (vrc *ViolationReportController) GetViolationReportExport(c *gin.Context) {
	c.Set("export", true)

	vrc.GetAllViolationReports(c)
}

func (vrc *ViolationReportController) CreateViolationReportComment(c *gin.Context) {

	var dto dtos.ViolationReportCommentFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	violationReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	createdComment, err := vrc.service.CreateViolationReportComment(
		uint(violationReportID),
		currentUserID,
		&dto,
	)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(vrc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

func (vrc *ViolationReportController) GetAllViolationReportComments(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := vrc.service.GetAllViolationReportComments(uint(id))
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (vrc *ViolationReportController) DeleteComment(c *gin.Context) {

	violationReportIDStr := c.Param("id")
	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	failedIDs, err := vrc.service.DeleteComment(
		uint(violationReportID),
		request.IDs,
		currentUserID,
	)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
