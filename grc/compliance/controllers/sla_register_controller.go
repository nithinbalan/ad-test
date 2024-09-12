package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLARegisterController struct {
	service services.SLARegisterServiceInterface
	db      *gorm.DB
}

func NewSLARegisterController(service services.SLARegisterServiceInterface) *SLARegisterController {
	return &SLARegisterController{
		service: service,
	}
}

func (src *SLARegisterController) GetAllSLARegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(src.db, c, "Invalid export parameter")
		return
	}
	// Define an array of allowed parameter names
	allowedParams := []string{
		"sla_category_ids",
		"severity_level_ids",
		"responsible_unit_ids",
		"lifecycle_status_ids",
		"compliance_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse  issue_date_from, issue_date_to, expiry_date_from, expiry_date_to parameters
	issueDateFrom := c.Query("issue_date_from")
	issueDateTo := c.Query("issue_date_to")
	expiryDateFrom := c.Query("expiry_date_from")
	expiryDateTo := c.Query("expiry_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"reference_number":     "sla_register.reference_number",
		"title_en":             "sla_register.title_en",
		"title_ar":             "sla_register.title_ar",
		"description":          "sla_register.description",
		"client":               "sla_register.client",
		"sla_category_en":      "sla_categories.title_en",
		"sla_category_ar":      "sla_categories.title_ar",
		"severity_level_en":    "severity_levels.title_en",
		"severity_level_ar":    "severity_levels.title_ar",
		"responsible_unit_en":  "organization_units.title_en",
		"responsible_unit_ar":  "organization_units.title_ar",
		"responsible_user_en":  "responsible_user.first_name_en",
		"responsible_user_ar":  "responsible_user.first_name_ar",
		"status_en":            "lifecycle_statuses.title_en",
		"status_ar":            "lifecycle_statuses.title_ar",
		"compliance_status_en": "compliance_statuses.title_en",
		"compliance_status_ar": "compliance_statuses.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := slaDTOs.SLARegisterListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		IssueDateFrom:  issueDateFrom,
		IssueDateTo:    issueDateTo,
		ExpiryDateFrom: expiryDateFrom,
		ExpiryDateTo:   expiryDateTo,
		Export:         export,
	}
	slaRegisters, err := src.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	// If export is true, generate and return the Excel file
	if export {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(src.db, c, err.Error())
			return
		}

		excelFile, err := src.service.ExportSLARegisters(slaRegisters.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(src.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="SLARegister.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, slaRegisters)

}

func (src *SLARegisterController) GetSLARegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sla, err := src.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, sla)
}

func (src *SLARegisterController) CreateSLARegister(c *gin.Context) {
	var slaRegFormDTO slaDTOs.SLARegisterFormDTO

	if err := c.ShouldBindJSON(&slaRegFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	id, err := src.service.Create(currentUserID, &slaRegFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (src *SLARegisterController) UpdateSLARegister(c *gin.Context) {
	var slaRegFormDTO slaDTOs.SLARegisterFormDTO

	slaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaRegFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}
	err = src.service.Update(uint(slaID), currentUserID, &slaRegFormDTO)
	// @Review
	// We have to manage file upload failure error message,
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (src *SLARegisterController) DeleteSLARegisters(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	failedIDs, err := src.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	response := commonDTOs.FailedIDsResponse{
		FailedIds: failedIDs,
	}

	c.JSON(http.StatusOK, response)
}

func (src *SLARegisterController) ArchiveSLARegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	successfulIds, err := src.service.Archive(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	response := commonDTOs.SuccessfulIDsResponse{
		SuccessfulIds: successfulIds,
	}

	c.JSON(http.StatusOK, response)
}

func (src *SLARegisterController) RestoreSLARegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	successfulIds, err := src.service.Restore(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	response := commonDTOs.SuccessfulIDsResponse{
		SuccessfulIds: successfulIds,
	}

	c.JSON(http.StatusOK, response)
}

func (src *SLARegisterController) GetSLARegisterDocumentPreview(c *gin.Context) {
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

	slaDocumentPath, err := src.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(src.db, c, err.Error())
		}
		return
	}
	c.File(slaDocumentPath)
}

func (src *SLARegisterController) DownloadSLARegisterDocument(c *gin.Context) {
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

	slaDocumentPath, err := src.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(src.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(slaDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(slaDocumentPath)
}

func (src *SLARegisterController) AssessmentStatusSummary(c *gin.Context) {
	idStr := c.Param("slaId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assessmentStatusSummary, err := src.service.AssessmentStatusSummary(uint(id))
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, assessmentStatusSummary)
}

func (src *SLARegisterController) GetSLARegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	src.GetAllSLARegisters(c)
}
