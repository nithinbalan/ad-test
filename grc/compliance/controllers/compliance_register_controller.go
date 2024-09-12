package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceRegisterController struct {
	service services.ComplianceRegisterServiceInterface
	db      *gorm.DB
}

func NewComplianceRegisterController(service services.ComplianceRegisterServiceInterface) *ComplianceRegisterController {
	return &ComplianceRegisterController{
		service: service,
	}
}

func (crc *ComplianceRegisterController) GetAllComplianceRegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_category_ids",
		"severity_level_ids",
		"compliance_status_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"review_frequency_ids",
		"review_responsible_user_ids",
		"lifecycle_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse deadline_from, deadline_to, issue_date_from, issue_date_to, expiry_date_from, expiry_date_to, next_review_date_from, next_review_date_to parameters
	deadlineFrom := c.Query("deadline_from")
	deadlineTo := c.Query("deadline_to")
	issueDateFrom := c.Query("issue_date_from")
	issueDateTo := c.Query("issue_date_to")
	expiryDateFrom := c.Query("expiry_date_from")
	expiryDateTo := c.Query("expiry_date_to")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"reference_number":     "compliance_register.reference_number",
		"title_en":             "compliance_register.title_en",
		"title_ar":             "compliance_register.title_ar",
		"category_title_en":    "compliance_categories.title_en",
		"category_title_ar":    "compliance_categories.title_ar",
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
	queryParams := complianceDTOs.ComplianceRegisterListQueryParamsDTO{
		Sort:               sort,
		SortBy:             sortField,
		Keyword:            keyword,
		Status:             status,
		Page:               page,
		PageSize:           pageSize,
		Filters:            filters,
		QuickFilter:        quickFilter,
		DeadlineFrom:       deadlineFrom,
		DeadlineTo:         deadlineTo,
		IssueDateFrom:      issueDateFrom,
		IssueDateTo:        issueDateTo,
		ExpiryDateFrom:     expiryDateFrom,
		ExpiryDateTo:       expiryDateTo,
		NextReviewDateFrom: nextReviewDateFrom,
		NextReviewDateTo:   nextReviewDateTo,
		Export:             export,
	}

	compliances, err := crc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	// If export is true, generate and return the Excel file
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, err.Error())
			return
		}

		excelFile, err := crc.service.ExportComplianceRegisters(compliances.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Compliance Register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	// Otherwise, return the data as JSON
	c.JSON(http.StatusOK, compliances)
}

func (crc *ComplianceRegisterController) GetComplianceRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	compliance, err := crc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, compliance)
}

func (crc *ComplianceRegisterController) CreateComplianceRegister(c *gin.Context) {
	var complianceRegisterFormDTO complianceDTOs.ComplianceRegisterFormDTO

	if err := c.ShouldBindJSON(&complianceRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	id, err := crc.service.Create(currentUserID, &complianceRegisterFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (crc *ComplianceRegisterController) UpdateComplianceRegister(c *gin.Context) {
	var complianceRegisterFormDTO complianceDTOs.ComplianceRegisterFormDTO

	complianceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.Update(uint(complianceID), currentUserID, &complianceRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crc *ComplianceRegisterController) DeleteComplianceRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	failedIDs, err := crc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (crc *ComplianceRegisterController) ArchiveComplianceRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	successfulIds, err := crc.service.Archive(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}

func (crc *ComplianceRegisterController) RestoreComplianceRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	successfulIds, err := crc.service.Restore(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}

func (crc *ComplianceRegisterController) GetComplianceRegisterDocumentPreview(c *gin.Context) {
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

	complianceDocumentPath, err := crc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crc.db, c, err.Error())
		}
		return
	}
	c.File(complianceDocumentPath)
}

func (crc *ComplianceRegisterController) DownloadComplianceRegisterDocument(c *gin.Context) {
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

	complianceDocumentPath, err := crc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(complianceDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(complianceDocumentPath)
}

func (crc *ComplianceRegisterController) GetComplianceRegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	crc.GetAllComplianceRegisters(c)
}
