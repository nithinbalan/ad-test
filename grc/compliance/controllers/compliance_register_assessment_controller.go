package controllers

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceRegisterAssessmentController struct {
	service services.ComplianceAssessmentServiceInterface
	db      *gorm.DB
}

func NewComplianceRegisterAssessmentController(service services.ComplianceAssessmentServiceInterface) *ComplianceRegisterAssessmentController {
	return &ComplianceRegisterAssessmentController{
		service: service,
	}
}

func (crac *ComplianceRegisterAssessmentController) GetAllComplianceRegisterAssessments(c *gin.Context) {
	complianceIDStr := c.Param("id")
	_, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "id")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(crac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(crac.db, c, "Invalid page_size parameter")
		return
	}

	filters := map[string][]string{
		"compliance_register_ids": {complianceIDStr},
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

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crac.db, c, err.Error())
		return
	}

	queryParams := complianceDTOs.ComplianceAssessmentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      false,
	}

	assessments, err := crac.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, assessments)
}

func (crac *ComplianceRegisterAssessmentController) CreateComplianceRegisterAssessment(c *gin.Context) {
	complianceIDStr := c.Param("id")
	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var complianceAssessmentFormDTO complianceDTOs.ComplianceAssessmentFormDTO
	if err := c.ShouldBindJSON(&complianceAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Set compliance register ID from URL
	complianceRegisterID := uint(complianceID)
	complianceAssessmentFormDTO.ComplianceRegisterID = &complianceRegisterID

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crac.db, c, err.Error())
		return
	}

	err = crac.service.Create(currentUserID, &complianceAssessmentFormDTO)
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
			utility.SendInternalErrResp(crac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (crac *ComplianceRegisterAssessmentController) DeleteComplianceRegisterAssessments(c *gin.Context) {
	complianceIDStr := c.Param("id")
	_, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crac.db, c, err.Error())
		return
	}

	failedIDs, err := crac.service.ForceDelete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(crac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
