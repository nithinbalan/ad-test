package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"
	"strings"

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
	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_category_ids",
		"severity_level_ids",
		"compliance_status_ids",
		"responsible_unit_ids",
		"review_frequency_ids",
		"review_responsible_user_ids",
		"lifecycle_status_ids",
	}

	// Parse dynamic filter parameters
	filters := make(map[string][]string)
	for key, values := range c.Request.URL.Query() {
		// Check if the key is one of the allowed parameters
		for _, allowedParam := range allowedParams {
			if key == allowedParam {
				// Extract column name
				columnName := strings.TrimSuffix(key, "s") // Remove the 's' at the end of the key
				if len(columnName) > 0 {
					// Split values by comma
					valueList := strings.Split(values[0], ",")
					// Add filtering condition only if the column name is not empty
					filters[columnName] = valueList
				}
			}
		}
	}

	// Parse quick filters
	quickFilter := c.QueryArray("quick_filter")

	// Parse deadline_from, deadline_to, issue_date_from, issue_date_to, expiry_date_from, expiry_date_to, next_review_date_from, next_review_date_to parameters
	deadlineFrom := c.Query("deadline_from")
	deadlineTo := c.Query("deadline_to")
	issueDateFrom := c.Query("issue_date_from")
	issueDateTo := c.Query("issue_date_to")
	expiryDateFrom := c.Query("expiry_date_from")
	expiryDateTo := c.Query("expiry_date_to")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

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

	compliances, err := crc.service.GetAll(sort, sortBy, keyword, status, page, pageSize, filters, quickFilter, deadlineFrom, deadlineTo, issueDateFrom, issueDateTo, expiryDateFrom, expiryDateTo, nextReviewDateFrom, nextReviewDateTo)

	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, compliances)
}

func (crc *ComplianceRegisterController) GetComplianceRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	compliance, err := crc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, compliance)
}

func (crc *ComplianceRegisterController) CreateComplianceRegister(c *gin.Context) {

	var complianceRegFormDTO complianceDTOs.ComplianceRegisterFormDTO

	if err := c.ShouldBindJSON(&complianceRegFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	id, err := crc.service.Create(currentUserID, &complianceRegFormDTO)

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

	complianceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var complianceRegFormDTO complianceDTOs.ComplianceRegisterFormDTO
	if err := c.ShouldBindJSON(&complianceRegFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.Update(uint(complianceID), currentUserID, &complianceRegFormDTO)
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
