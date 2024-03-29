package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceRequirementController struct {
	service services.ComplianceRequirementServiceInterface
	db      *gorm.DB
}

func NewComplianceRequirementController(service services.ComplianceRequirementServiceInterface) *ComplianceRequirementController {
	return &ComplianceRequirementController{
		service: service,
	}
}

func (crc *ComplianceRequirementController) GetAllComplianceRequirements(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_requirement_category_ids",
		"compliance_clause_ids",
		"severity_level_ids",
		"compliance_status_ids",
		"responsible_user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse deadline_from, deadline_to
	deadlineFrom := c.Query("deadline_from")
	deadlineTo := c.Query("deadline_to")

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

	queryParams := complianceDTOs.ComplianceRequirementListQueryParamsDTO{
		Sort:         sort,
		SortBy:       sortBy,
		Keyword:      keyword,
		Status:       status,
		Page:         page,
		PageSize:     pageSize,
		Filters:      filters,
		QuickFilter:  quickFilter,
		DeadlineFrom: deadlineFrom,
		DeadlineTo:   deadlineTo,
	}
	compliances, err := crc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, compliances)
}

func (crc *ComplianceRequirementController) GetComplianceRequirementDetails(c *gin.Context) {
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

func (crc *ComplianceRequirementController) CreateComplianceRequirement(c *gin.Context) {

	var ComplianceRequirementFormDTO complianceDTOs.ComplianceRequirementFormDTO

	if err := c.ShouldBindJSON(&ComplianceRequirementFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	complianceIDStr := c.Param("complianceId")
	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid complianceId"})
		return
	}

	id, err := crc.service.Create(uint(complianceID), currentUserID, &ComplianceRequirementFormDTO)

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

func (crc *ComplianceRequirementController) UpdateComplianceRequirement(c *gin.Context) {

	complianceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var ComplianceRequirementFormDTO complianceDTOs.ComplianceRequirementFormDTO
	if err := c.ShouldBindJSON(&ComplianceRequirementFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.Update(uint(complianceID), currentUserID, &ComplianceRequirementFormDTO)
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

func (crc *ComplianceRequirementController) DeleteComplianceRequirements(c *gin.Context) {

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
