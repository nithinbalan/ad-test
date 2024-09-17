package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentRegisterController struct {
	service services.IncidentRegisterServiceInterface
	db      *gorm.DB
}

func NewIncidentRegisterController(service services.IncidentRegisterServiceInterface) *IncidentRegisterController {
	return &IncidentRegisterController{
		service: service,
	}
}

func (irc *IncidentRegisterController) GetAllIncidentRegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"incident_category_ids",
		"incident_sub_category_ids",
		"severity_level_ids",
		"incident_status_ids",
		"reported_user_ids",
		"incident_type_ids",
	}
	validSortFields := map[string]string{
		"incident_at":                    "incident_register.incident_at",
		"reported_at":                    "incident_register.reported_at",
		"reference_number":               "incident_register.reference_number",
		"title":                          "incident_register.title",
		"description":                    "incident_register.description",
		"incident_category_en":           "incident_categories.title_en",
		"incident_category_ar":           "incident_categories.title_ar",
		"incident_type_en":               "incident_types.title_en",
		"incident_type_ar":               "incident_types.title_ar",
		"incident_status_en":             "incident_statuses.title_en",
		"incident_status_ar":             "incident_statuses.title_ar",
		"severity_level_en":              "severity_levels.title_en",
		"severity_level_ar":              "severity_levels.title_ar",
		"severity_level_color_code":      "severity_levels.color_code",
		"responsible_unit_en":            "organization_units.title_en",
		"responsible_unit_ar":            "organization_units.title_ar",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
		"responsible_user_initial":       "responsible_user.initial",
		"reported_by_first_name_en":      "reported_user.first_name_en",
		"reported_by_first_name_ar":      "reported_user.first_name_ar",
		"reported_by_last_name_en":       "reported_user.last_name_en",
		"reported_by_last_name_ar":       "reported_user.last_name_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	incidentAt := c.Query("incident_at")
	reportedAt := c.Query("reported")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := incidentDTOs.IncidentRegisterListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		IncidentAt:  incidentAt,
		ReportedAt:  reportedAt,
		Export:      export,
	}
	incidents, err := irc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(irc.db, c, err.Error())
			return
		}

		excelFile, err := irc.service.ExportIncidentRegister(incidents.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(irc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Incident Register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, incidents)
}

func (irc *IncidentRegisterController) GetIncidentRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incident, err := irc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incident)
}

func (irc *IncidentRegisterController) CreateIncidentRegister(c *gin.Context) {
	var incidentRegisterFormDTO incidentDTOs.IncidentRegisterFormDTO

	if err := c.ShouldBindJSON(&incidentRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	id, err := irc.service.Create(currentUserID, &incidentRegisterFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (irc *IncidentRegisterController) UpdateIncidentRegister(c *gin.Context) {
	var incidentRegisterFormDTO incidentDTOs.IncidentRegisterFormDTO

	incidentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&incidentRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	err = irc.service.Update(uint(incidentID), currentUserID, &incidentRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (irc *IncidentRegisterController) DeleteIncidentRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	failedIDs, err := irc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (irc *IncidentRegisterController) GetIncidentRegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	irc.GetAllIncidentRegisters(c)
}
