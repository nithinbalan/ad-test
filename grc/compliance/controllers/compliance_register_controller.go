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
	sortBy := c.DefaultQuery("sort_by", "created_at")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(crc.db, c, "Invalid export parameter")
			return
		}
	}

	// Define an array of allowed parameter names (matching repository available fields)
	allowedParams := []string{
		"compliance_status_ids",
		"compliance_folder_ids",
		"legislative_body_ids",
		"legislation_type_ids",
		"importance_classification_ids",
		"secrecy_classification_ids",
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

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

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
		"reference_number":                   "compliance_register.reference_number",
		"standard_code":                      "compliance_register.standard_code",
		"title_en":                           "compliance_register.title_en",
		"title_ar":                           "compliance_register.title_ar",
		"description_en":                     "compliance_register.description_en",
		"description_ar":                     "compliance_register.description_ar",
		"deadline":                           "compliance_register.deadline",
		"compliance_status_en":               "compliance_statuses.title_en",
		"compliance_status_ar":               "compliance_statuses.title_ar",
		"created_at":                         "compliance_register.created_at",
		"updated_at":                         "compliance_register.updated_at",
		"legislative_body_title_en":          "legislative_bodies.title_en",
		"legislative_body_title_ar":          "legislative_bodies.title_ar",
		"legislation_type_title_en":          "legislation_types.title_en",
		"legislation_type_title_ar":          "legislation_types.title_ar",
		"importance_classification_title_en": "importance_classifications.title_en",
		"importance_classification_title_ar": "importance_classifications.title_ar",
		"secrecy_classification_title_en":    "secrecy_classifications.title_en",
		"secrecy_classification_title_ar":    "secrecy_classifications.title_ar",
		"created_user_first_name_en":         "created_user.first_name_en",
		"created_user_first_name_ar":         "created_user.first_name_ar",
		"updated_user_first_name_en":         "updated_user.first_name_en",
		"updated_user_first_name_ar":         "updated_user.first_name_ar",
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
		DateFrom:           dateFrom,
		DateTo:             dateTo,
		Export:             export,
	}
	if isControlStr := c.Query("is_control_assessment"); isControlStr != "" {
		val, err := strconv.ParseBool(isControlStr)
		if err == nil {
			queryParams.IsControlAssessment = &val
		}
	}
	if isRequirementStr := c.Query("is_requirement_assessment"); isRequirementStr != "" {
		val, err := strconv.ParseBool(isRequirementStr)
		if err == nil {
			queryParams.IsRequirementAssessment = &val
		}
	}
	if isRequirementStr := c.Query("is_latest"); isRequirementStr != "" {
		val, err := strconv.ParseBool(isRequirementStr)
		if err == nil {
			queryParams.IsLatest = &val
		}
	}
	if isRequirementStr := c.Query("is_parent"); isRequirementStr != "" {
		val, err := strconv.ParseBool(isRequirementStr)
		if err == nil {
			queryParams.IsParent = &val
		}
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
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
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

func (crc *ComplianceRegisterController) GetIssuingAuthorityDocumentPreview(c *gin.Context) {
	authorityIDStr := c.Param("authorityId")
	authorityID, err := strconv.ParseUint(authorityIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := crc.service.GetIssuingAuthorityDocumentPath(uint(authorityID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (crc *ComplianceRegisterController) DownloadIssuingAuthorityDocument(c *gin.Context) {
	authorityIDStr := c.Param("authorityId")
	authorityID, err := strconv.ParseUint(authorityIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := crc.service.GetIssuingAuthorityDocumentPath(uint(authorityID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(documentPath)
}

func (crc *ComplianceRegisterController) GetComplianceRegisterExport(c *gin.Context) {
	c.Set("export", true)

	crc.GetAllComplianceRegisters(c)
}

func (crc *ComplianceRegisterController) GetAllVersions(c *gin.Context) {
	complianceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	compliance, err := crc.service.GetAllVersions(uint(complianceID))
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

func (crc *ComplianceRegisterController) CreateComplianceVersion(c *gin.Context) {
	var complianceRegisterFormDTO complianceDTOs.ComplianceVersionFormDTO

	if err := c.ShouldBindJSON(&complianceRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	id, err := crc.service.CreateComplianceVersion(currentUserID, &complianceRegisterFormDTO)

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

func (crc *ComplianceRegisterController) DeleteForceComplianceRegisters(c *gin.Context) {
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

	failedIDs, err := crc.service.ForceDelete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
