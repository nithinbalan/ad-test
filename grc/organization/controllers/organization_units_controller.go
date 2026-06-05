package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationUnitsController struct {
	service services.OrganizationUnitsServiceInterface
	db      *gorm.DB
}

func NewOrganizationUnitsController(
	service services.OrganizationUnitsServiceInterface,
) *OrganizationUnitsController {
	return &OrganizationUnitsController{
		service: service,
	}
}

func (osc *OrganizationUnitsController) GetOrganizationUnitsHierarchies(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")

	organizationStructure, err := osc.service.GetOrganizationUnitsHierarchies(sort, sortBy)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, organizationStructure)
}

func (osc *OrganizationUnitsController) GetUsersByOrganizationUnitID(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	// Get organization unit ID from URL parameter
	unitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"first_name_en":  "users.first_name_en",
		"first_name_ar":  "users.first_name_ar",
		"designation_en": "designations.title_en",
		"designation_ar": "designations.title_ar",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Page:     page,
		PageSize: pageSize,
	}

	// Call service to get users by organization unit ID
	users, err := osc.service.GetUsersByOrganizationUnit(&queryParams, uint(unitID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, users)
}

func (osc *OrganizationUnitsController) GetOrganizationUnitSummary(c *gin.Context) {
	unitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := osc.service.GetOrganizationUnitSummary(uint(unitID))
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (osc *OrganizationUnitsController) UpdatePrimaryOrganizationUnitHead(c *gin.Context) {

	unitID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	err = osc.service.UpdateOrganizationUnitHead(uint(unitID), uint(userID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
		}
		return
	}
	c.Status(http.StatusOK)
}

func (osc *OrganizationUnitsController) GetAllPoliciesByOrganizationUnit(c *gin.Context) {
	unitIDStr := c.Param("id")

	unitID, err := strconv.ParseUint(unitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
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
		"created_by_en":       "CONCAT_WS(' ', created_user.first_name_en, created_user.last_name_en)",
		"created_by_ar":       "CONCAT_WS(' ', created_user.first_name_ar, created_user.last_name_ar)",
		"responsible_unit_en": "organization_units.title_en",
		"responsible_unit_ar": "organization_units.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	policies, err := osc.service.GetAllPoliciesByOrganizationUnit(uint(unitID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, policies)
}

func (osc *OrganizationUnitsController) GetAllProcessesByOrganizationUnit(c *gin.Context) {
	unitIDStr := c.Param("id")

	unitID, err := strconv.ParseUint(unitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_group_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"lifecycle_status_ids",
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

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	processes, err := osc.service.GetAllProcessesByOrganizationUnit(uint(unitID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, processes)
}

func (osc *OrganizationUnitsController) GetAllRiskByOrganizationUnit(c *gin.Context) {

	unitIDStr := c.Param("id")

	unitID, err := strconv.ParseUint(unitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"risk_source_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"risk_status_ids",
		"review_frequency_ids",
		"inherent_risk_rating_ids",
		"residual_risk_rating_ids",
		"current_likelihood_ids",
		"current_impact_ids",
		"current_risk_rating_ids",
		"risk_response_strategy_ids",
	}

	validSortFields := map[string]string{
		"reference_number":          "risk_register.reference_number",
		"title":                     "risk_register.title",
		"description":               "risk_register.description",
		"risk_category_en":          "risk_categories.title_en",
		"risk_category_ar":          "risk_categories.title_ar",
		"organization_unit_en":      "organization_units.title_en",
		"organization_unit_ar":      "organization_units.title_ar",
		"responsible_user_en":       "users.first_name_en",
		"responsible_user_ar":       "users.first_name_ar",
		"risk_source_en":            "risk_sources.title_en",
		"risk_source_ar":            "risk_sources.title_ar",
		"risk_status_en":            "risk_statuses.title_en",
		"risk_status_ar":            "risk_statuses.title_ar",
		"inherent_likelihood_en":    "inherent_likelihood.title_en",
		"inherent_likelihood_ar":    "inherent_likelihood.title_ar",
		"inherent_impact_en":        "inherent_impact.title_en",
		"inherent_impact_ar":        "inherent_impact.title_ar",
		"inherent_risk_rating_en":   "inherent_risk_rating.title_en",
		"inherent_risk_rating_ar":   "inherent_risk_rating.title_ar",
		"risk_response_strategy_en": "risk_response_strategies.title_en",
		"risk_response_strategy_ar": "risk_response_strategies.title_ar",
		"residual_likelihood_en":    "residual_likelihood.title_en",
		"residual_likelihood_ar":    "residual_likelihood.title_ar",
		"residual_impact_en":        "residual_impact.title_en",
		"residual_impact_ar":        "residual_impact.title_ar",
		"residual_risk_rating_en":   "residual_risk_rating.title_en",
		"residual_risk_rating_ar":   "residual_risk_rating.title_ar",
		"reported_at":               "risk_register.reported_at",
		"review_frequency_en":       "review_frequencies.title_en",
		"review_frequency_ar":       "review_frequencies.title_ar",
		"next_review_date":          "risk_register.next_review_date",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(osc.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Value:       value,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	risks, err := osc.service.GetAllRiskByOrganizationUnit(uint(unitID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risks)
}

func (osc *OrganizationUnitsController) GetParentOrganizationUnitAndCurrentUnitHead(c *gin.Context) {
	unitIDStr := c.Param("id")

	unitID, err := strconv.ParseUint(unitIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := osc.service.GetParentOrganizationUnitAndCurrentUnitHead(uint(unitID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, result)
}

func (osc *OrganizationUnitsController) GetAllOrganizationUnitReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "organization_unit_reports.title",
		"start_date": "organization_unit_reports.start_date",
		"end_date":   "organization_unit_reports.end_date",
		"created_at": "organization_unit_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.OrganizationUnitReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := osc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (osc *OrganizationUnitsController) CreateOrganizationUnitReport(c *gin.Context) {
	var request organizationDTOs.OrganizationUnitReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	organnizationUnitReportID, err := osc.service.CreateOrganizationUnitReport(currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":                     "Organization unit report generated successfully",
		"organization_unit_report_id": organnizationUnitReportID,
	})
}

func (osc *OrganizationUnitsController) GetOrganizationUnitReportDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	report, err := osc.service.GetOrganizationUnitReport(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, report)
}

func (osc *OrganizationUnitsController) DeleteOrganizationUnitReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	err = osc.service.DeleteOrganizationUnitReport(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(osc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (osc *OrganizationUnitsController) ExportOrganizationUnitReportToPDF(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := strconv.ParseUint(IDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(osc.db, c, err.Error())
		return
	}

	// Delegate to service to handle PDF generation
	fileName, err := osc.service.GetOrganizationUnitReportExport(uint(ID), uint(currentUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(osc.db, c, fmt.Sprintf("Failed to export Organization unit Report PDF: %s", err.Error()))
			return
		}
	}

	// Set headers to ensure the file is served with the correct file name
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")

	// Serve the PDF file
	c.File(fileName)
	// Cleanup after the response has been sent
	c.Writer.Flush()

	go func() {
		time.Sleep(2 * time.Second) // Add a small delay to ensure file is no longer in use
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Successfully deleted file %s\n", fileName)
		}
	}()
}
