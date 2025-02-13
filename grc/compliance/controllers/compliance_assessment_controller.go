package controllers

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentController struct {
	service services.ComplianceAssessmentServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentController(service services.ComplianceAssessmentServiceInterface) *ComplianceAssessmentController {
	return &ComplianceAssessmentController{
		service: service,
	}
}

func (cac *ComplianceAssessmentController) GetAllComplianceAssessments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "id")
	status := c.Query("status")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_register_ids",
		"organization_unit_ids",
		"review_user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"id":                           "id",
		"organization_unit_title_en":   "organization_units.title_en",
		"organization_unit_title_ar":   "organization_units.title_ar",
		"compliance_register_title_en": "compliance_register.title_en",
		"compliance_register_title_ar": "compliance_register.title_ar",
		"review_user_en":               "review_user.first_name_en",
		"review_user_ar":               "review_user.first_name_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := complianceDTOs.ComplianceAssessmentListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	complianceAssessments, err := cac.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	// Otherwise, return the data as JSON
	c.JSON(http.StatusOK, complianceAssessments)
}

func (cac *ComplianceAssessmentController) GetComplianceAssessmentDetails(c *gin.Context) {
	idStr := c.Param("assessmentId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceAssessment, err := cac.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceAssessment)
}

func (cac *ComplianceAssessmentController) CreateComplianceAssessment(c *gin.Context) {
	var complianceAssessmentFormDTO complianceDTOs.ComplianceAssessmentFormDTO

	if err := c.ShouldBindJSON(&complianceAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.Create(currentUserID, &complianceAssessmentFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (cac *ComplianceAssessmentController) UpdateComplianceAssessment(c *gin.Context) {
	var complianceAssessmentFormDTO complianceDTOs.ComplianceAssessmentUpdateFormDTO

	complianceAssessmentID, err := strconv.Atoi(c.Param("assessmentId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.Update(uint(complianceAssessmentID), currentUserID, &complianceAssessmentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (cac *ComplianceAssessmentController) DeleteComplianceAssessments(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	failedIDs, err := cac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cac *ComplianceAssessmentController) GetComplianceAssessmentSummary(c *gin.Context) {
	idStr := c.Param("assessmentId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assessmentCompletionSummary, err := cac.service.AssessmentCompletionSummary(uint(id))
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, assessmentCompletionSummary)
}

func (cac *ComplianceAssessmentController) GetComplianceStatusSummary(c *gin.Context) {
	idStr := c.Param("assessmentId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusCounts, err := cac.service.GetComplianceStatusSummary(uint(id))
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusCounts)
}

func (cac *ComplianceAssessmentController) GetComplianceScoreSummary(c *gin.Context) {
	idStr := c.Param("assessmentId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusCounts, err := cac.service.GetComplianceScoreSummary(uint(id))
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusCounts)
}

func (cac *ComplianceAssessmentController) ComplianceAssessmentSubmit(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid assessment ID"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.SubmitComplianceAssessment(uint(assessmentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (cac *ComplianceAssessmentController) ComplianceAssessmentApprove(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid assessment ID"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.ApproveComplianceAssessment(uint(assessmentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
