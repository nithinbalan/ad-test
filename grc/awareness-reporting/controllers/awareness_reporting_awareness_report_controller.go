package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	"jethur/grc/awareness-reporting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessReportController struct {
	service services.AwarenessReportServiceInterface
	db      *gorm.DB
}

func NewAwarenessReportController(service services.AwarenessReportServiceInterface) *AwarenessReportController {
	return &AwarenessReportController{
		service: service,
	}
}

func (arc *AwarenessReportController) GetAllAwarenessReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "reference_number")
	keyword := c.Query("keyword")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(arc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"awareness_report_type_ids",
		"sub_category_ids",
		"severity_ids",
		"awareness_handling_status_ids",
		"handling_owner_ids",
	}

	validSortFields := map[string]string{
		"reference_number":                       "awareness_reports.reference_number",
		"report_source_channel":                  "awareness_reports.report_source_channel",
		"report_date":                            "awareness_reports.report_date",
		"resulted_from_campaign":                 "awareness_reports.resulted_from_campaign",
		"closure_date":                           "awareness_reports.closure_date",
		"awareness_report_type_title_en":         "awareness_campaign_types.title_en",
		"awareness_report_type_title_ar":         "awareness_campaign_types.title_ar",
		"sub_category_title_en":                  "awareness_report_sub_categories.title_en",
		"sub_category_title_ar":                  "awareness_report_sub_categories.title_ar",
		"severity_title_en":                      "severity_levels.title_en",
		"severity_title_ar":                      "severity_levels.title_ar",
		"awareness_handling_status_title_en":     "awareness_handling_statuses.title_en",
		"awareness_handling_status_title_ar":     "awareness_handling_statuses.title_ar",
		"handling_owner_first_name_en":           "handling_owner_user.first_name_en",
		"handling_owner_first_name_ar":           "handling_owner_user.first_name_ar",
		"handling_owner_last_name_en":            "handling_owner_user.last_name_en",
		"handling_owner_last_name_ar":            "handling_owner_user.last_name_ar",
		"handling_owner_email":                   "handling_owner_user.email",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, "Invalid page_size parameter")
		return
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := awarenessDTOs.AwarenessReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	reports, err := arc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (arc *AwarenessReportController) GetAwarenessReportDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	report, err := arc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, report)
}

func (arc *AwarenessReportController) CreateAwarenessReport(c *gin.Context) {
	var formDTO awarenessDTOs.AwarenessReportFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	id, err := arc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (arc *AwarenessReportController) UpdateAwarenessReport(c *gin.Context) {
	var formDTO awarenessDTOs.AwarenessReportFormDTO

	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	err = arc.service.Update(uint(reportID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arc *AwarenessReportController) DeleteAwarenessReports(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	failedIDs, err := arc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (arc *AwarenessReportController) GetAwarenessReportExport(c *gin.Context) {
	c.Set("export", true)

	arc.GetAllAwarenessReports(c)
}