package controllers

import (
	"errors"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationReportController struct {
	service services.OrganizationReportServiceInterface
	userSvc organizationServices.UserServiceInterface
	db      *gorm.DB
}

func NewOrganizationReportController(service services.OrganizationReportServiceInterface,
	userSvc organizationServices.UserServiceInterface) *OrganizationReportController {
	return &OrganizationReportController{
		service: service,
		userSvc: userSvc,
	}
}

func (orc *OrganizationReportController) GetAllUsers(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	userIdsStr := c.Query("user_ids")
	roleIdsStr := c.Query("role_ids")

	userIds, err := utility.ParseCommaSeparatedStringToIntArray(userIdsStr)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'user_ids' parameter"))
		return
	}

	roleIds, err := utility.ParseCommaSeparatedStringToIntArray(roleIdsStr)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'role_ids' parameter"))
		return
	}

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"designation_ids",
		"organization_unit_1_ids",
		"organization_unit_2_ids",
		"organization_unit_3_ids",
		"organization_unit_4_ids",
		"organization_unit_5_ids",
		"user_status_ids",
	}

	validSortFields := map[string]string{
		"first_name_en":  "users.first_name_en",
		"first_name_ar":  "users.first_name_ar",
		"email":          "users.email",
		"mobile_number":  "users.mobile_number",
		"designation_en": "designations.title_en",
		"designation_ar": "designations.title_ar",
		"role_en":        "roles_agg.roles_en",
		"role_ar":        "roles_agg.roles_ar",
		"user_status_en": "user_statuses.title_en",
		"user_status_ar": "user_statuses.title_ar",
		"name_en":        "CONCAT(users.first_name_en, ' ', users.last_name_en)",
		"name_ar":        "CONCAT(users.first_name_ar, ' ', users.last_name_ar)",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := organizationDTOs.UserListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		UserIDs:     userIds,
		RoleIDs:     roleIds,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	users, err := orc.userSvc.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	if export {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}

		excelFile, err := orc.userSvc.ExportUsers(users.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(orc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Users.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, users)
}

func (orc *OrganizationReportController) GetUsersExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	orc.GetAllUsers(c)
}

func (orc *OrganizationReportController) GetUsersByStatus(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	usersByStatusReport, err := orc.service.GetUsersByStatus(&queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(orc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, usersByStatusReport)
}

func (orc *OrganizationReportController) GetUsersByRole(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	usersByRoleReport, err := orc.service.GetUsersByRole(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, usersByRoleReport)
}

func (orc *OrganizationReportController) GetUsersByDesignation(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	usersByDesignationReport, err := orc.service.GetUsersByDesignation(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, usersByDesignationReport)
}

func (orc *OrganizationReportController) GetProcessesByResponsibleUnit(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	queryParams := dtos.ReportListQueryParamsDTO{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	processByResponsibleUnitReport, err := orc.service.GetProcessesByResponsibleUnit(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, processByResponsibleUnitReport)
}
