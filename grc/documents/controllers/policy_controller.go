package controllers

import (
	"errors"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyController struct {
	service services.DocumentServiceInterface
	db      *gorm.DB
}

func NewPolicyController(service services.DocumentServiceInterface) *PolicyController {
	return &PolicyController{
		service: service,
	}
}

func (pc *PolicyController) GetAllPolicies(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	enableAcknowledgedStatus := c.Query("enable_acknowledge_status")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"document_category_ids",
		"document_type_ids",
		"document_access_type_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"document_status_ids",
	}

	validSortFields := map[string]string{
		"type_ar":             "document_types.title_ar",
		"type_en":             "document_types.title_en",
		"title_ar":            "documents.title_ar",
		"title_en":            "documents.title_en",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"category_en":         "document_categories.title_en",
		"category_ar":         "document_categories.title_ar",
		"status_en":           "document_statuses.title_ar",
		"status_ar":           "document_statuses.title_en",
		"created_at":          "documents.created_at",
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

	queryParams := documentDTOs.DocumentListQueryParamsDTO{
		Sort:                     sort,
		SortBy:                   dbField,
		Keyword:                  keyword,
		Status:                   status,
		Page:                     page,
		PageSize:                 pageSize,
		Filters:                  filters,
		QuickFilter:              quickFilter,
		EnableAcknowledgedStatus: enableAcknowledgedStatus == "true",
		DocumentTypes:            []string{"policy"},
		DocumentStatus:           []string{"published"},
		Export:                   export,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pc.db, c, err.Error())
		return
	}

	policies, err := pc.service.GetAll(&queryParams, currentUserID)
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

		excelFile, err := pc.service.ExportDocuments(policies.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(pc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Document.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, policies)
}

func (pc *PolicyController) GetPolicyExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	pc.GetAllPolicies(c)
}
